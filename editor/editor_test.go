package editor

import (
	"testing"
)

func TestInsertChar(t *testing.T) {
	// Helper function
	newTestEditor := func(content []string) *Editor {
		ed := NewEditor(80, 24) // Use arbitrary dimensions
		ed.EditorContent = content
		return ed
	}

	tests := []struct {
		name            string
		initialContent  []string
		initialCursorX  int
		initialCursorY  int
		charToInsert    byte
		expectedContent []string
		expectedCursorX int
		expectedIsDirty bool
	}{
		{
			name:            "Insert into empty line",
			initialContent:  []string{""},
			initialCursorX:  0,
			initialCursorY:  0,
			charToInsert:    'a',
			expectedContent: []string{"a"},
			expectedCursorX: 1,
			expectedIsDirty: true,
		},
		{
			name:            "Insert at beginning of line",
			initialContent:  []string{"bc"},
			initialCursorX:  0,
			initialCursorY:  0,
			charToInsert:    'a',
			expectedContent: []string{"abc"},
			expectedCursorX: 1,
			expectedIsDirty: true,
		},
		{
			name:            "Insert in middle of line",
			initialContent:  []string{"ac"},
			initialCursorX:  1,
			initialCursorY:  0,
			charToInsert:    'b',
			expectedContent: []string{"abc"},
			expectedCursorX: 2,
			expectedIsDirty: true,
		},
		{
			name:            "Insert at end of line",
			initialContent:  []string{"ab"},
			initialCursorX:  2,
			initialCursorY:  0,
			charToInsert:    'c',
			expectedContent: []string{"abc"},
			expectedCursorX: 3,
			expectedIsDirty: true,
		},
		{
			name:            "Insert into second line",
			initialContent:  []string{"line1", ""},
			initialCursorX:  0,
			initialCursorY:  1,
			charToInsert:    'X',
			expectedContent: []string{"line1", "X"},
			expectedCursorX: 1,
			expectedIsDirty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditor(tt.initialContent)
			ed.CursorX = tt.initialCursorX
			ed.CursorY = tt.initialCursorY
			ed.IsDirty = false // Reset dirty flag for test

			ed.InsertChar(tt.charToInsert)

			// Check content
			if len(ed.EditorContent) != len(tt.expectedContent) {
				t.Fatalf("Expected %d lines, got %d", len(tt.expectedContent), len(ed.EditorContent))
			}
			for i := range tt.expectedContent {
				if ed.EditorContent[i] != tt.expectedContent[i] {
					t.Errorf("Line %d: Expected %q, got %q", i, tt.expectedContent[i], ed.EditorContent[i])
				}
			}

			// Check cursor X
			if ed.CursorX != tt.expectedCursorX {
				t.Errorf("Expected CursorX %d, got %d", tt.expectedCursorX, ed.CursorX)
			}

			// Check cursor Y (should not change on InsertChar)
			if ed.CursorY != tt.initialCursorY {
				t.Errorf("Expected CursorY %d, got %d", tt.initialCursorY, ed.CursorY)
			}

			// Check dirty flag
			if ed.IsDirty != tt.expectedIsDirty {
				t.Errorf("Expected IsDirty %t, got %t", tt.expectedIsDirty, ed.IsDirty)
			}
		})
	}
}

