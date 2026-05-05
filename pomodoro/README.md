# Pomodoro TUI

A clean, keyboard-driven Pomodoro timer built in Go. Designed as a real-time terminal application using the [Charm](https://charm.sh/) ecosystem, with a focus on precise timing, smooth state transitions, and queue-based session control.

---

## Features

A minimal yet powerful Pomodoro engine right in your terminal:

- **Full-Screen TUI:** Runs in an alternate screen buffer for a distraction-free experience.
- **Session Queue Engine:** Automatically cycles through work and break sessions (e.g., 25m → 5m → 25m → 15m).
- **Real-Time Progress Bar:** Smoothly updating progress visualization using a high-frequency tick loop.
- **Playback Controls:**  
  - `n` → skip to next session  
  - `p` → go to previous session  
  - `space` → pause/resume  
- **Accurate Time Tracking:** Handles pause/resume without drifting using real elapsed time.
- **State Machine Design:** Built around deterministic state transitions (not naive countdown loops).

---

## Usage 

Run the program directly using Go:

```bash
go run main.go
```

**Example Output:**
```text
[Work]

████████████████████░░░░░░░░░░░░░░░░

Time remaining: 12:43

Press 'q' to quit | 'n' next | 'p' prev | space pause
[ |> Running]
```

---

### Build it yourself

To compile the app into a standalone executable:

```bash
go build -o bin/pomodoro
./bin/pomodoro
```

---

## Under the Hood

1. **The Elm Architecture:**  
   Built using `charmbracelet/bubbletea`, separating:
   - `Model` → application state  
   - `Update` → event handling  
   - `View` → rendering  

2. **Progress Visualization:**  
   Uses `charmbracelet/bubbles/progress` to render a dynamic progress bar, updated via commands (`SetPercent`) rather than direct mutation.

3. **Time Engine (Core Logic):**  
   Instead of decrementing counters, the app calculates time using:

```go
elapsed := time.Since(m.startTime)
percent := float64(elapsed) / float64(totalDuration)
```

This ensures:
- No drift
- Accurate pause/resume
- Smooth UI updates

4. **Session Queue System:**  
   The Pomodoro cycle is implemented as a queue:

```go
sessionQueue := []int{
	25 * 60,
	5 * 60,
	25 * 60,
	15 * 60,
}
```

Progression is handled through a state transition function:

```go
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
```

5. **Tick-Driven Updates:**  
   A continuous tick loop drives the UI:

```go
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
```

This enables:
- Smooth animations
- Responsive controls
- Real-time progress updates

