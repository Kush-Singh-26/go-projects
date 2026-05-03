package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	clipboard "github.com/atotto/clipboard"
	textinput "github.com/charmbracelet/bubbles/textinput"
	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	glamour "github.com/charmbracelet/glamour"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	appStyle   = lipgloss.NewStyle().Margin(1, 2)
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF7CCB")).
			Padding(0, 1).
			Width(50)
)

type model struct {
	notes   []string
	current textinput.Model
	cursor  int
	vp      viewport.Model
	ready   bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type a note ..."
	ti.Focus()

	var notes []string
	content, err := os.ReadFile("notes.md")
	if err == nil && len(content) > 0 {
		notes = strings.Split(string(content), "\n")
	} else {
		notes = []string{}
	}

	return model{
		notes:   notes,
		current: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) ensureCursorVisible() {
	if m.cursor < m.vp.YOffset {
		m.vp.SetYOffset(m.cursor)
	} else if m.cursor >= m.vp.YOffset+m.vp.Height {
		m.vp.SetYOffset(m.cursor - m.vp.Height + 1)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd, vpCmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "esc", "ctrl+c":
			finalText := strings.Join(m.notes, "\n")
			os.WriteFile("notes.md", []byte(finalText), 0644)
			return m, tea.Quit

		case "enter":
			cv := m.current.Value()
			if strings.TrimSpace(cv) != "" {
				m.notes = append(m.notes, cv)
				m.cursor = len(m.notes) - 1
			}
			m.current.Reset()

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.notes)-1 {
				m.cursor++
			}

		case "ctrl+y":
			if len(m.notes) > 0 {
				clipboard.WriteAll(m.notes[m.cursor])
			}

		case "delete", "ctrl+d":
			if len(m.notes) > 0 {
				m.notes = slices.Delete(m.notes, m.cursor, m.cursor+1)
			}
			if m.cursor >= len(m.notes) && m.cursor > 0 {
				m.cursor--
			}
			if len(m.notes) == 0 {
				m.cursor = 0
			}
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			m.vp = viewport.New(msg.Width, msg.Height-10)
			m.vp.KeyMap.Down.SetEnabled(false)
			m.vp.KeyMap.Up.SetEnabled(false)
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - 10
		}
	}

	// Update components
	m.current, tiCmd = m.current.Update(msg)
	m.vp, vpCmd = m.vp.Update(msg)

	// 🔧 Build and set viewport content HERE (important fix)
	var formattedNotes []string
	for i, note := range m.notes {
		if strings.TrimSpace(note) != "" {
			if i == m.cursor {
				formattedNotes = append(formattedNotes, "> "+note)
			} else {
				formattedNotes = append(formattedNotes, "- "+note)
			}
		}
	}

	notesDisplay := strings.Join(formattedNotes, "\n")

	renderedNotes, err := glamour.Render(notesDisplay, "dark")
	if err != nil {
		renderedNotes = notesDisplay
	}

	m.vp.SetContent(renderedNotes)

	// 🔧 Ensure cursor visibility AFTER content update
	m.ensureCursorVisible()

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	styledInput := inputStyle.Render(m.current.View())

	ui := fmt.Sprintf(
		"What's on your mind?\n\n%s\n\n%s\n\n(press esc to quit.)\n",
		m.vp.View(),
		styledInput,
	)

	return appStyle.Render(ui)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}