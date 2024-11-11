package teecp_server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/jeffque/teecp/teecp"
)

func ServerTeecp(port int) error {
	// When creating the teecp.Clients, always have a local client so we can see the echo.
	clients := teecp.Clients{}
	clients.Attach(func(msg string) bool {
		fmt.Print(msg)
		return true
	})

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("could not open socket to port %d: %w", port, err)
	}
	fmt.Printf("Listening at %d", port)

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
