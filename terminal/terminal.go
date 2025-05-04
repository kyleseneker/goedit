package terminal

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// Represents special key codes outside the ASCII range.
const (
	KeyNull       byte = 0
	KeyEsc        byte = 27
	KeyArrowUp    byte = 250
	KeyArrowDown  byte = 251
	KeyArrowLeft  byte = 252
	KeyArrowRight byte = 253
	// TODO: Add PageUp, PageDown, Home, End, Del
)

// EnableRawMode puts the terminal into raw mode with VMIN=0, VTIME=1.
// It returns the original terminal state (termios) so it can be restored later.
func EnableRawMode() (*unix.Termios, error) {
	fd := int(os.Stdin.Fd())

	// Get original termios settings
	originalTermios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return nil, fmt.Errorf("error getting initial termios: %w", err)
	}

	// Apply raw mode settings
	_, err = term.MakeRaw(fd)
	if err != nil {
		_ = unix.IoctlSetTermios(fd, unix.TIOCSETA, originalTermios) // Restore on error
		return nil, fmt.Errorf("error setting terminal to raw mode: %w", err)
	}

	// Get the state *after* MakeRaw to explicitly set VMIN/VTIME
	// This ensures our desired settings override any defaults from MakeRaw.
	currentTermios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		_ = unix.IoctlSetTermios(fd, unix.TIOCSETA, originalTermios) // Restore on error
		return nil, fmt.Errorf("error getting termios after MakeRaw: %w", err)
	}

	// Explicitly set VMIN = 0, VTIME = 1 (100ms timeout)
	currentTermios.Cc[unix.VMIN] = 0
	currentTermios.Cc[unix.VTIME] = 1

	// Apply the final state
	err = unix.IoctlSetTermios(fd, unix.TIOCSETA, currentTermios)
	if err != nil {
		_ = unix.IoctlSetTermios(fd, unix.TIOCSETA, originalTermios) // Restore on error
		return nil, fmt.Errorf("error setting final termios (VMIN/VTIME): %w", err)
	}

	fmt.Print("\x1b[?1049h") // Enter alternate screen buffer

	return originalTermios, nil
}

// DisableRawMode restores the terminal to its original termios state.
func DisableRawMode(originalTermios *unix.Termios) {
	fd := int(os.Stdin.Fd())

	fmt.Print("\x1b[?1049l") // Leave alternate screen buffer

	fmt.Print("\x1b[2J\x1b[H")
	if originalTermios != nil {
		if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, originalTermios); err != nil {
			log.Printf("Error restoring terminal state: %v", err)
		}
	}
}

// ReadKey reads a single key press, attempting to handle escape sequences
// using a non-blocking simulation after reading Escape.
func ReadKey() byte {
	var readBuf [3]byte
	n, err := os.Stdin.Read(readBuf[:1]) // Blocking read for first byte
	if err != nil {
		// Handle EOF or other errors gracefully
		return KeyNull
	}
	if n == 0 {
		// This might happen if VMIN=0/VTIME=0, but shouldn't with VTIME=1
		return KeyNull
	}

	key := readBuf[0]

	if key == KeyEsc {
		// Attempt to read the rest of the sequence (non-blocking due to VMIN=0, VTIME=1)
		n_seq, err_seq := os.Stdin.Read(readBuf[1:])

		if err_seq != nil {
			return KeyEsc // Assume it was just Esc
		}

		if n_seq >= 2 && readBuf[1] == '[' {
			switch readBuf[2] {
			case 'A':
				return KeyArrowUp
			case 'B':
				return KeyArrowDown
			case 'C':
				return KeyArrowRight
			case 'D':
				return KeyArrowLeft
				// TODO: Add Home, End, PgUp, PgDn, Del sequences
			}
			return KeyEsc
		} else if n_seq >= 1 { // Esc followed by something else (e.g., Alt+key?)
			return KeyEsc
		} else {
			return KeyEsc
		}
	}

	return key
}

// GetSize returns the current width and height of the terminal.
func GetSize() (width, height int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}
