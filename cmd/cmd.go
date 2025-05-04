package cmd

import (
	"fmt"
	"os"

	"goedit/editor"
	"goedit/terminal" // Need KeyEsc constant
)

// commandFuncMap defines the mapping from command strings to functions.
var commandFuncMap = map[string]func(e *editor.Editor){
	"w":  func(e *editor.Editor) { _ = SaveFile(e) }, // Wrapper ignores return
	"wq": saveAndQuit,
	"q":  QuitEditor,
	"q!": quitWithoutSaving,
}

// ProcessCommandInput handles a single key press when in Command mode.
func ProcessCommandInput(e *editor.Editor, key byte) {
	switch key {
	case terminal.KeyEsc:
		e.CurrentMode = editor.ModeNormal
		e.CommandBuffer = "" // Clear buffer on escape
	case 13: // Enter
		originalMode := e.CurrentMode // Remember mode before command
		executeCommand(e)
		// If executeCommand didn't change the mode (e.g., unknown command
		// or a command like :set that doesn't prompt), return to Normal.
		// If it *did* change (like :w prompting), let the new mode stick.
		if e.CurrentMode == originalMode {
			e.CurrentMode = editor.ModeNormal
		}
	case 127, 8: // Backspace
		if len(e.CommandBuffer) > 0 {
			e.CommandBuffer = e.CommandBuffer[:len(e.CommandBuffer)-1]
		}
	default:
		if key >= 32 && key <= 126 { // Allow printable ASCII
			e.CommandBuffer += string(key)
		}
	}
}

// executeCommand looks up and runs the command in the command buffer.
func executeCommand(e *editor.Editor) { // Keep unexported, helper for ProcessCommandInput
	command := e.CommandBuffer // Store command before clearing
	e.CommandBuffer = ""       // Clear buffer for next command
	// Mode change back to Normal is handled by the input processing logic
	// after the command function returns, or by the command func itself if needed.
	// SetStatusMessage needs to happen *before* mode change if we want it visible.

	if cmdFunc, exists := commandFuncMap[command]; exists {
		cmdFunc(e) // This might set a status message or change mode
	} else {
		e.SetStatusMessage(fmt.Sprintf("Unknown command: %s", command))
		// If command unknown, explicitly return to Normal mode
		e.CurrentMode = editor.ModeNormal
	}
	// If cmdFunc changed the mode (e.g., to ModeFileNamePrompt),
	// let that mode persist for the next input loop iteration.
	// If cmdFunc did NOT change the mode, we should probably return to Normal here?
	// Let's assume command functions handle necessary mode changes or
	// that staying in Command mode briefly isn't an issue if the command
	// didn't explicitly change it (it gets reset on next ':' anyway).
	// Revisit if needed. Let's try without explicit reset first.

	// The most robust is perhaps for ProcessCommandInput to handle the mode reset
	// after calling executeCommand, IF executeCommand didn't change it.
	// Let's modify ProcessCommandInput slightly.
}

// --- Command Implementations ---

// SaveFile writes the editor content or prompts for filename if needed.
// Returns true if a save was attempted (success or error), false if prompt was initiated.
func SaveFile(e *editor.Editor) bool {
	if e.Filename == "" {
		// No filename, enter prompt mode
		e.CurrentMode = editor.ModeFileNamePrompt
		e.CommandBuffer = ""                 // Clear buffer for filename input
		e.SetStatusMessage("Save file as: ") // Initial prompt
		e.PromptOriginCommand = "w"          // Mark that :w triggered the prompt
		return false                         // Didn't attempt save, initiated prompt
	}

	// Filename exists, proceed with saving
	content := e.ContentAsString()
	err := os.WriteFile(e.Filename, []byte(content), 0644)
	if err != nil {
		e.SetStatusMessage(fmt.Sprintf("Error saving file: %v", err))
	} else {
		e.SetStatusMessage(fmt.Sprintf("File '%s' saved successfully.", e.Filename))
		e.IsDirty = false
	}
	return true // Attempted save
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
		return // Don't set ShouldQuit if dirty
	}
	e.ShouldQuit = true
}

// quitWithoutSaving signals the main loop to exit regardless of dirty state.
func quitWithoutSaving(e *editor.Editor) {
	e.ShouldQuit = true
}
