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

	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/logging"
)

const RUN_SSH_SERVER = false
const RUN_TUI = true

const (
	host = "localhost"
	port = "23234"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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

	if RUN_TUI {
		client, _ := NewClient(0, toServer)

		// tuiMain(client)
		newClients <- client
	}

	go gameManager.RunGameServer(ctx)

	if RUN_SSH_SERVER {
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

	sig := <-sigChan
	fmt.Printf("\nReceived signal (%v)\n", sig)

	cancel()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\nDone")
}
