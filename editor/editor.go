package editor

import (
	"strings"
	"time"
)

// Editor holds the state of the text editor
type Editor struct {
	TermWidth           int
	TermHeight          int
	CursorX             int // Cursor position relative to file content (0-based col)
	CursorY             int // Cursor position relative to file content (0-based row)
	RowOffset           int // Top row of the file visible on screen (0-based file index)
	ColOffset           int // Leftmost column of the file visible on screen (0-based file index)
	EditorContent       []string
	CurrentMode         Mode
	Filename            string    // Name of the file being edited
	CommandBuffer       string    // Stores the currently typed command
	StatusMessage       string    // Message to show at the bottom
	StatusMessageTime   time.Time // When the status message was set
	ShouldQuit          bool      // Flag to signal graceful exit
	IsDirty             bool      // Flag for unsaved changes
	PromptOriginCommand string    // Command (:w or :wq) that triggered filename prompt
}

// Mode defines the current state of the editor
type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
	ModeFileNamePrompt // Mode for entering filename on save
)

// NewEditor creates and initializes a new Editor instance.
func NewEditor(width, height int) *Editor {
	return &Editor{
		TermWidth:     width,
		TermHeight:    height,
		CursorX:       0,
		CursorY:       0,
		RowOffset:     0,
		ColOffset:     0,
		EditorContent: []string{""},
		CurrentMode:   ModeNormal,
	}
}

// SetStatusMessage sets the status message and the time it was set.
func (e *Editor) SetStatusMessage(msg string) {
	e.StatusMessage = msg
	e.StatusMessageTime = time.Now()
}

// ensureLineExists appends empty lines if needed to reach target row y.
func (e *Editor) ensureLineExists(y int) {
	for len(e.EditorContent) <= y {
		e.EditorContent = append(e.EditorContent, "")
	}
}

// InsertChar inserts a character at the current cursor position.
func (e *Editor) InsertChar(char byte) {
	e.ensureLineExists(e.CursorY)
	line := e.EditorContent[e.CursorY]
	// TODO: Adjust for colOffset when inserting/deleting
	if e.CursorX >= len(line) {
		line += string(char)
	} else {
		line = line[:e.CursorX] + string(char) + line[e.CursorX:]
	}
	e.EditorContent[e.CursorY] = line
	e.CursorX++
	e.IsDirty = true
}

// InsertNewline inserts a newline by splitting the current line.
func (e *Editor) InsertNewline() {
	e.ensureLineExists(e.CursorY)
	line := e.EditorContent[e.CursorY]
	// TODO: Adjust for colOffset
	beforeCursor := line[:e.CursorX]
	afterCursor := line[e.CursorX:]

	e.EditorContent[e.CursorY] = beforeCursor

	// Insert new line into EditorContent slice
	nextLine := afterCursor
	e.EditorContent = append(e.EditorContent[:e.CursorY+1], append([]string{nextLine}, e.EditorContent[e.CursorY+1:]...)...)

	e.CursorY++
	e.CursorX = 0
	e.IsDirty = true
}

// DeleteChar handles backspace: deleting char or joining lines.
func (e *Editor) DeleteChar() {
	if e.CursorX == 0 && e.CursorY == 0 {
		return
	}

	originalContentLen := len(e.EditorContent)
	originalLineLen := 0
	if e.CursorY < len(e.EditorContent) {
		originalLineLen = len(e.EditorContent[e.CursorY])
	}

	if e.CursorX == 0 { // At start of a line (not the first line)
		// Join with the previous line
		prevLineIndex := e.CursorY - 1
		currentLine := e.EditorContent[e.CursorY]
		prevLine := e.EditorContent[prevLineIndex]
		newCursorX := len(prevLine)
		e.EditorContent[prevLineIndex] = prevLine + currentLine
		// Remove current line
		e.EditorContent = append(e.EditorContent[:e.CursorY], e.EditorContent[e.CursorY+1:]...)
		e.CursorY--
		e.CursorX = newCursorX
	} else {
		e.ensureLineExists(e.CursorY)
		line := e.EditorContent[e.CursorY]
		// TODO: Adjust for colOffset
		if e.CursorX > 0 && e.CursorX <= len(line) {
			line = line[:e.CursorX-1] + line[e.CursorX:]
			e.EditorContent[e.CursorY] = line
			e.CursorX--
		} else if e.CursorX > 0 {
			e.CursorX--
		}
	}

	// Check if content actually changed before marking dirty
	if len(e.EditorContent) != originalContentLen || (e.CursorY < len(e.EditorContent) && len(e.EditorContent[e.CursorY]) != originalLineLen) {
		e.IsDirty = true
	}
}

// EnsureCursorBounds adjusts cursorX if it's beyond the end of the current line
// after a vertical move.
func (e *Editor) EnsureCursorBounds() {
	if e.CursorY >= len(e.EditorContent) {
		// Cursor is on a tilde line (below content)
		e.CursorX = 0
	} else {
		// Cursor is on a content line
		lineLen := len(e.EditorContent[e.CursorY])
		if e.CursorX > lineLen {
			// Snap cursor to the end of the shorter line
			e.CursorX = lineLen
		}
	}
	// Ensure cursor stays within terminal width bounds
	if e.CursorX >= e.TermWidth {
		e.CursorX = e.TermWidth - 1
	}
	if e.CursorX < 0 {
		e.CursorX = 0
	}
}

// LoadFile reads a file into the editorContent buffer.
func (e *Editor) LoadFile(content []byte) {
	// Replace CRLF with LF and split into lines
	fileStr := string(content)
	fileStr = strings.ReplaceAll(fileStr, "\r\n", "\n")
	e.EditorContent = strings.Split(fileStr, "\n")

	// Handle files ending with newline correctly (split adds trailing "")
	if len(e.EditorContent) > 0 && e.EditorContent[len(e.EditorContent)-1] == "" {
		if len(fileStr) > 0 && fileStr != "\n" {
			// File had content and ended with \n, remove the empty string from split
			e.EditorContent = e.EditorContent[:len(e.EditorContent)-1]
		} // Keep the single "" if the file was empty or just "\n"
	}
	// Ensure at least one empty line exists if file was truly empty
	if len(e.EditorContent) == 0 {
		e.EditorContent = []string{""}
	}
	e.IsDirty = false // Loading resets dirty flag
}

// ContentAsString joins the editor content into a single string for saving.
func (e *Editor) ContentAsString() string {
	content := strings.Join(e.EditorContent, "\n")
	// Ensure file ends with a newline
	if len(content) > 0 {
		content += "\n"
	}
	return content
}