func TestInsertNewline(t *testing.T) {
	// Helper function
	newTestEditor := func(content []string) *Editor {
		ed := NewEditor(80, 24)
		ed.EditorContent = content
		return ed
	}

	tests := []struct {
		name            string
		initialContent  []string
		initialCursorX  int
		initialCursorY  int
		expectedContent []string
		expectedCursorX int
		expectedCursorY int
		expectedIsDirty bool
	}{
		{
			name:            "Split line in middle",
			initialContent:  []string{"abcd"},
			initialCursorX:  2,
			initialCursorY:  0,
			expectedContent: []string{"ab", "cd"},
			expectedCursorX: 0,
			expectedCursorY: 1,
			expectedIsDirty: true,
		},
		{
			name:            "Insert newline at end of line",
			initialContent:  []string{"abc"},
			initialCursorX:  3,
			initialCursorY:  0,
			expectedContent: []string{"abc", ""},
			expectedCursorX: 0,
			expectedCursorY: 1,
			expectedIsDirty: true,
		},
		{
			name:            "Insert newline at beginning of line",
			initialContent:  []string{"abc"},
			initialCursorX:  0,
			initialCursorY:  0,
			expectedContent: []string{"", "abc"},
			expectedCursorX: 0,
			expectedCursorY: 1,
			expectedIsDirty: true,
		},
		{
			name:            "Insert newline in empty file",
			initialContent:  []string{""},
			initialCursorX:  0,
			initialCursorY:  0,
			expectedContent: []string{"", ""},
			expectedCursorX: 0,
			expectedCursorY: 1,
			expectedIsDirty: true,
		},
		{
			name:            "Insert newline between existing lines",
			initialContent:  []string{"line1", "line2"},
			initialCursorX:  3,
			initialCursorY:  0,
			expectedContent: []string{"lin", "e1", "line2"},
			expectedCursorX: 0,
			expectedCursorY: 1,
			expectedIsDirty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditor(tt.initialContent)
			ed.CursorX = tt.initialCursorX
			ed.CursorY = tt.initialCursorY
			ed.IsDirty = false

			ed.InsertNewline()

			// Check content
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

			// Check cursor X
			if ed.CursorX != tt.expectedCursorX {
				t.Errorf("Expected CursorX %d, got %d", tt.expectedCursorX, ed.CursorX)
			}

			// Check cursor Y
			if ed.CursorY != tt.expectedCursorY {
				t.Errorf("Expected CursorY %d, got %d", tt.expectedCursorY, ed.CursorY)
			}

			// Check dirty flag
			if ed.IsDirty != tt.expectedIsDirty {
				t.Errorf("Expected IsDirty %t, got %t", tt.expectedIsDirty, ed.IsDirty)
			}
		})
	}
}

func TestDeleteChar(t *testing.T) {
	// Helper function
	newTestEditor := func(content []string) *Editor {
		ed := NewEditor(80, 24)
		ed.EditorContent = content
		return ed
	}

	tests := []struct {
		name            string
		initialContent  []string
		initialCursorX  int
		initialCursorY  int
		expectedContent []string
		expectedCursorX int
		expectedCursorY int
		expectedIsDirty bool
	}{
		{
			name:            "Delete from middle of line",
			initialContent:  []string{"abc"},
			initialCursorX:  2,
			initialCursorY:  0,
			expectedContent: []string{"ac"},
			expectedCursorX: 1,
			expectedCursorY: 0,
			expectedIsDirty: true,
		},
		{
			name:            "Delete from end of line",
			initialContent:  []string{"abc"},
			initialCursorX:  3,
			initialCursorY:  0,
			expectedContent: []string{"ab"},
			expectedCursorX: 2,
			expectedCursorY: 0,
			expectedIsDirty: true,
		},
		{
			name:            "Delete - join lines",
			initialContent:  []string{"ab", "cd"},
			initialCursorX:  0,
			initialCursorY:  1,
			expectedContent: []string{"abcd"},
			expectedCursorX: 2, // Cursor moves to end of first line
			expectedCursorY: 0,
			expectedIsDirty: true,
		},
		{
			name:            "Delete at start of file",
			initialContent:  []string{"abc"},
			initialCursorX:  0,
			initialCursorY:  0,
			expectedContent: []string{"abc"}, // Should do nothing
			expectedCursorX: 0,
			expectedCursorY: 0,
			expectedIsDirty: false, // Content didn't change
		},
		{
			name:            "Delete on empty line after line",
			initialContent:  []string{"abc", ""},
			initialCursorX:  0,
			initialCursorY:  1,
			expectedContent: []string{"abc"},
			expectedCursorX: 3,
			expectedCursorY: 0,
			expectedIsDirty: true,
		},
		{
			name:            "Delete only char on line",
			initialContent:  []string{"a"},
			initialCursorX:  1,
			initialCursorY:  0,
			expectedContent: []string{""},
			expectedCursorX: 0,
			expectedCursorY: 0,
			expectedIsDirty: true,
		},
		{
			name:            "Delete on empty file",
			initialContent:  []string{""},
			initialCursorX:  0,
			initialCursorY:  0,
			expectedContent: []string{""},
			expectedCursorX: 0,
			expectedCursorY: 0,
			expectedIsDirty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditor(tt.initialContent)
			ed.CursorX = tt.initialCursorX
			ed.CursorY = tt.initialCursorY
			ed.IsDirty = false // Start clean

			ed.DeleteChar()

			// Check content
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

			// Check cursor X
			if ed.CursorX != tt.expectedCursorX {
				t.Errorf("Expected CursorX %d, got %d", tt.expectedCursorX, ed.CursorX)
			}

			// Check cursor Y
			if ed.CursorY != tt.expectedCursorY {
				t.Errorf("Expected CursorY %d, got %d", tt.expectedCursorY, ed.CursorY)
			}

			// Check dirty flag
			if ed.IsDirty != tt.expectedIsDirty {
				t.Errorf("Expected IsDirty %t, got %t", tt.expectedIsDirty, ed.IsDirty)
			}
		})
	}
}

