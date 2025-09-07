package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
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
	}
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
