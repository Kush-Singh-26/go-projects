# Quick-Note Taker TUI

A fast, distraction-free command-line interface (CLI) scratchpad written in Go. Built using the [Charm](https://charm.sh/) ecosystem to provide a beautiful, full-screen terminal experience.

## Features

An instant digital notebook that lives right in your terminal, featuring:
- **Distraction-Free UI:** Opens in an alternate full-screen buffer, keeping your terminal history clean.
- **Rich Rendering:** Automatically parses and renders your notes as a clean Markdown list.
- **Smart Scrolling:** Features a dynamic viewport that tracks your active cursor natively.
- **Quick Actions:** Easily add new thoughts, delete old ones, or copy a specific note directly to your system clipboard.
- **Persistent Storage:** Auto-saves your session to a local `notes.md` file upon exiting.

## Usage 

You can run the program directly using Go:
```bash
go run main.go
```

**Example Output:**
```text
What's on your mind?

  - Learn the Elm Architecture
  > Build the GoTune Application
  - Master Bubble Tea layouts

╭──────────────────────────────────────────────────╮                            
│ > Type a note ...                                │                            
╰──────────────────────────────────────────────────╯                            

(press esc to quit.)
```

### Build it yourself
To compile the app into a standalone executable that you can use anywhere:
```bash
go build -o bin/quicknote
./bin/quicknote
```

---

## Under the Hood

1. **The Elm Architecture:** Powered by `charmbracelet/bubbletea` to cleanly separate application state (`Model`), event handling (`Update`), and UI rendering (`View`).
2. **Pre-built Components:** Integrates `charmbracelet/bubbles` to handle complex terminal interactions, specifically utilizing the `textinput` for user typing and `viewport` for the scrolling container.
3. **Markdown Rendering:** Passes the raw note strings through `charmbracelet/glamour` to render the text as a formatted Markdown list on the fly.
4. **Declarative Styling:** Uses `charmbracelet/lipgloss` to define the colorful, rounded borders and padding for the input box without messy ANSI codes.
5. **Smart Camera Logic:** Implements custom bounding-box logic to ensure the scrolling viewport perfectly tracks the user's cursor without visual jitter:

```go
func (m *model) ensureCursorVisible() {
	if m.cursor < m.vp.YOffset {
		m.vp.SetYOffset(m.cursor)
	} else if m.cursor >= m.vp.YOffset+m.vp.Height {
		m.vp.SetYOffset(m.cursor - m.vp.Height + 1)
	}
}
```