func TestEnsureCursorBounds(t *testing.T) {
	// Helper function
	newTestEditor := func(content []string, termWidth int) *Editor {
		ed := NewEditor(termWidth, 24) // Height arbitrary
		ed.EditorContent = content
		return ed
	}

	tests := []struct {
		name            string
		initialContent  []string
		termWidth       int
		initialCursorX  int
		initialCursorY  int // Y position to check against
		expectedCursorX int
	}{
		{
			name:            "Cursor stays within shorter line bounds",
			initialContent:  []string{"long line", "short"},
			termWidth:       80,
			initialCursorX:  8, // Initially past end of "short"
			initialCursorY:  1, // Check against line 1 ("short")
			expectedCursorX: 5, // Should snap to length of "short"
		},
		{
			name:            "Cursor stays same on longer line",
			initialContent:  []string{"short", "long line"},
			termWidth:       80,
			initialCursorX:  3, // Initially within "long line"
			initialCursorY:  1, // Check against line 1 ("long line")
			expectedCursorX: 3, // Should stay at 3
		},
		{
			name:            "Cursor goes to 0 when below content",
			initialContent:  []string{"line1"},
			termWidth:       80,
			initialCursorX:  5,
			initialCursorY:  1, // Check against row below content
			expectedCursorX: 0,
		},
		{
			name:            "Cursor clamps to terminal width",
			initialContent:  []string{"a very very long line that exceeds the terminal width"},
			termWidth:       20,
			initialCursorX:  25, // Try to set past termWidth
			initialCursorY:  0,
			expectedCursorX: 19, // Should clamp to termWidth - 1
		},
		{
			name:            "Cursor clamps to 0 if negative",
			initialContent:  []string{"abc"},
			termWidth:       80,
			initialCursorX:  -5, // Try to set negative X
			initialCursorY:  0,
			expectedCursorX: 0, // Should clamp to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditor(tt.initialContent, tt.termWidth)
			// Directly set cursor X/Y *before* calling EnsureCursorBounds
			ed.CursorX = tt.initialCursorX
			ed.CursorY = tt.initialCursorY

			ed.EnsureCursorBounds()

			// Check cursor X
			if ed.CursorX != tt.expectedCursorX {
				t.Errorf("Expected CursorX %d, got %d", tt.expectedCursorX, ed.CursorX)
			}

			// Check cursor Y (should not change)
			if ed.CursorY != tt.initialCursorY {
				t.Errorf("Expected CursorY %d, got %d", tt.initialCursorY, ed.CursorY)
			}
		})
	}
}

