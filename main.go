package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/jeffque/teecp/teecp"
)

type appState = int32
type appStateDescription struct {
	state       appState
	description string
}

var appTypeStates = struct {
	undefined appStateDescription
	server    appStateDescription
	client    appStateDescription
}{
	appStateDescription{0, "undefined"},
	appStateDescription{1, "server"},
	appStateDescription{2, "client"},
}

func (s appStateDescription) isServer() bool {
	return s.state != appTypeStates.client.state
}

func defineState(desiredVal appStateDescription, currAppState *appStateDescription) func(s string) error {
	return func(s string) error {
		if currAppState.state != appTypeStates.undefined.state {
			return fmt.Errorf("Already defined as a [%s], cannot be redefined as a [%s]", currAppState.description, desiredVal.description)
		}
		*currAppState = desiredVal
		return nil
	}
}

func main() {
	var port int

	serverClientSetted := appTypeStates.undefined

	flag.IntVar(&port, "port", 6667, "A listener port")
	flag.BoolFunc("server", "Define a server teecp instance (conflict with --client)", defineState(appTypeStates.server, &serverClientSetted))
	flag.BoolFunc("client", "Define a client teecp instance (conflicts with --server)", defineState(appTypeStates.client, &serverClientSetted))
	flag.Parse()

	var err error
	if serverClientSetted.isServer() {
		err = serverTeecp(port)
	} else {
		err = listenerTeecp(port)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func listenerTeecp(port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return fmt.Errorf("Could not open socket to port %d: %w\n", port, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		txt, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error reading stream: %w\nclosing", err)
		}

		// Fprint not strictly needed, but doing so for consistency.
		fmt.Fprint(os.Stdout, txt)
	}

	return nil
}

func serverTeecp(port int) error {
	// When creating the teecp.Clients, always have a local client so we can see the echo.
	clients := teecp.Clients{}
	clients.Attach(func(msg string) bool {
		fmt.Print(msg)
		return true
	})

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("Could not open socket to port %d: %w\n", port, err)
	}
	defer ln.Close()

	// Create a channel so we can signal to the goroutine that it can quit.
	quit := make(chan bool)
	defer close(quit)

	go acceptNewConns(ln, &clients, quit)

	reader := bufio.NewReader(os.Stdin)
	for {
		txt, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error reading form stdin: %w\nclosing teecp", err)
		}
		clients.Broadcast(txt)
	}

	return nil
}

func acceptNewConns(ln net.Listener, clients *teecp.Clients, quit chan bool) {
	// We need the label to break out of the for loop because otherwise we would only break out of the select.
LOOP:
	for {
		select {
		case <-quit:
			// Break out of the loop.
			break LOOP
		default:
			conn, err := ln.Accept()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("tried to connect but failed %s\n", err.Error()))
				return
			}

			// Add the connection as a client.
			clients.Attach(func(msg string) bool {
				if _, err := fmt.Fprint(conn, msg); err != nil {
					conn.Close()
					return false
				}
				return true
			})
		}
	}
}
