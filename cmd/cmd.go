package cmd

import (
	"fmt"
	"os"

	"goedit/editor"
	"goedit/terminal"
)

// commandFuncMap defines the mapping from command strings to functions.
var commandFuncMap = map[string]func(e *editor.Editor){
	"w":  func(e *editor.Editor) { _ = SaveFile(e) },
	"wq": saveAndQuit,
	"q":  QuitEditor,
	"q!": quitWithoutSaving,
}

// ProcessCommandInput handles a single key press when in Command mode.
func ProcessCommandInput(e *editor.Editor, key byte) {
	switch key {
	case terminal.KeyEsc:
		e.CurrentMode = editor.ModeNormal
		e.CommandBuffer = ""
	case 13:
		originalMode := e.CurrentMode
		executeCommand(e)
		if e.CurrentMode == originalMode {
			e.CurrentMode = editor.ModeNormal
		}
	case 127, 8:
		if len(e.CommandBuffer) > 0 {
			e.CommandBuffer = e.CommandBuffer[:len(e.CommandBuffer)-1]
		}
	default:
		if key >= 32 && key <= 126 {
			e.CommandBuffer += string(key)
		}
	}
}

// executeCommand looks up and runs the command in the command buffer.
func executeCommand(e *editor.Editor) {
	command := e.CommandBuffer
	e.CommandBuffer = ""

	if cmdFunc, exists := commandFuncMap[command]; exists {
		cmdFunc(e)
	} else {
		e.SetStatusMessage(fmt.Sprintf("Unknown command: %s", command))
		e.CurrentMode = editor.ModeNormal // If command unknown, explicitly return to Normal mode
	}
}

// SaveFile writes the editor content or prompts for filename if needed.
// Returns true if a save was attempted (success or error), false if prompt was initiated.
func SaveFile(e *editor.Editor) bool {
	if e.Filename == "" {
		e.CurrentMode = editor.ModeFileNamePrompt
		e.CommandBuffer = ""
		e.SetStatusMessage("Save file as: ")
		e.PromptOriginCommand = "w"
		return false
	}

	content := e.ContentAsString()
	err := os.WriteFile(e.Filename, []byte(content), 0644)
	if err != nil {
		e.SetStatusMessage(fmt.Sprintf("Error saving file: %v", err))
	} else {
		e.SetStatusMessage(fmt.Sprintf("File '%s' saved successfully.", e.Filename))
		e.IsDirty = false
	}
	return true
}

// saveAndQuit saves the file and then signals quit, only if save was attempted.
func saveAndQuit(e *editor.Editor) {
	attemptedSave := SaveFile(e)
	if attemptedSave {
		// Only quit if SaveFile actually tried to save (had filename)
		// and wasn't blocked by IsDirty flag itself.
		QuitEditor(e)
	} else {
		// SaveFile returned false, meaning it entered prompt mode.
		// Mark that :wq triggered this prompt.
		e.PromptOriginCommand = "wq"
	}
}

// QuitEditor signals the main loop to exit if buffer isn't dirty.
func QuitEditor(e *editor.Editor) {
	if e.IsDirty {
		e.SetStatusMessage("Unsaved changes! Use :q! or :wq to save and quit.")
		return
	}
	e.ShouldQuit = true
}

// quitWithoutSaving signals the main loop to exit regardless of dirty state.
func quitWithoutSaving(e *editor.Editor) {
	e.ShouldQuit = true
}
