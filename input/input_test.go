package input

import (
	"os"
	"strings"
	"testing"

	"goedit/editor"
	"goedit/terminal"
)

// Helper to create a test editor instance
func newTestEditor(content []string, cursorX, cursorY int) *editor.Editor {
	ed := editor.NewEditor(80, 24) // Use arbitrary dimensions
	ed.EditorContent = content
	ed.CursorX = cursorX
	ed.CursorY = cursorY
	ed.CurrentMode = editor.ModeNormal // Start in Normal mode for these tests
	return ed
}

func TestProcessNormalModeInput(t *testing.T) {
	tests := []struct {
		name            string
		initialContent  []string
		initialCursorX  int
		initialCursorY  int
		key             byte
		expectedMode    editor.Mode
		expectedCursorX int
		expectedCursorY int
		expectedCommand string // For : command
	}{
		// --- Mode Change ---
		{
			name:            "i enters Insert mode",
			initialContent:  []string{""},
			initialCursorX:  0,
			initialCursorY:  0,
			key:             'i',
			expectedMode:    editor.ModeInsert,
			expectedCursorX: 0, // Position doesn't change yet
			expectedCursorY: 0,
		},
		{
			name:            ": enters Command mode",
			initialContent:  []string{""},
			initialCursorX:  0,
			initialCursorY:  0,
			key:             ':',
			expectedMode:    editor.ModeCommand,
			expectedCursorX: 0,
			expectedCursorY: 0,
			expectedCommand: "", // Buffer starts empty
		},
		// --- Movement ---
		{
			name:            "h moves left",
			initialContent:  []string{"abc"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             'h',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 0,
			expectedCursorY: 0,
		},
		{
			name:            "h at start of line",
			initialContent:  []string{"abc"},
			initialCursorX:  0,
			initialCursorY:  0,
			key:             'h',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 0, // Stays at 0
			expectedCursorY: 0,
		},
		{
			name:            "l moves right",
			initialContent:  []string{"abc"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             'l',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 2,
			expectedCursorY: 0,
		},
		{
			name:            "l at end of line",
			initialContent:  []string{"abc"},
			initialCursorX:  3,
			initialCursorY:  0,
			key:             'l',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 3, // Stays at end
			expectedCursorY: 0,
		},
		{
			name:            "j moves down",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             'j',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 1, // EnsureCursorBounds will handle if line is shorter
			expectedCursorY: 1,
		},
		{
			name:            "j at bottom",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  1,
			key:             'j',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 1,
			expectedCursorY: 1, // Stays at bottom
		},
		{
			name:            "k moves up",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  1,
			key:             'k',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 1,
			expectedCursorY: 0,
		},
		{
			name:            "k at top",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             'k',
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 1,
			expectedCursorY: 0, // Stays at top
		},
		// Arrow keys (should behave same as hjkl in Normal mode)
		{
			name:            "ArrowLeft moves left",
			initialContent:  []string{"abc"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             terminal.KeyArrowLeft,
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 0,
			expectedCursorY: 0,
		},
		{
			name:            "ArrowRight moves right",
			initialContent:  []string{"abc"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             terminal.KeyArrowRight,
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 2,
			expectedCursorY: 0,
		},
		{
			name:            "ArrowDown moves down",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             terminal.KeyArrowDown,
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 1,
			expectedCursorY: 1,
		},
		{
			name:            "ArrowUp moves up",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  1,
			key:             terminal.KeyArrowUp,
			expectedMode:    editor.ModeNormal,
			expectedCursorX: 1,
			expectedCursorY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditor(tt.initialContent, tt.initialCursorX, tt.initialCursorY)

			// processNormalModeInput is not exported, so we call ProcessInput
			// ensuring the editor starts in Normal mode.
			if ed.CurrentMode != editor.ModeNormal {
				t.Fatalf("Test setup error: Editor should start in Normal mode")
			}
			ProcessInput(ed, tt.key)

			// Check Mode
			if ed.CurrentMode != tt.expectedMode {
				t.Errorf("Expected Mode %v, got %v", tt.expectedMode, ed.CurrentMode)
			}

			// Check Cursor X
			if ed.CursorX != tt.expectedCursorX {
				t.Errorf("Expected CursorX %d, got %d", tt.expectedCursorX, ed.CursorX)
			}

			// Check Cursor Y
			if ed.CursorY != tt.expectedCursorY {
				t.Errorf("Expected CursorY %d, got %d", tt.expectedCursorY, ed.CursorY)
			}

			// Check command buffer if relevant
			if tt.key == ':' && ed.CommandBuffer != tt.expectedCommand {
				t.Errorf("Expected CommandBuffer %q, got %q", tt.expectedCommand, ed.CommandBuffer)
			}
		})
	}
}

func TestProcessInsertModeInput(t *testing.T) {
	// Helper
	newTestEditorInsert := func(content []string, cursorX, cursorY int) *editor.Editor {
		ed := newTestEditor(content, cursorX, cursorY) // Use helper from Normal test
		ed.CurrentMode = editor.ModeInsert             // Start in Insert mode
		return ed
	}

	tests := []struct {
		name            string
		initialContent  []string
		initialCursorX  int
		initialCursorY  int
		key             byte
		expectedMode    editor.Mode
		expectedContent []string // Check basic content change for insert/delete
		expectedCursorX int
		expectedCursorY int
	}{
		// --- Mode Change ---
		{
			name:            "Esc enters Normal mode",
			initialContent:  []string{"abc"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             terminal.KeyEsc,
			expectedMode:    editor.ModeNormal,
			expectedContent: []string{"abc"}, // Content doesn't change
			expectedCursorX: 1,
			expectedCursorY: 0,
		},
		// --- Editing --- (Basic check, full logic tested in editor_test)
		{
			name:            "Printable char inserts",
			initialContent:  []string{"ac"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             'b',
			expectedMode:    editor.ModeInsert,
			expectedContent: []string{"abc"},
			expectedCursorX: 2,
			expectedCursorY: 0,
		},
		{
			name:            "Enter creates newline",
			initialContent:  []string{"ab"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             13, // Enter key code
			expectedMode:    editor.ModeInsert,
			expectedContent: []string{"a", "b"}, // Basic check
			expectedCursorX: 0,
			expectedCursorY: 1,
		},
		{
			name:            "Backspace deletes char",
			initialContent:  []string{"abc"},
			initialCursorX:  2,
			initialCursorY:  0,
			key:             127, // Backspace key code
			expectedMode:    editor.ModeInsert,
			expectedContent: []string{"ac"}, // Basic check
			expectedCursorX: 1,
			expectedCursorY: 0,
		},
		// --- Movement ---
		{
			name:            "ArrowLeft moves left in insert",
			initialContent:  []string{"abc"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             terminal.KeyArrowLeft,
			expectedMode:    editor.ModeInsert,
			expectedContent: []string{"abc"},
			expectedCursorX: 0,
			expectedCursorY: 0,
		},
		{
			name:            "ArrowRight moves right in insert",
			initialContent:  []string{"abc"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             terminal.KeyArrowRight,
			expectedMode:    editor.ModeInsert,
			expectedContent: []string{"abc"},
			expectedCursorX: 2,
			expectedCursorY: 0,
		},
		{
			name:            "ArrowDown moves down in insert",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  0,
			key:             terminal.KeyArrowDown,
			expectedMode:    editor.ModeInsert,
			expectedContent: []string{"line1", "line2"},
			expectedCursorX: 1, // EnsureCursorBounds handles length
			expectedCursorY: 1,
		},
		{
			name:            "ArrowUp moves up in insert",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  1,
			initialCursorY:  1,
			key:             terminal.KeyArrowUp,
			expectedMode:    editor.ModeInsert,
			expectedContent: []string{"line1", "line2"},
			expectedCursorX: 1,
			expectedCursorY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditorInsert(tt.initialContent, tt.initialCursorX, tt.initialCursorY)

			// processInsertModeInput is not exported, call ProcessInput
			if ed.CurrentMode != editor.ModeInsert {
				t.Fatalf("Test setup error: Editor should start in Insert mode")
			}
			ProcessInput(ed, tt.key)

			// Check Mode
			if ed.CurrentMode != tt.expectedMode {
				t.Errorf("Expected Mode %v, got %v", tt.expectedMode, ed.CurrentMode)
			}

			// Check Cursor X
			if ed.CursorX != tt.expectedCursorX {
				t.Errorf("Expected CursorX %d, got %d", tt.expectedCursorX, ed.CursorX)
			}

			// Check Cursor Y
			if ed.CursorY != tt.expectedCursorY {
				t.Errorf("Expected CursorY %d, got %d", tt.expectedCursorY, ed.CursorY)
			}

			// Check Content (basic check for edit keys)
			if len(ed.EditorContent) != len(tt.expectedContent) {
				t.Fatalf("Expected %d lines, got %d", len(tt.expectedContent), len(ed.EditorContent))
			}
			for i := range tt.expectedContent {
				if i >= len(ed.EditorContent) {
					t.Errorf("Line %d missing: expected %q", i, tt.expectedContent[i])
					continue
				}
				if ed.EditorContent[i] != tt.expectedContent[i] {
					t.Errorf("Line %d: Expected %q, got %q", i, tt.expectedContent[i], ed.EditorContent[i])
				}
			}
		})
	}
}

func TestProcessFileNamePrompt(t *testing.T) {
	// Helper
	newTestEditorPrompt := func(buffer string, origin string) *editor.Editor {
		// Using Normal mode helper is fine
		ed := newTestEditor([]string{"content"}, 0, 0)
		ed.CurrentMode = editor.ModeFileNamePrompt
		ed.CommandBuffer = buffer
		ed.PromptOriginCommand = origin
		ed.SetStatusMessage("Save file as: " + buffer) // Mimic initial prompt state
		return ed
	}

	tests := []struct {
		name               string
		initialBuffer      string
		initialOrigin      string
		key                byte
		expectedMode       editor.Mode
		expectedBuffer     string
		expectedStatusMsg  string // Check contains substring
		expectedFilename   string // Check if set on Enter
		expectedShouldQuit bool   // Check if set on Enter with :wq origin
	}{
		{
			name:              "Esc cancels prompt",
			initialBuffer:     "test",
			initialOrigin:     "w",
			key:               terminal.KeyEsc,
			expectedMode:      editor.ModeNormal,
			expectedBuffer:    "", // Buffer cleared
			expectedStatusMsg: "aborted",
		},
		{
			name:              "Enter with empty buffer aborts",
			initialBuffer:     "",
			initialOrigin:     "w",
			key:               13,
			expectedMode:      editor.ModeNormal,
			expectedBuffer:    "",
			expectedStatusMsg: "aborted",
		},
		{
			name:              "Character appends to buffer",
			initialBuffer:     "file",
			initialOrigin:     "w",
			key:               'n',
			expectedMode:      editor.ModeFileNamePrompt, // Stays in prompt
			expectedBuffer:    "filen",
			expectedStatusMsg: "Save file as: filen",
		},
		{
			name:              "Backspace removes from buffer",
			initialBuffer:     "filen",
			initialOrigin:     "w",
			key:               127, // Backspace
			expectedMode:      editor.ModeFileNamePrompt,
			expectedBuffer:    "file",
			expectedStatusMsg: "Save file as: file",
		},
		{
			name:              "Backspace on empty buffer",
			initialBuffer:     "",
			initialOrigin:     "w",
			key:               127,
			expectedMode:      editor.ModeFileNamePrompt,
			expectedBuffer:    "",
			expectedStatusMsg: "Save file as: ", // Stays as initial prompt
		},
		// --- Enter with filename needs assumptions about SaveFile/QuitEditor ---
		// We assume SaveFile succeeds (returns true, sets IsDirty=false, sets status)
		// We assume QuitEditor respects IsDirty flag
		{
			name:               "Enter sets filename, saves (origin :w)",
			initialBuffer:      "newfile.txt",
			initialOrigin:      "w",
			key:                13,
			expectedMode:       editor.ModeNormal,
			expectedBuffer:     "",
			expectedFilename:   "newfile.txt", // Expect this name to be SET initially
			expectedShouldQuit: false,
			// Status message set by SaveFile, not checked here
		},
		{
			name:               "Enter sets filename, saves and quits (origin :wq)",
			initialBuffer:      "another.txt",
			initialOrigin:      "wq",
			key:                13,
			expectedMode:       editor.ModeNormal,
			expectedBuffer:     "",
			expectedFilename:   "another.txt", // Expect this name to be SET initially
			expectedShouldQuit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditorPrompt(tt.initialBuffer, tt.initialOrigin)
			initialFilename := ed.Filename // Should be empty from helper
			var finalExpectedFilename string

			// For tests involving Enter+filename, prepare temp dir path for SaveFile
			if tt.key == 13 && tt.initialBuffer != "" {
				t.TempDir()
				// Set the expected filename based on the buffer content.
				finalExpectedFilename = tt.initialBuffer
			}

			// processFileNamePrompt is not exported, call ProcessInput
			if ed.CurrentMode != editor.ModeFileNamePrompt {
				t.Fatalf("Test setup error: Editor should start in FileNamePrompt mode")
			}
			ProcessInput(ed, tt.key)

			// Check Mode
			if ed.CurrentMode != tt.expectedMode {
				t.Errorf("Expected Mode %v, got %v", tt.expectedMode, ed.CurrentMode)
			}

			// Check CommandBuffer
			if ed.CommandBuffer != tt.expectedBuffer {
				t.Errorf("Expected CommandBuffer %q, got %q", tt.expectedBuffer, ed.CommandBuffer)
			}

			// Check Status Message (contains)
			if tt.key == 13 && tt.initialBuffer != "" {
				// Don't check status message for save cases, as real SaveFile sets it.
			} else {
				if !strings.Contains(ed.StatusMessage, tt.expectedStatusMsg) {
					t.Errorf("Expected StatusMessage to contain %q, got %q", tt.expectedStatusMsg, ed.StatusMessage)
				}
			}

			// Check Filename (only if Enter was pressed with non-empty buffer)
			if tt.key == 13 && tt.initialBuffer != "" {
				if ed.Filename != finalExpectedFilename {
					t.Errorf("Expected Filename to be %q, got %q", finalExpectedFilename, ed.Filename)
				}
				// Assume SaveFile succeeded (wrote file, set IsDirty=false) and check ShouldQuit
				if ed.IsDirty { // Check if SaveFile correctly marked as not dirty
					t.Error("Expected IsDirty to be false after successful save, but it was true")
				}
				if ed.ShouldQuit != tt.expectedShouldQuit {
					t.Errorf("Expected ShouldQuit to be %t, got %t", tt.expectedShouldQuit, ed.ShouldQuit)
				}
			} else if ed.Filename != initialFilename {
				// Filename should not change unless Enter was pressed with a name
				t.Errorf("Filename unexpectedly changed to %q", ed.Filename)
			}

			// Manual cleanup needed if we let real SaveFile run
			if tt.key == 13 && tt.initialBuffer != "" {
				_ = os.Remove(tt.initialBuffer) // Attempt cleanup
			}
		})
	}
}
