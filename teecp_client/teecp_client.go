package teecp_client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func connectSocket(port int, waitConnection time.Duration, retryInterval time.Duration) (net.Conn, error) {
	var conn net.Conn
	var err error
	start := time.Now()

	if waitConnection > 0 {
		fmt.Fprintf(os.Stderr, "Trying to connect to server for %f seconds\n", waitConnection.Seconds())
	}

	for {
		conn, err = net.Dial("tcp", fmt.Sprintf("localhost:%d", port))

		if waitConnection == 0 || time.Since(start) > waitConnection || waitConnection < retryInterval {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Fprintf(os.Stderr, "Waiting for %f seconds\n", retryInterval.Seconds())
		time.Sleep(retryInterval)
	}

	return conn, err
}

func ListenerTeecp(port int, waitConnection time.Duration, retryInterval time.Duration) error {
	conn, err := connectSocket(port, waitConnection, retryInterval)

	if err != nil {
		return fmt.Errorf("%w", err)
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
