package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Runs the actual Git Bash/Unix host entrypoint, including native argument conversion.
func TestWorkshopWrapper(t *testing.T) {
	wrapper := os.Getenv("WORKSHOP_TEST_WRAPPER")
	if wrapper == "" {
		t.Skip("set WORKSHOP_TEST_WRAPPER to test the installed host entrypoint")
	}
	cmd := exec.Command("bash", filepath.ToSlash(wrapper), "--allow-notes")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	client := mcp.NewClient(&mcp.Implementation{Name: "wrapper-test", Version: "1"}, nil)
	session, err := client.Connect(ctx(t), &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("wrapper failed: %v\n%s", err, stderr.String())
	}
	defer session.Close()
	result, err := session.CallTool(ctx(t), &mcp.CallToolParams{Name: "get_task", Arguments: map[string]any{"id": "task-101"}})
	if err != nil || result.IsError {
		t.Fatalf("wrapper read: result=%v error=%v", result, err)
	}
	note := "Native wrapper verified: café, 日本語, \"quotes\".\nSecond line."
	result, err = session.CallTool(ctx(t), &mcp.CallToolParams{Name: "add_task_note", Arguments: map[string]any{"id": "task-101", "note": note}})
	if err != nil || result.IsError {
		t.Fatalf("wrapper write: result=%v error=%v", result, err)
	}
	var got struct {
		Body string `json:"body"`
	}
	decodeStructured(t, result, &got)
	if !strings.Contains(got.Body, note) {
		t.Fatalf("note changed in transit: %q", got.Body)
	}
}
