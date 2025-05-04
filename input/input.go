package input

import (
	"time"

	"goedit/cmd"
	"goedit/editor"
	"goedit/terminal"
)

// ProcessInput routes the key press to the appropriate mode handler.
func ProcessInput(e *editor.Editor, key byte) {
	switch e.CurrentMode {
	case editor.ModeNormal:
		processNormalModeInput(e, key)
	case editor.ModeInsert:
		processInsertModeInput(e, key)
	case editor.ModeCommand:
		cmd.HandleCommandKey(e, key)
	case editor.ModeFileNamePrompt:
		processFileNamePrompt(e, key)
	}
}

// processNormalModeInput handles input when in Normal mode.
func processNormalModeInput(e *editor.Editor, key byte) {
	switch key {
	case 'q': // Do nothing (require :q)
	case 'i':
		e.CurrentMode = editor.ModeInsert
	case 'h', terminal.KeyArrowLeft:
		if e.CursorX > 0 {
			e.CursorX--
		}
	case 'j', terminal.KeyArrowDown:
		if e.CursorY < len(e.EditorContent)-1 {
			e.CursorY++
			if e.CursorY >= e.RowOffset+e.TermHeight-1 {
				e.RowOffset++
			}
			e.EnsureCursorBounds()
		}
	case 'k', terminal.KeyArrowUp:
		if e.CursorY > 0 {
			e.CursorY--
			if e.CursorY < e.RowOffset {
				e.RowOffset--
			}
			e.EnsureCursorBounds()
		}
	case 'l', terminal.KeyArrowRight:
		if e.CursorY < len(e.EditorContent) {
			lineLen := len(e.EditorContent[e.CursorY])
			if e.CursorX < lineLen {
				e.CursorX++
			}
		}
	case ':':
		e.CurrentMode = editor.ModeCommand
		e.CommandBuffer = ""
		e.SetStatusMessage("")
	}
}

// processInsertModeInput handles input when in Insert mode.
func processInsertModeInput(e *editor.Editor, key byte) {
	switch key {
	case terminal.KeyEsc:
		e.CurrentMode = editor.ModeNormal
		e.StatusMessageTime = time.Time{}
	case 13:
		e.InsertNewline()
	case 127, 8:
		e.DeleteChar()
	case terminal.KeyArrowUp:
		if e.CursorY > 0 {
			e.CursorY--
			e.EnsureCursorBounds()
		}
	case terminal.KeyArrowDown:
		if e.CursorY < len(e.EditorContent)-1 {
			e.CursorY++
			e.EnsureCursorBounds()
		}
	case terminal.KeyArrowLeft:
		if e.CursorX > 0 {
			e.CursorX--
		} else if e.CursorY > 0 {
			e.CursorY--
			if e.CursorY < len(e.EditorContent) {
				e.CursorX = len(e.EditorContent[e.CursorY])
			} else {
				e.CursorX = 0
			}
		}
	case terminal.KeyArrowRight:
		if e.CursorY < len(e.EditorContent) {
			lineLen := len(e.EditorContent[e.CursorY])
			if e.CursorX < lineLen {
				e.CursorX++
			} else if e.CursorY < len(e.EditorContent)-1 {
				e.CursorY++
				e.CursorX = 0
			}
		}
	default:
		if key >= 32 && key < 127 {
			e.InsertChar(key)
		}
	}
}

// processFileNamePrompt handles input when prompting for a filename to save.
func processFileNamePrompt(e *editor.Editor, key byte) {
	switch key {
	case terminal.KeyEsc:
		// Cancel prompt
		e.SetStatusMessage("Save aborted.")
		e.CurrentMode = editor.ModeNormal
		e.CommandBuffer = ""
	case 13:
		filename := e.CommandBuffer
		if filename == "" { // No filename entered
			e.SetStatusMessage("Save aborted.")
			e.CurrentMode = editor.ModeNormal
		} else {
			e.Filename = filename
			e.SetStatusMessage("")
			// Note: SaveFile might set its own status message ("saved" or "error")
			// Don't change mode here yet, let SaveFile finish
			attemptedSave := cmd.SaveFile(e)

			if attemptedSave {
				// If save was attempted (even if it failed), check if we need to quit.
				// Only quit if the original command was :wq.
				if e.PromptOriginCommand == "wq" {
					// quitEditor will check IsDirty flag. If save failed,
					// IsDirty is still true, and quit will be blocked (correctly).
					// If save succeeded, IsDirty is false, and quit will proceed.
					cmd.QuitEditor(e)
				}
				// If origin was :w, we don't quit here.
			} // If save wasn't attempted (shouldn't happen here), do nothing more.

			// Always return to normal mode after attempting save/handling quit from prompt
			e.CurrentMode = editor.ModeNormal
			e.PromptOriginCommand = "" // Clear origin after handling
		}
		e.CommandBuffer = ""
	case 127, 8:
		if len(e.CommandBuffer) > 0 {
			e.CommandBuffer = e.CommandBuffer[:len(e.CommandBuffer)-1]
			// Update prompt dynamically
			e.SetStatusMessage("Save file as: " + e.CommandBuffer)
		}
	default:
		if key >= 32 && key <= 126 {
			e.CommandBuffer += string(key)
			newPrompt := "Save file as: " + e.CommandBuffer
			e.SetStatusMessage(newPrompt)
		}
	}
}
