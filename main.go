package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/logging"
)

const RUN_SERVER = false

const (
	host = "localhost"
	port = "23234"
)

func main() {
	// TODO: Start goroutine for server
	// Create a new game manager (struct)
	//      GM has things like players, ids, etc
	// The GM has a channel for receiving new connections
	// GM = newGM
	// go GM.run
	//
	// Every time a new player sshes in, they get assigned a Client
	// with the ToServer chan and a new ToClient chan
	// via the middleware handler and the Client gets sent to
	// the newClients chan to the GM
	//
	// Main owns all ToServer chans and creates new ones for ToClients

	newClients := make(chan Client, 8)
	toServer := make(chan ClientMessage, 8)
	gameManager := NewGameManager(newClients, toServer)

	go gameManager.RunGameServer()

	fmt.Println(gameManager) // TODO: Use game manager
	if RUN_SERVER {
		s, err := wish.NewServer(
			wish.WithAddress(net.JoinHostPort(host, port)),
			wish.WithHostKeyPath(".ssh/id_ed25519"),
			ssh.AllocatePty(),
			wish.WithMiddleware(
				// Bubble Tea TUI middleware
				// TODO: Figure out how to turn a TUI into a middleware
				// Then put the middleware here
				// TODO: Figure out how the TUI communicates with the server
				// The handler needs to request a Client
				func(next ssh.Handler) ssh.Handler {
					return func(sess ssh.Session) {
						toClient := make(chan ServerMessage, 8)
						client := Client{
							ID: 0, // ??? TODO: How to get ID?
							// Idea: start with id -1
							// Then have a server message to assign an ID
							// ID -1 will not be recognized later (?)
							ToServer: toServer,
							ToClient: toClient,
						}
						newClients <- client
						pty, _, _ := sess.Pty()
						wish.Printf(sess, "%s\r\n", pty.Term)
						wish.Printf(sess, "Hello World!\r\n")

						next(sess)
					}
				},

				logging.Middleware(),
			),
		)
		if err != nil {
			log.Error("Could not start server", "error", err)
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Info("Starting SSH server", "host", host, "port", port)
		go func() {
			if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
				log.Error("Could not start server", "error", err)
				done <- nil
			}
		}()

		<-done
		log.Info("Stopping SSH server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer func() { cancel() }()
		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not stop server", "error", err)
		}
	}
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := model{
		term:      pty.Term,
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		txtStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		quitStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		bg:        "light",
	}
	return m, []tea.ProgramOption{}
}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term      string
	profile   string
	width     int
	height    int
	bg        string
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	// default values
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.ColorProfileMsg:
		m.profile = msg.String()
	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			m.bg = "dark"
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	s := fmt.Sprintf("Your term is %s\nYour window size is %dx%d\nBackground: %s\nColor Profile: %s", m.term, m.width, m.height, m.bg, m.profile)
	v := tea.NewView(m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render("Press 'q' to quit\n"))
	v.AltScreen = true
	return v
}
