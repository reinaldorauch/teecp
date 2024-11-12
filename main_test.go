package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestShouldTryToConnectDefaultPort(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	client := exec.Command("./teecp.exe", "--client", "--wait-connection")
	start := time.Now()
	output, err := client.CombinedOutput()
	runDuration := time.Since(start).Seconds()

	if err == nil {
		t.Error("must show error when connecting to a port without server")
	}

	if pattern := regexp.MustCompile("6667"); !pattern.Match(output) {
		t.Error("the default port that it tries should be 6667")
	}

	outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(outputLines) != 4 {
		t.Error("should have tried 2 times to connect")
	}

	if runDuration < 1.0 || runDuration > 1.05 {
		t.Error("the default wait time should be around 1 second")
	}
}

func TestShouldWaitForDefinedDuration(t *testing.T) {
	t.Parallel()
	var waitSeconds float64 = 2.0
	client := exec.Command("./teecp.exe", "--client", fmt.Sprintf("--wait-connection=%.f", waitSeconds))
	start := time.Now()
	output, err := client.CombinedOutput()
	runDuration := time.Since(start).Seconds()

	if err == nil {
		t.Error("client should show error when not connecting")
	}

	if lines := strings.Split(strings.TrimSpace(string(output)), "\n"); len(lines) != 6 {
		t.Error("should output retry info for 3 times")
	}

	if runDuration < waitSeconds || runDuration > (waitSeconds+0.2) {
		t.Errorf("client should wait for about %.f seconds", waitSeconds)
	}
}

func TestUnderstandTimeUnitsWhenWaiting(t *testing.T) {
	t.Parallel()
	var waitMicroseconds int64 = 1000
	client := exec.Command(
		"./teecp.exe",
		"--client",
		fmt.Sprintf("--wait-connection=%dms", waitMicroseconds),
		"--retry-interval=500ms",
	)
	start := time.Now()
	output, err := client.CombinedOutput()
	runDuration := time.Since(start).Milliseconds()

	if err == nil {
		t.Error("client should show error when not connecting")
	}

	if lines := strings.Split(strings.TrimSpace(string(output)), "\n"); len(lines) != 6 {
		t.Error("should output retry info for 3 times")
	}

	if runDuration < waitMicroseconds || runDuration > (waitMicroseconds+60) {
		t.Errorf("client should wait for about 1 second")
	}
}

func TestSettingRetryIntervalWithoutWaitingShouldHaveNoEffect(t *testing.T) {
	t.Parallel()
	client := exec.Command(
		"./teecp.exe",
		"--client",
		"--retry-interval=500ms",
	)
	output, err := client.CombinedOutput()

	if err == nil {
		t.Error("client should show error when not connecting")
	}

	if lines := strings.Split(strings.TrimSpace(string(output)), "\n"); len(lines) != 1 {
		t.Error("should only have errored 1 time")
	}
	fmt.Print("Done.\n")
}
