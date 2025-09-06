package main

import (
	"github.com/rivo/tview"
	"os/exec"
	"strings"
)

func runCommand(cmdStr string, output *tview.TextView, app *tview.Application) {
	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)

	// Capture output instead of writing directly
	result, err := cmd.CombinedOutput()

	// Use QueueUpdateDraw for thread-safe UI updates
	app.QueueUpdateDraw(func() {
		if err != nil {
			output.Write([]byte(string(result) + "\n[ERROR] " + err.Error()))
		} else {
			output.Write(result)
		}
	})
}

func main() {
	app := tview.NewApplication()

	menu := tview.NewList().
		AddItem("Encrypt Disk", "Run encrypt.sh", '1', nil).
		AddItem("Generate Certificate", "Run cert_generator.py", '2', nil).
		AddItem("Verify Certificate", "Run verify_cert.py", '3', nil).
		AddItem("Exit", "Quit application", 'q', func() { app.Stop() })

	output := tview.NewTextView()
	output.SetBorder(true)
	output.SetTitle("Output")

	// Wire menu actions
	menu.SetSelectedFunc(func(idx int, mainText, secText string, shortcut rune) {
		// Clear output in main goroutine (thread-safe)
		output.Clear()

		switch idx {
		case 0:
			go runCommand("bash encrypt.sh", output, app)
		case 1:
			go runCommand("python3 cert_generator.py", output, app)
		case 2:
			go runCommand("python3 verify_cert.py", output, app)
		case 3:
			app.Stop()
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(menu, 0, 1, true).
		AddItem(output, 0, 3, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
