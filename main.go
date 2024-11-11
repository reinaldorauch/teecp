package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/jeffque/teecp/teecp_client"
	"github.com/jeffque/teecp/teecp_server"
)

type AppState = int32
type AppStateDescription struct {
	state          AppState
	description    string
	waitConnection time.Duration
	retryInterval  time.Duration
}

var appTypeStates = struct {
	undefined AppStateDescription
	server    AppStateDescription
	client    AppStateDescription
}{
	AppStateDescription{0, "undefined", 0, 0},
	AppStateDescription{1, "server", 0, 0},
	AppStateDescription{2, "client", 0, time.Duration(1000000000)},
}

func (s AppStateDescription) isServer() bool {
	return s.state != appTypeStates.client.state
}

func defineState(desiredVal AppStateDescription, currAppState *AppStateDescription) func(s string) error {
	return func(s string) error {
		if currAppState.state != appTypeStates.undefined.state {
			return fmt.Errorf("already defined as a [%s], cannot be redefined as a [%s]", currAppState.description, desiredVal.description)
		}
		*currAppState = desiredVal
		return nil
	}
}

func parseDurationOption(s string) (time.Duration, error) {
	matched, err := regexp.MatchString("^\\d*$", s)

	if matched && err == nil {
		s += "s"
	}

	duration, err := time.ParseDuration(s)

	if err != nil {
		return duration, errors.New("invalid duration")
	}

	return duration, nil
}

func setWaitConnectionState(appState *AppStateDescription) func(s string) error {
	return func(s string) error {
		if s == "true" {
			s = "1s"
		}
		duration, err := parseDurationOption(s)

		if err != nil {
			return err
		}

		appState.waitConnection = duration

		return nil
	}
}

func setRetryIntervalState(appState *AppStateDescription) func(s string) error {
	return func(s string) error {
		if s == "true" {
			s = "1s"
		}

		duration, err := parseDurationOption(s)

		if err != nil {
			return err
		}

		appState.retryInterval = duration

		return nil
	}
}

func main() {
	var port int

	serverClientSetted := appTypeStates.undefined

	flag.IntVar(&port, "port", 6667, "A listener port")
	flag.BoolFunc("server", "Define a server teecp instance (conflict with --client)", defineState(appTypeStates.server, &serverClientSetted))
	flag.BoolFunc("wait-connection", "Makes the client wait for a connection retrying until specified (requires --client)", setWaitConnectionState(&serverClientSetted))
	flag.BoolFunc("retry-interval", "Sets the retry time interval for waiting a connection (requires --client and --wait-connection)", setRetryIntervalState(&serverClientSetted))
	flag.BoolFunc("client", "Define a client teecp instance (conflicts with --server)", defineState(appTypeStates.client, &serverClientSetted))
	flag.Parse()

	var err error
	if serverClientSetted.isServer() {
		err = teecp_server.ServerTeecp(port)
	} else {
		err = teecp_client.ListenerTeecp(port, serverClientSetted.waitConnection, serverClientSetted.retryInterval)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
