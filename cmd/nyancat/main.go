package main

import (
	"log"
	"os"

	"github.com/aymanbagabas/nyancatsh/bubble"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func main() {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	debug := os.Getenv("DEBUG")
	if debug != "" {
		f, err := tea.LogToFile(debug, "")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
	b := bubble.New(w, h)
	p := tea.NewProgram(b, tea.WithAltScreen())
	if _, err = p.Run(); err != nil {
		log.Fatal(err)
	}
}
