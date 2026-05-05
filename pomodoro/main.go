package main

import (
	"fmt"
	"os"
	"time"

	progress "charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
)

var initTime = 25 * 60

type model struct {
	secondsLeft   int
	isPaused      bool
	bar           progress.Model
	totalDuration int

	startTime time.Time
	pauseTime time.Time

	sessionQueue []int
	currentIndex int
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

// 🔁 Move to next session
func (m *model) nextSession() tea.Cmd {
	m.currentIndex++

	if m.currentIndex >= len(m.sessionQueue) {
		return tea.Quit
	}

	next := m.sessionQueue[m.currentIndex]

	m.totalDuration = next
	m.secondsLeft = next
	m.startTime = time.Now()
	m.isPaused = false

	return m.bar.SetPercent(0)
}

// ⏮️ Move to previous session
func (m *model) prevSession() tea.Cmd {
	if m.currentIndex == 0 {
		return nil
	}

	m.currentIndex--

	prev := m.sessionQueue[m.currentIndex]

	m.totalDuration = prev
	m.secondsLeft = prev
	m.startTime = time.Now()
	m.isPaused = false

	return m.bar.SetPercent(0)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newBar, barCmd := m.bar.Update(msg)
	m.bar = newBar

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.bar.SetWidth(msg.Width - 4)
		return m, barCmd

	case tickMsg:
		if m.isPaused {
			return m, tea.Batch(tick(), barCmd)
		}

		total := time.Duration(m.totalDuration) * time.Second
		elapsed := time.Since(m.startTime)

		percent := float64(elapsed) / float64(total)
		if percent > 1 {
			percent = 1
		}

		remaining := total - elapsed
		if remaining < time.Second {
			remaining = time.Second
		}

		m.secondsLeft = int(remaining.Seconds())

		percentCmd := m.bar.SetPercent(percent)

		if percent >= 1 {
			cmd := m.nextSession()
			return m, tea.Batch(tick(), cmd, barCmd)
		}

		return m, tea.Batch(tick(), percentCmd, barCmd)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "n":
			cmd := m.nextSession()
			return m, tea.Batch(tick(), cmd, barCmd)

		case "p":
			cmd := m.prevSession()
			return m, tea.Batch(tick(), cmd, barCmd)

		case "space":
			m.isPaused = !m.isPaused
			if m.isPaused {
				m.pauseTime = time.Now()
			} else {
				pausedDuration := time.Since(m.pauseTime)
				m.startTime = m.startTime.Add(pausedDuration)
			}
			return m, barCmd
		}
	}

	return m, barCmd
}

func (m model) currentLabel() string {
	switch m.currentIndex {
	case 0, 2:
		return "Work"
	case 1:
		return "Short Break"
	case 3:
		return "Long Break"
	default:
		return "Session"
	}
}

func (m model) View() tea.View {
	minutes := m.secondsLeft / 60
	seconds := m.secondsLeft % 60

	barView := m.bar.View()

	status := "[ |> Running]"
	if m.isPaused {
		status = "[ || Paused]"
	}

	msg := fmt.Sprintf(
		"\n[%s]\n\n%s\n\nTime remaining: %02d:%02d\n\nPress 'q' to quit | 'n' next | 'p' prev | space pause\n%s\n",
		m.currentLabel(),
		barView,
		minutes,
		seconds,
		status,
	)

	return tea.View{
		Content:   msg,
		AltScreen: true,
	}
}

func main() {
	myBar := progress.New()

	queue := []int{
		25 * 60,
		5 * 60,
		25 * 60,
		15 * 60,
	}

	if len(queue) == 0 {
		fmt.Println("session queue cannot be empty")
		os.Exit(1)
	}

	initialModel := model{
		secondsLeft:   queue[0],
		totalDuration: queue[0],
		bar:           myBar,
		startTime:     time.Now(),
		sessionQueue:  queue,
		currentIndex:  0,
	}

	p := tea.NewProgram(&initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}