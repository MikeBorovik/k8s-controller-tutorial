package cmd

import (
	"testing"
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

func TestGetServerKubeClient_InvalidPath(t *testing.T) {
	_, err := getServerKubeClient("/invalid/path", false)
	if err == nil {
		t.Error("expected error for invalid kubeconfig path")
	}
}
