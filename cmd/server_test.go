package cmd

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestServerCmd_HasPortFlag(t *testing.T) {
	flag := serverCmd.Flags().Lookup("port")
	if flag == nil {
		t.Fatal("Expected 'port' flag to be defined")
	}
	if flag.DefValue != "8080" {
		t.Errorf("Expected default port to be 8080, got %s", flag.DefValue)
	}
}

func TestServerCmd_Run_StartsServer(t *testing.T) {
	// Find a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to get a free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()

	serverPort = port

	// Run server in a goroutine
	done := make(chan struct{})
	go func() {
		// Use a copy of the command to avoid side effects
		cmd := &cobra.Command{}
		*cmd = *serverCmd
		cmd.SetArgs([]string{})
		// Run the command (it will block, so we need to stop it)
		go func() {
			defer close(done)
			cmd.Run(cmd, []string{})
		}()
	}()

	// Give the server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Try to connect to the server
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", // check port
		strings.TrimPrefix(serverCmd.Flags().Lookup("port").Value.String(), ":")))
	if err != nil {
		t.Errorf("Expected server to be listening, but got error: %v", err)
	} else {
		conn.Close()
	}

	// The server blocks forever, so we can't cleanly stop it in this test.
	// In a real-world scenario, refactor to allow graceful shutdown for testing.
}
