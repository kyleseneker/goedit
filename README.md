# goedit

[![Go Report Card](https://goreportcard.com/badge/github.com/kyleseneker/goedit)](https://goreportcard.com/report/github.com/kyleseneker/goedit)

`goedit` is a simple terminal text editor written in Go. It provides basic modal editing capabilities (Normal, Insert, Command modes) inspired by editors like Vim.

<!-- Consider adding a GIF screencast here showing basic editing, mode switching, and commands -->
<!-- ![goedit demo](link/to/your/demo.gif) -->

## Features

*   **Modal Editing:** Switch between Normal, Insert, and Command modes.
*   **Basic Text Manipulation:** Insert/delete characters, insert newlines.
*   **Vim-like Navigation:** Use `h`, `j`, `k`, `l` or Arrow Keys for cursor movement.
*   **File Operations:**
    *   Open files from the command line.
    *   Save files (`:w`, `:wq`).
    *   Filename prompting on save if needed.
    *   Dirty file indicator (`+`).
*   **Essential Commands:** `:w`, `:wq`, `:q`, `:q!`.
*   **Terminal UI:**
    *   Uses raw mode and alternate screen buffer for clean interaction.
    *   Status bar showing mode, filename, position, and messages.
    *   Basic vertical scrolling.

## Getting Started

### Prerequisites

*   Go (version 1.18 or later recommended)

### Installation & Running

1.  Clone the repository:
    ```sh
    git clone https://github.com/kyleseneker/goedit.git # Replace if needed
    cd goedit
    ```
2.  Build the executable:
    ```sh
    go build
    ```
3.  Run the editor:
    ```sh
    ./goedit [optional_filename]
    ```
    If `optional_filename` is provided, the editor attempts to load it. Otherwise, it starts with an empty buffer.

## Usage

### Modes

*   **Normal Mode:** Default mode for navigation and entering commands.
*   **Insert Mode:** Entered by pressing `i` in Normal Mode. Allows text insertion. Press `Esc` to return to Normal Mode.
*   **Command Mode:** Entered by pressing `:` in Normal Mode. Allows executing commands like `:w` or `:q`. Press `Enter` to execute, `Esc` to cancel.

### Key Bindings

*   **Normal Mode:**
    *   `i`: Enter Insert Mode
    *   `h`, `j`, `k`, `l` / Arrow Keys: Navigate
    *   `: `: Enter Command Mode
*   **Insert Mode:**
    *   `Esc`: Exit to Normal Mode
    *   `Enter`: Insert Newline
    *   `Backspace`: Delete previous character / join lines
    *   Arrow Keys: Navigate
    *   *(Printable Characters)*: Insert text
*   **Command Mode:**
    *   `Enter`: Execute command
    *   `Esc`: Cancel and return to Normal Mode
    *   `Backspace`: Edit command

### Commands

*   `:w`: Write (save) the file. Prompts for filename if needed.
*   `:wq`: Write (save) and quit.
*   `:q`: Quit if the file is not modified.
*   `:q!`: Quit without saving changes (force quit).

## Project Structure

The codebase is organized into several packages:

*   `main`: Entry point, initialization, main loop.
*   `editor`: Core editor state (`Editor` struct) and text manipulation methods.
*   `terminal`: Low-level terminal handling (raw mode, key reading, size).
*   `ui`: Screen rendering logic (drawing text, status bar, cursor).
*   `cmd`: Command mode processing and command implementations.
*   `input`: Normal and Insert mode input handling.

## Known Issues / Future Work

*   No horizontal scrolling (`colOffset` implementation is pending).
*   Limited command set (no search, replace, settings, etc.).
*   No support for advanced features like syntax highlighting, undo/redo, configuration.
*   Basic escape sequence handling in `ReadKey` (might not cover all terminals/keys like Home, End, PgUp/Dn).

## Contributing

This project was primarily a learning exercise. However, contributions, bug reports, or suggestions are welcome! Please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.