func TestLoadFile(t *testing.T) {
	newTestEditor := func() *Editor {
		// Dimensions don't matter much for LoadFile
		ed := NewEditor(10, 5)
		return ed
	}

	tests := []struct {
		name            string
		fileContent     []byte
		expectedLines   []string
		expectedIsDirty bool // Should always be false after load
	}{
		{
			name:            "Empty file content",
			fileContent:     []byte(""),
			expectedLines:   []string{""}, // Should result in one empty line
			expectedIsDirty: false,
		},
		{
			name:            "Single line no newline",
			fileContent:     []byte("hello"),
			expectedLines:   []string{"hello"},
			expectedIsDirty: false,
		},
		{
			name:            "Single line with LF",
			fileContent:     []byte("hello\n"),
			expectedLines:   []string{"hello"},
			expectedIsDirty: false,
		},
		{
			name:            "Multiple lines with LF",
			fileContent:     []byte("line1\nline2\nline3"),
			expectedLines:   []string{"line1", "line2", "line3"},
			expectedIsDirty: false,
		},
		{
			name:            "Multiple lines with CRLF",
			fileContent:     []byte("line1\r\nline2\r\n"),
			expectedLines:   []string{"line1", "line2"},
			expectedIsDirty: false,
		},
		{
			name:            "Multiple lines mixed endings ending LF",
			fileContent:     []byte("line1\r\nline2\nline3\n"),
			expectedLines:   []string{"line1", "line2", "line3"},
			expectedIsDirty: false,
		},
		{
			name:            "Content just newline",
			fileContent:     []byte("\n"),
			expectedLines:   []string{""}, // Should be one empty line
			expectedIsDirty: false,
		},
		{
			name:            "Content just CRLF",
			fileContent:     []byte("\r\n"),
			expectedLines:   []string{""},
			expectedIsDirty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditor()
			ed.IsDirty = true // Set dirty initially to ensure LoadFile resets it

			ed.LoadFile(tt.fileContent)

			// Check content lines
			if len(ed.EditorContent) != len(tt.expectedLines) {
				t.Fatalf("Expected %d lines, got %d", len(tt.expectedLines), len(ed.EditorContent))
			}
			for i := range tt.expectedLines {
				if ed.EditorContent[i] != tt.expectedLines[i] {
					t.Errorf("Line %d: Expected %q, got %q", i, tt.expectedLines[i], ed.EditorContent[i])
				}
			}

			// Check dirty flag
			if ed.IsDirty != tt.expectedIsDirty {
				t.Errorf("Expected IsDirty %t, got %t", tt.expectedIsDirty, ed.IsDirty)
			}
		})
	}
}

func TestContentAsString(t *testing.T) {
	newTestEditor := func(content []string) *Editor {
		// Dimensions don't matter for ContentAsString
		ed := NewEditor(10, 5)
		ed.EditorContent = content
		return ed
	}

	tests := []struct {
		name          string
		editorContent []string
		expectedStr   string
	}{
		{
			name:          "Empty content",
			editorContent: []string{""},
			expectedStr:   "\n", // Should still add trailing newline
		},
		{
			name:          "Single line",
			editorContent: []string{"hello"},
			expectedStr:   "hello\n",
		},
		{
			name:          "Multiple lines",
			editorContent: []string{"line1", "line2", "line3"},
			expectedStr:   "line1\nline2\nline3\n",
		},
		{
			name:          "Lines with empty strings",
			editorContent: []string{"line1", "", "line3"},
			expectedStr:   "line1\n\nline3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newTestEditor(tt.editorContent)
			result := ed.ContentAsString()
			if result != tt.expectedStr {
				t.Errorf("Expected %q, got %q", tt.expectedStr, result)
			}
		})
	}
}
