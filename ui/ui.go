package ui

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"goedit/editor"
	"goedit/terminal"
)

// RefreshScreen clears the screen, draws the editor content, status bar, and cursor.
func RefreshScreen(e *editor.Editor) {
	var screenBuf bytes.Buffer

	var err error
	e.TermWidth, e.TermHeight, err = terminal.GetSize()
	if err != nil {
		// Keep previous size on error? Log it.
		log.Printf("Error getting terminal size: %v", err)
	}

	screenBuf.WriteString("\x1b[?25l") // Hide cursor
	screenBuf.WriteString("\x1b[H")    // Move cursor home

	// Draw visible portion of the file content
	drawTextRows(e, &screenBuf)

	// Draw Status Bar
	drawStatusBar(e, &screenBuf)

	// Position the actual terminal cursor
	positionCursor(e, &screenBuf)

	screenBuf.WriteString("\x1b[?25h") // Show cursor

	_, err = os.Stdout.Write(screenBuf.Bytes())
	if err != nil {
		log.Printf("Error writing to stdout: %v", err)
	}
}

// drawTextRows draws the visible lines of the file content or tildes.
func drawTextRows(e *editor.Editor, buf *bytes.Buffer) {
	for y := 0; y < e.TermHeight-1; y++ {
		fileRow := e.RowOffset + y

		if fileRow >= len(e.EditorContent) {
			buf.WriteString("~") // Draw tilde
		} else {
			line := e.EditorContent[fileRow]
			// TODO: Handle colOffset for horizontal scrolling
			lineLen := len(line)
			if lineLen > e.TermWidth {
				line = line[:e.TermWidth] // Truncate long lines
			}
			buf.WriteString(line)
		}

		buf.WriteString("\x1b[K") // Clear rest of line
		// Add newline (except for the last text row before status bar)
		if y < e.TermHeight-2 {
			buf.WriteString("\r\n")
		}
	}
}

// drawStatusBar renders the status bar at the bottom line.
func drawStatusBar(e *editor.Editor, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("\x1b[%d;%dH", e.TermHeight, 1))
	buf.WriteString("\x1b[7m")

	// Message content logic
	msg := ""
	if e.CurrentMode == editor.ModeCommand {
		msg = ":" + e.CommandBuffer
	} else if e.CurrentMode == editor.ModeFileNamePrompt {
		msg = e.StatusMessage
	} else if time.Since(e.StatusMessageTime) < 5*time.Second {
		msg = e.StatusMessage
	} else {
		e.StatusMessage = ""
		modeStr := "NORMAL"
		if e.CurrentMode == editor.ModeInsert {
			modeStr = "INSERT"
		}
		fn := e.Filename
		if fn == "" {
			fn = "[No Name]"
		}
		if e.IsDirty {
			fn += " +"
		}
		maxFnLen := 20
		if len(fn) > maxFnLen {
			fn = fn[:maxFnLen-3] + "..."
		}
		leftStatus := fmt.Sprintf(" %s | %s ", modeStr, fn)
		rightStatus := fmt.Sprintf(" %d/%d ", e.CursorY+1, len(e.EditorContent))
		spaces := e.TermWidth - len(leftStatus) - len(rightStatus)
		if spaces < 0 {
			spaces = 0
		}
		msg = leftStatus + strings.Repeat(" ", spaces) + rightStatus
	}

	// Truncate and write message
	if len(msg) > e.TermWidth {
		msg = msg[:e.TermWidth]
	}
	buf.WriteString(msg)
	buf.WriteString(strings.Repeat(" ", e.TermWidth-len(msg))) // Fill rest of bar

	buf.WriteString("\x1b[m") // Reset colors
}

// positionCursor moves the terminal cursor to the calculated screen position.
func positionCursor(e *editor.Editor, buf *bytes.Buffer) {
	// Calculate screen position based on file cursor and viewport offset
	screenCursorY := e.CursorY - e.RowOffset + 1
	screenCursorX := e.CursorX - e.ColOffset + 1

	// Clamp cursor position to valid screen area
	if screenCursorY < 1 {
		screenCursorY = 1
	}
	if screenCursorY >= e.TermHeight {
		screenCursorY = e.TermHeight - 1
	}
	if screenCursorX < 1 {
		screenCursorX = 1
	}
	if screenCursorX > e.TermWidth {
		screenCursorX = e.TermWidth
	}

	cursorPosCmd := fmt.Sprintf("\x1b[%d;%dH", screenCursorY, screenCursorX)
	buf.WriteString(cursorPosCmd)
}
