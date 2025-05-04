package main

import (
	"log"
	"os"

	"goedit/editor"
	"goedit/input"
	"goedit/terminal"
	"goedit/ui"
)

func main() {
	// Get initial terminal size for editor creation
	width, height, err := terminal.GetSize()
	if err != nil {
		log.Printf("Error getting terminal size on init: %v. Using defaults.", err)
		width = 80
		height = 24
	}

	// Initialize editor state using the new package
	ed := editor.NewEditor(width, height)

	// Handle file loading
	if len(os.Args) > 1 {
		ed.Filename = os.Args[1]
		content, err := os.ReadFile(ed.Filename)
		if err != nil {
			if !os.IsNotExist(err) { // Log errors other than file not found
				log.Printf("Error opening file '%s': %v", ed.Filename, err)
			} // If file doesn't exist, editor starts empty anyway
		} else {
			ed.LoadFile(content)
		}
	}

	// Enter raw mode and ensure it's disabled on exit
	originalState, err := terminal.EnableRawMode()
	if err != nil {
		log.Fatalf("Failed to enable raw mode: %v", err)
	}
	defer terminal.DisableRawMode(originalState)

	ui.RefreshScreen(ed)

	// Main input loop
	for {
		key := terminal.ReadKey()

		input.ProcessInput(ed, key)

		ui.RefreshScreen(ed)

		if ed.ShouldQuit {
			break // Exit the loop gracefully
		}
	}
}
