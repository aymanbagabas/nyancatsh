package bubble

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg struct{}

var (
	interval = 90 * time.Millisecond
	// Set output character
	outputChar = "  "
	// Set colors
	colors = map[string]string{
		"+": "226",
		"@": "223",
		",": "17",
		"-": "205",
		"#": "82",
		".": "15",
		"$": "219",
		"%": "217",
		";": "99",
		"&": "214",
		"=": "39",
		"'": "0",
		">": "196",
		"*": "245",
	}
)

type Bubble struct {
	frame     int
	width     int
	height    int
	minRow    int
	maxRow    int
	minCol    int
	maxCol    int
	animWidth int
	startTime time.Time
}

func New(width, height int) *Bubble {
	b := &Bubble{
		frame:     0,
		startTime: time.Now(),
	}
	b.setSize(width, height)
	return b
}

func (b *Bubble) setSize(w, h int) {
	termWidth, termHeight := w, h
	// Calculate the width in terms of the output char
	termWidth = termWidth / len(outputChar)
	minRow := 0
	maxRow := len(frames[0])
	minCol := 0
	maxCol := len(frames[0][0])
	if maxRow > termHeight {
		minRow = (maxRow - termHeight) / 2
		maxRow = minRow + termHeight
	}
	if maxRow > 0 {
		// account for timer text line
		maxRow -= 1
	}
	if maxCol > termWidth {
		minCol = (maxCol - termWidth) / 2
		maxCol = minCol + termWidth
	}
	// Calculate the final animation width
	animWidth := (maxCol - minCol) * len(outputChar)
	b.width = w
	b.height = h
	b.minRow = minRow
	b.maxRow = maxRow
	b.minCol = minCol
	b.maxCol = maxCol
	b.animWidth = animWidth
}

func (b *Bubble) Init() tea.Cmd {
	return tick
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		b.setSize(msg.Width, msg.Height)
	case tickMsg:
		b.frame++
		if b.frame >= len(frames) {
			b.frame = 0
		}
		return b, tick
	}
	return b, nil
}

func tick() tea.Msg {
	time.Sleep(interval)
	return tickMsg{}
}

func (b *Bubble) View() string {
	var s strings.Builder
	frame := frames[b.frame]
	for _, line := range frame[b.minRow:b.maxRow] {
		for _, char := range line[b.minCol:b.maxCol] {
			style := lipgloss.NewStyle().Background(lipgloss.Color(colors[string(char)]))
			s.WriteString(style.Render(outputChar))
		}
		s.WriteString("\n")
	}
	s.WriteString(b.viewTime())
	return s.String()
}

func (b *Bubble) viewTime() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color(colors[","])).
		Foreground(lipgloss.Color(colors["."])).
		Bold(true).
		Width(b.animWidth).
		Align(lipgloss.Center)
	message := fmt.Sprintf("You have nyaned for %.f seconds!", time.Since(b.startTime).Seconds())
	return style.Render(message)
}
