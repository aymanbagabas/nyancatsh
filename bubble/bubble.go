package bubble

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	termWidth, termHeight := width, height
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
	if maxCol > termWidth {
		minCol = (maxCol - termWidth) / 2
		maxCol = minCol + termWidth
	}
	// Calculate the final animation width
	animWidth := (maxCol - minCol) * len(outputChar)
	b := &Bubble{
		frame:     0,
		width:     width,
		height:    height,
		minRow:    minRow,
		maxRow:    maxRow,
		minCol:    minCol,
		maxCol:    maxCol,
		animWidth: animWidth,
		startTime: time.Now(),
	}
	return b
}

func (b *Bubble) Init() tea.Cmd {
	fmt.Print("\033[H\033[2J\033[?25l")
	return tick
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
	case tickMsg:
		b.frame++
		if b.frame >= len(frames) {
			b.frame = 0
		}
		return b, tick
	}
	return b, tea.Batch(cmds...)
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
			fmt.Fprintf(&s, "\033[48;5;%sm%s", colors[string(char)], outputChar)
		}
		fmt.Fprintln(&s, "\033[0m")
	}
	s.WriteString(b.viewTime())
	fmt.Fprint(&s, "\033[H")
	return s.String()
}

func (b *Bubble) viewTime() string {
	var s strings.Builder
	pr := func(m string) { fmt.Fprintf(&s, "\033[1;37;48;5;17m%s", m) }
	message := fmt.Sprintf("You have nyaned for %.f seconds!", time.Since(b.startTime).Seconds())
	padding := (b.animWidth - (len(message) + 4)) / 2

	pr(strings.Repeat(" ", padding))
	pr(message)
	pr(strings.Repeat(" ", padding+4))
	fmt.Fprintln(&s, "\033[0m")
	return s.String()
}
