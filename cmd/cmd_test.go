package cmd

import (
	"path/filepath"
	"testing"

	"goedit/editor"
)

// Helper to create a test editor instance
func newTestEditor(isDirty bool) *editor.Editor {
	ed := editor.NewEditor(80, 24) // Size arbitrary for these tests
	ed.IsDirty = isDirty
	ed.ShouldQuit = false // Ensure start state
	ed.StatusMessage = "" // Ensure start state
	return ed
}

func TestQuitEditor(t *testing.T) {
	t.Run("Quit when not dirty", func(t *testing.T) {
		ed := newTestEditor(false)
		QuitEditor(ed)

		if !ed.ShouldQuit {
			t.Errorf("Expected ShouldQuit to be true when not dirty, got false")
		}
		if ed.StatusMessage != "" {
			t.Errorf("Expected StatusMessage to be empty, got %q", ed.StatusMessage)
		}
	})

	t.Run("Do not quit when dirty", func(t *testing.T) {
		ed := newTestEditor(true)
		QuitEditor(ed)

		if ed.ShouldQuit {
			t.Errorf("Expected ShouldQuit to be false when dirty, got true")
		}
		if ed.StatusMessage == "" {
			t.Errorf("Expected StatusMessage to be set when dirty and quitting, but it was empty")
		}
	})
}

func TestQuitWithoutSaving(t *testing.T) {
	t.Run("Force quit when not dirty", func(t *testing.T) {
		ed := newTestEditor(false)
		quitWithoutSaving(ed)

		if !ed.ShouldQuit {
			t.Errorf("Expected ShouldQuit to be true, got false")
		}
	})

	t.Run("Force quit when dirty", func(t *testing.T) {
		ed := newTestEditor(true)
		quitWithoutSaving(ed)

		if !ed.ShouldQuit {
			t.Errorf("Expected ShouldQuit to be true even when dirty, got false")
		}
	})
}

func TestSaveFile(t *testing.T) {
	// Extend helper to optionally set filename
	newTestEditorWithFile := func(filename string, content []string) *editor.Editor {
		ed := editor.NewEditor(80, 24)
		ed.Filename = filename
		ed.EditorContent = content
		ed.IsDirty = true // Assume dirty before save
		return ed
	}

	t.Run("Save with no filename (enter prompt mode)", func(t *testing.T) {
		ed := newTestEditorWithFile("", []string{"content"})
		ed.CurrentMode = editor.ModeCommand // Start in command mode for realism

		attemptedSave := SaveFile(ed)

		if attemptedSave {
			t.Error("Expected SaveFile to return false when no filename exists, got true")
		}
		if ed.CurrentMode != editor.ModeFileNamePrompt {
			t.Errorf("Expected mode to be ModeFileNamePrompt, got %v", ed.CurrentMode)
		}
		if ed.StatusMessage != "Save file as: " {
			t.Errorf("Expected status message 'Save file as: ', got %q", ed.StatusMessage)
		}
		if ed.PromptOriginCommand != "w" {
			t.Errorf("Expected PromptOriginCommand to be 'w', got %q", ed.PromptOriginCommand)
		}
		if !ed.IsDirty { // Should still be dirty, save wasn't attempted
			t.Error("Expected IsDirty to remain true, got false")
		}
	})

	t.Run("Save with existing filename", func(t *testing.T) {
		// We cannot easily test the actual file write here without mocking os.WriteFile
		// or using temporary files. We focus on the state changes assuming write succeeds.
		// Use t.TempDir() for cleanup, even though we aren't writing yet.
		tempDir := t.TempDir()
		tempFilename := filepath.Join(tempDir, "test_save.txt") // Use path.Join

		ed := newTestEditorWithFile(tempFilename, []string{"line1", "line2"})
		initialMode := editor.ModeCommand
		ed.CurrentMode = initialMode

		attemptedSave := SaveFile(ed)

		if !attemptedSave {
			t.Error("Expected SaveFile to return true when filename exists, got false")
		}
		if ed.CurrentMode != initialMode { // Mode should not change
			t.Errorf("Expected mode to remain %v, got %v", initialMode, ed.CurrentMode)
		}
		if ed.IsDirty { // Should be false after successful save
			t.Error("Expected IsDirty to become false, got true")
		}
		if ed.StatusMessage == "" { // Should have a success message
			t.Error("Expected StatusMessage to be set on success, but it was empty")
		}
		// Note: Cannot assert exact success message content as it includes filename
	})
}
