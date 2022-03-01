package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aymanbagabas/nyancatsh/bubble"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
)

var port = flag.Int("port", 2226, "port to listen on")

func main() {
	flag.Parse()
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("0.0.0.0:%d", *port)),
		wish.WithHostKeyPath(".ssh/nyancatsh"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler()),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Starting SSH server on 0.0.0.0:%d", *port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	<-done
	if err := s.Close(); err != nil {
		log.Fatalln(err)
	}
}

func teaHandler() func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, active := s.Pty()
		if !active {
			s.Write([]byte("not active"))
			s.Exit(0)
			return nil, nil
		}
		w, h := pty.Window.Width, pty.Window.Height
		return bubble.New(w, h), []tea.ProgramOption{tea.WithAltScreen()}
	}
}
