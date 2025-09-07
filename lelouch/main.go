package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
	"zerotrace/lib"
)

// BlockDevice represents a storage device
type BlockDevice struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	Type  string `json:"type"`
	Model string `json:"model"`
}

// LsblkOutput represents the JSON structure from lsblk -J
type LsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

// Get all block devices
func getBlockDevices() ([]BlockDevice, error) {
	cmd := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,MODEL")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsblk: %v", err)
	}

	var lsblkOutput LsblkOutput
	if err := json.Unmarshal(output, &lsblkOutput); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// remove loop devices and partitions
	filteredDevices := []BlockDevice{}
	for _, device := range lsblkOutput.BlockDevices {
		if device.Type == "disk" {
			filteredDevices = append(filteredDevices, device)
		}
	}
	lsblkOutput.BlockDevices = filteredDevices

	// Add /dev/ prefix to device names
	for i := range lsblkOutput.BlockDevices {
		lsblkOutput.BlockDevices[i].Name = "/dev/" + lsblkOutput.BlockDevices[i].Name
	}

	return lsblkOutput.BlockDevices, nil
}

func runCommand(cmdStr string, output *tview.TextView, app *tview.Application) {
	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)

	// Create pipes for real-time output
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Start the command
	cmd.Start()

	// Stream stdout
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				text := string(buf[:n])
				// Clean up carriage returns
				text = strings.ReplaceAll(text, "\r", "\n")
				app.QueueUpdateDraw(func() {
					output.Write([]byte(text))
				})
			}
			if err != nil {
				break
			}
		}
	}()

	// Stream stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				text := string(buf[:n])
				// Clean up carriage returns
				text = strings.ReplaceAll(text, "\r", "\n")
				app.QueueUpdateDraw(func() {
					output.Write([]byte(text))
				})
			}
			if err != nil {
				break
			}
		}
	}()

	// Wait for command to finish
	err := cmd.Wait()
	if err != nil {
		app.QueueUpdateDraw(func() {
			output.Write([]byte(fmt.Sprintf("\n[ERROR] %s\n", err.Error())))
		})
	} else {
		// Command completed successfully, now sign the wipe data
		go func() {
			if err := signWipeData(output, app); err != nil {
				app.QueueUpdateDraw(func() {
					output.Write([]byte(fmt.Sprintf("\n[ERROR] Failed to sign wipe data: %s\n", err.Error())))
				})
			}
		}()
	}
}

// SignedWipeData represents the structure of the signed wipe data
type SignedWipeData struct {
	Data      interface{} `json:"data"`
	PublicKey string      `json:"publickey"`
	Signature string      `json:"signature"`
}

// signWipeData reads the wipe_log.json, signs it, and updates the file
func signWipeData(output *tview.TextView, app *tview.Application) error {
	// Read the existing wipe_log.json
	wipeLogPath := "wipe_log.json"
	data, err := ioutil.ReadFile(wipeLogPath)
	if err != nil {
		return fmt.Errorf("failed to read wipe_log.json: %v", err)
	}

	// Parse the existing wipe data
	var wipeData interface{}
	if err := json.Unmarshal(data, &wipeData); err != nil {
		return fmt.Errorf("failed to parse wipe_log.json: %v", err)
	}

	// Generate a new key pair
	keyPair, err := lib.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %v", err)
	}

	// Get the public key in PEM format
	publicKeyPEM, err := keyPair.PublicKeyToPEM()
	if err != nil {
		return fmt.Errorf("failed to convert public key to PEM: %v", err)
	}

	// Sign the wipe data
	signature, err := keyPair.SignData(wipeData)
	if err != nil {
		return fmt.Errorf("failed to sign wipe data: %v", err)
	}

	// Create the signed wipe data structure
	signedData := SignedWipeData{
		Data:      wipeData,
		PublicKey: publicKeyPEM,
		Signature: signature,
	}

	// Convert to JSON
	signedJSON, err := json.MarshalIndent(signedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal signed data: %v", err)
	}

	// Write back to the file
	if err := ioutil.WriteFile(wipeLogPath, signedJSON, 0644); err != nil {
		return fmt.Errorf("failed to write signed data to file: %v", err)
	}

	// Upload the certificate to 0x0.st
	url, err := lib.UploadCert(wipeLogPath)
	if err != nil {
		return fmt.Errorf("failed to upload certificate: %v", err)
	}

	// Generate ASCII QR code for the certificate URL
	qrASCII, err := lib.GenerateQRCodeASCII(url)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Display the certificate URL and QR code in the output view
	app.QueueUpdateDraw(func() {
		output.Write([]byte(fmt.Sprintf("\n\nHere is your cert: %s\n\n", url)))
		output.Write([]byte("QR Code:\n"))
		output.Write([]byte(qrASCII))
		output.Write([]byte("\n"))
	})

	// Remove the wipe_log.json file after successful upload
	if err := os.Remove(wipeLogPath); err != nil {
		// Log the error but don't fail the entire process
		app.QueueUpdateDraw(func() {
			output.Write([]byte(fmt.Sprintf("[WARNING] Failed to remove %s: %s\n", wipeLogPath, err.Error())))
		})
	}

	return nil
}

func main() {
	app := tview.NewApplication()

	// Get devices
	devices, err := getBlockDevices()
	if err != nil {
		panic(err)
	}

	// Create list
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Block Devices - Select a Drive")

	// Create output view
	output := tview.NewTextView()
	output.SetBorder(true)
	output.SetTitle("Output")

	// Create flex layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(output, 0, 1, false)

	for i, device := range devices {
		model := device.Model
		if model == "" {
			model = "Unknown"
		}

		// Capture device in closure
		currentDevice := device

		list.AddItem(
			currentDevice.Name,
			fmt.Sprintf("%s - %s - %s", currentDevice.Size, currentDevice.Type, model),
			rune('1'+i),
			func() {
				// Clear previous output
				output.Clear()
				// Fixed command format - pass device name as argument
				cmd := fmt.Sprintf("bash ./scripts/wipe.sh %s", currentDevice.Name)
				// Run command in background
				go runCommand(cmd, output, app)
			},
		)
	}

	list.AddItem("Exit", "Press to exit", 'q', func() {
		app.Stop()
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
