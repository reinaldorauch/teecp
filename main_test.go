package main

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestShouldTryToConnectDefaultPort(t *testing.T) {
	buildProgram()
	client := exec.Command("./teecp.exe", "--client")
	output, err := client.CombinedOutput()

	if err == nil {
		t.Error("must show error when connecting to a port without server")
	}

	if pattern := regexp.MustCompile("6667"); !pattern.Match(output) {
		t.Error("the default port that it tries should be 6667")
	}
}

func TestShouldWaitForAtLeast1Second(t *testing.T) {
	buildProgram()

	client := exec.Command("./teecp.exe", "--client", "--wait-connection")
	output, err := client.CombinedOutput()

	if err == nil {
		t.Error("must show error when connecting to a port without server")
	}

	if pattern := regexp.MustCompile("6667"); !pattern.Match(output) {
		t.Error("the default port that it tries should be 6667")
	}

	outputLines := strings.Split(string(output), "\n")
	if len(outputLines) != 4 {
		t.Error("should have tried 2 times to connect")
	}

	client.ProcessState.UserTime()
}

func buildProgram() {
	buildCommand := exec.Command("go", "build")
	if err := buildCommand.Run(); err != nil {
		log.Fatal("could not build program")
	}
}
