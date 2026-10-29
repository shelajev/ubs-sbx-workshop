// Protocol tests: a real MCP client speaks to a real beans-mcp subprocess over stdio,
// backed by a real (disposable) Beans backlog and the pinned Beans CLI.
//
// These are not unit tests with a fake backlog. They exercise the same code path the
// SBX gateway uses: process launch, initialization, capability negotiation, tool
// listing, tool calls, typed errors, and shutdown.
package main_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	binaryPath string
	beansPath  string
)

// TestMain builds the adapter and locates the pinned Beans CLI. Without Beans the
// protocol tests are skipped rather than silently passing against nothing.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "beans mcp test *")
	if err != nil {
		panic(err)
	}

	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	binaryPath = os.Getenv("WORKSHOP_TEST_ADAPTER_BIN")
	if binaryPath == "" {
		binaryPath = filepath.Join(dir, "beans-mcp"+suffix)
		build := exec.Command("go", "build", "-o", binaryPath, ".")
		build.Stderr = os.Stderr
		if err := build.Run(); err != nil {
			panic("failed to build beans-mcp: " + err.Error())
		}

	}
	beansPath = os.Getenv("WORKSHOP_TEST_BEANS_BIN")
	if beansPath == "" {
		// Tests run from mcp/beans/cmd/beans-mcp.
		if wd, err := os.Getwd(); err == nil {
			candidate := filepath.Join(wd, "..", "..", "..", "..", ".local", "chapters", "bin", "beans"+suffix)
			if abs, err := filepath.Abs(candidate); err == nil {
				if _, err := os.Stat(abs); err == nil {
					beansPath = abs
				}
			}
		}
	}
	if beansPath == "" {
		if p, err := exec.LookPath("beans"); err == nil {
			beansPath = p
		}
	}
	if beansPath == "" && os.Getenv("WORKSHOP_REQUIRE_PROTOCOL_TESTS") == "1" {
		panic("real Beans CLI required: set WORKSHOP_TEST_BEANS_BIN")
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

// backlog creates a disposable Beans backlog with the given tasks and returns its
// config and data paths.
func backlog(t *testing.T, tasks map[string]string) (configFile, dataDir string) {
	t.Helper()
	if beansPath == "" {
		t.Skip("pinned beans CLI not found; set WORKSHOP_TEST_BEANS_BIN")
	}
	root := filepath.Join(t.TempDir(), "workshop backlog with spaces")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	dataDir = filepath.Join(root, ".beans")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configFile = filepath.Join(root, ".beans.yml")
	config := "beans:\n    path: .beans\n    prefix: task-\n    id_length: 4\n" +
		"    default_status: todo\n    default_type: task\n"
	if err := os.WriteFile(configFile, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	for id, body := range tasks {
		name := id + "--" + id + ".md"
		if err := os.WriteFile(filepath.Join(dataDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return configFile, dataDir
}

func task(id, title, status, body string) string {
	return "---\n# " + id + "\ntitle: '" + title + "'\nstatus: " + status +
		"\ntype: task\npriority: normal\ntags:\n    - workshop\n" +
		"created_at: 2026-09-16T09:00:00Z\nupdated_at: 2026-09-16T09:00:00Z\n---\n\n" + body + "\n"
}

// connect starts a beans-mcp subprocess and returns an initialized client session.
func connect(t *testing.T, configFile, dataDir string, extraArgs ...string) (*mcp.ClientSession, *strings.Builder) {
	t.Helper()
	args := append([]string{
		"--beans-bin=" + beansPath,
		"--beans-config=" + configFile,
		"--beans-data=" + dataDir,
	}, extraArgs...)

	cmd := exec.Command(binaryPath, args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr

	client := mcp.NewClient(&mcp.Implementation{Name: "protocol-test", Version: "1.0.0"}, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect failed: %v\nstderr:\n%s", err, stderr.String())
	}
	t.Cleanup(func() { _ = session.Close() })
	return session, &stderr
}

func ctx(t *testing.T) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return c
}

func sampleTasks() map[string]string {
	return map[string]string{
		"task-101": task("task-101", "Warm-up: active-filter result count", "todo",
			"Fix the count text.\n\n## Acceptance criteria\n- singular for one result\n- plural otherwise\n\n## Out of scope\n- filtering"),
		"task-102": task("task-102", "Add assignment and resolution", "in-progress",
			"Implement assignment.\n\n## Acceptance criteria\n- note required at the API"),
		"task-109": task("task-109", "Completed example", "completed", "Nothing to do."),
	}
}

func TestInitializeAndCapabilityNegotiation(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, _ := connect(t, cfg, data)

	init := session.InitializeResult()
	if init.ServerInfo.Name != "workshop-beans" {
		t.Errorf("server name = %q, want workshop-beans", init.ServerInfo.Name)
	}
	if init.ProtocolVersion == "" {
		t.Error("no protocol version was negotiated")
	}
	if init.Capabilities == nil || init.Capabilities.Tools == nil {
		t.Error("server did not advertise the tools capability")
	}
	if !strings.Contains(init.Instructions, "get_task") {
		t.Errorf("instructions do not mention get_task: %q", init.Instructions)
	}
	if init.ServerInfo.Version == "" {
		t.Error("server did not report a version")
	}
	t.Logf("negotiated protocol %s with %s %s", init.ProtocolVersion, init.ServerInfo.Name, init.ServerInfo.Version)
}

func TestToolListIsReadOnlyByDefault(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, _ := connect(t, cfg, data)

	tools, err := session.ListTools(ctx(t), nil)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]*mcp.Tool{}
	for _, tool := range tools.Tools {
		names[tool.Name] = tool
	}
	if len(names) != 2 {
		t.Errorf("tool count = %d, want 2 (get_task, list_tasks): %v", len(names), keys(names))
	}
	for _, want := range []string{"get_task", "list_tasks"} {
		tool, ok := names[want]
		if !ok {
			t.Fatalf("missing tool %s; got %v", want, keys(names))
		}
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Errorf("%s is not annotated read-only", want)
		}
		if tool.InputSchema == nil {
			t.Errorf("%s has no input schema", want)
		}
	}
	if _, present := names["add_task_note"]; present {
		t.Error("the presenter write tool must not be registered by default")
	}
}

func TestGetTaskReturnsTheHostTaskWithAVersion(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, _ := connect(t, cfg, data)

	result, err := session.CallTool(ctx(t), &mcp.CallToolParams{
		Name: "get_task", Arguments: map[string]any{"id": "task-101"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("get_task returned an error: %s", textOf(result))
	}

	var got struct {
		ID                 string   `json:"id"`
		Title              string   `json:"title"`
		Body               string   `json:"body"`
		Status             string   `json:"status"`
		Version            string   `json:"version"`
		AcceptanceCriteria []string `json:"acceptance_criteria"`
		Backlog            string   `json:"backlog"`
	}
	decodeStructured(t, result, &got)

	if got.ID != "task-101" {
		t.Errorf("id = %q, want task-101", got.ID)
	}
	if !strings.Contains(got.Title, "active-filter result count") {
		t.Errorf("title = %q", got.Title)
	}
	if !strings.Contains(got.Body, "Fix the count text") {
		t.Errorf("body does not contain the task body: %q", got.Body)
	}
	if got.Status != "todo" {
		t.Errorf("status = %q, want todo", got.Status)
	}
	if got.Version == "" {
		t.Error("no content version (etag) was returned")
	}
	if len(got.AcceptanceCriteria) != 2 {
		t.Errorf("acceptance criteria = %v, want 2 entries", got.AcceptanceCriteria)
	}
	if got.Backlog != "workshop" {
		t.Errorf("backlog label = %q, want workshop", got.Backlog)
	}
	if strings.Contains(textOf(result), data) {
		t.Error("the text response leaked the host data directory path")
	}
}

func TestGetTaskRejectsMalformedIDs(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, _ := connect(t, cfg, data)

	cases := []struct{ name, id, wantCode string }{
		{"empty", "", "invalid_id"},
		{"flag injection", "--beans-path", "invalid_id"},
		{"path traversal", "../../etc/passwd", "invalid_id"},
		{"shell metacharacters", "task-101; rm -rf /", "invalid_id"},
		{"command substitution", "$(whoami)", "invalid_id"},
		{"uppercase", "TASK-101", "invalid_id"},
		{"single group", "task", "invalid_id"},
		{"too long", strings.Repeat("a", 30) + "-" + strings.Repeat("b", 40), "invalid_id"},
		{"unknown but well formed", "task-999", "task_not_found"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(ctx(t), &mcp.CallToolParams{
				Name: "get_task", Arguments: map[string]any{"id": tc.id},
			})
			if err != nil {
				t.Fatalf("transport error (the server should answer with a tool error): %v", err)
			}
			if !result.IsError {
				t.Fatalf("id %q was accepted; response: %s", tc.id, textOf(result))
			}
			code := errorCode(t, result)
			if code != tc.wantCode {
				t.Errorf("error code = %q, want %q (text: %s)", code, tc.wantCode, textOf(result))
			}
		})
	}
}

func TestListTasksIsBoundedAndFiltered(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, _ := connect(t, cfg, data)

	type listResult struct {
		Tasks []struct {
			ID      string `json:"id"`
			Status  string `json:"status"`
			Version string `json:"version"`
		} `json:"tasks"`
		Count     int  `json:"count"`
		Limit     int  `json:"limit"`
		Truncated bool `json:"truncated"`
	}

	t.Run("default limit", func(t *testing.T) {
		result, err := session.CallTool(ctx(t), &mcp.CallToolParams{Name: "list_tasks"})
		if err != nil {
			t.Fatal(err)
		}
		var got listResult
		decodeStructured(t, result, &got)
		if got.Count != 3 {
			t.Errorf("count = %d, want 3", got.Count)
		}
		if got.Limit != 20 {
			t.Errorf("default limit = %d, want 20", got.Limit)
		}
		for _, task := range got.Tasks {
			if task.Version == "" {
				t.Errorf("task %s has no version", task.ID)
			}
		}
	})

	t.Run("status filter", func(t *testing.T) {
		result, err := session.CallTool(ctx(t), &mcp.CallToolParams{
			Name: "list_tasks", Arguments: map[string]any{"status": "completed"},
		})
		if err != nil {
			t.Fatal(err)
		}
		var got listResult
		decodeStructured(t, result, &got)
		if got.Count != 1 || got.Tasks[0].ID != "task-109" {
			t.Errorf("status filter returned %+v", got.Tasks)
		}
	})

	t.Run("limit truncates", func(t *testing.T) {
		result, err := session.CallTool(ctx(t), &mcp.CallToolParams{
			Name: "list_tasks", Arguments: map[string]any{"limit": 1},
		})
		if err != nil {
			t.Fatal(err)
		}
		var got listResult
		decodeStructured(t, result, &got)
		if got.Count != 1 || !got.Truncated {
			t.Errorf("limit 1 gave count=%d truncated=%v", got.Count, got.Truncated)
		}
	})

	t.Run("invalid filters are rejected", func(t *testing.T) {
		for name, args := range map[string]map[string]any{
			"unknown status":      {"status": "sideways"},
			"status flag":         {"status": "--beans-path"},
			"tag with spaces":     {"tag": "a b"},
			"limit above maximum": {"limit": 5000},
			"negative limit":      {"limit": -1},
		} {
			result, err := session.CallTool(ctx(t), &mcp.CallToolParams{Name: "list_tasks", Arguments: args})
			if err != nil {
				t.Fatalf("%s: transport error: %v", name, err)
			}
			if !result.IsError {
				t.Errorf("%s: filter %v was accepted", name, args)
				continue
			}
			if code := errorCode(t, result); code != "invalid_filter" {
				t.Errorf("%s: error code = %q, want invalid_filter", name, code)
			}
		}
	})
}

func TestMultipleReadsSeeHostSideChanges(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, _ := connect(t, cfg, data)

	first, err := session.CallTool(ctx(t), &mcp.CallToolParams{
		Name: "get_task", Arguments: map[string]any{"id": "task-101"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var before struct {
		Version string `json:"version"`
		Body    string `json:"body"`
	}
	decodeStructured(t, first, &before)

	// Host-side edit, exactly as a maintainer changing the task would do it.
	path := filepath.Join(data, "task-101--task-101.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(content), "Fix the count text.",
		"Fix the count text. UPDATED BY THE HOST.", 1)
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}

	second, err := session.CallTool(ctx(t), &mcp.CallToolParams{
		Name: "get_task", Arguments: map[string]any{"id": "task-101"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var after struct {
		Version string `json:"version"`
		Body    string `json:"body"`
	}
	decodeStructured(t, second, &after)

	if !strings.Contains(after.Body, "UPDATED BY THE HOST") {
		t.Error("the second read did not see the host-side change")
	}
	if before.Version == after.Version {
		t.Errorf("the content version did not change: %s", after.Version)
	}
	t.Logf("version changed from %s to %s", before.Version, after.Version)
}

func TestTwoIndependentProcessesShareTheBacklog(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	first, _ := connect(t, cfg, data)
	second, _ := connect(t, cfg, data)

	for i, session := range []*mcp.ClientSession{first, second} {
		result, err := session.CallTool(ctx(t), &mcp.CallToolParams{
			Name: "get_task", Arguments: map[string]any{"id": "task-102"},
		})
		if err != nil {
			t.Fatalf("session %d: %v", i, err)
		}
		if result.IsError {
			t.Fatalf("session %d: %s", i, textOf(result))
		}
		var got struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		decodeStructured(t, result, &got)
		if got.ID != "task-102" || got.Status != "in-progress" {
			t.Errorf("session %d read %+v", i, got)
		}
	}
}

func TestStdoutCarriesProtocolTrafficOnly(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, stderr := connect(t, cfg, data, "--log-level=debug")

	for _, id := range []string{"task-101", "task-999", "task-102"} {
		if _, err := session.CallTool(ctx(t), &mcp.CallToolParams{
			Name: "get_task", Arguments: map[string]any{"id": id},
		}); err != nil {
			t.Fatalf("call for %s failed, which would also happen if diagnostics corrupted stdout: %v", id, err)
		}
	}
	// If any diagnostic had gone to stdout the JSON-RPC stream above would have failed
	// to parse. Confirm the diagnostics were produced at all, on stderr.
	if !strings.Contains(stderr.String(), "workshop-beans starting") {
		t.Errorf("expected startup diagnostics on stderr, got:\n%s", stderr.String())
	}
	if !strings.Contains(stderr.String(), "get_task served") {
		t.Errorf("expected per-call diagnostics on stderr, got:\n%s", stderr.String())
	}
}

func TestWriteToolIsRefusedWithoutADisposableBacklog(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	session, _ := connect(t, cfg, data, "--enable-presenter-note-tool")

	tools, err := session.ListTools(ctx(t), nil)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, tool := range tools.Tools {
		if tool.Name == "add_task_note" {
			found = true
		}
	}
	if !found {
		t.Fatal("add_task_note should be registered when explicitly enabled")
	}

	result, err := session.CallTool(ctx(t), &mcp.CallToolParams{
		Name: "add_task_note", Arguments: map[string]any{"id": "task-101", "note": "should not be written"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("the write was accepted against a backlog with no disposable marker")
	}
	if code := errorCode(t, result); code != "write_disabled" {
		t.Errorf("error code = %q, want write_disabled", code)
	}

	// The task body must be untouched.
	content, err := os.ReadFile(filepath.Join(data, "task-101--task-101.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "should not be written") {
		t.Error("the refused note was written to the task anyway")
	}
}

func TestWriteToolWorksOnlyOnADisposableBacklog(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	if err := os.WriteFile(filepath.Join(data, ".workshop-disposable"), []byte("disposable\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	session, _ := connect(t, cfg, data, "--enable-presenter-note-tool")

	result, err := session.CallTool(ctx(t), &mcp.CallToolParams{
		Name: "add_task_note", Arguments: map[string]any{"id": "task-101", "note": "Review complete: \"quoted\" paths, café, 日本語.\nSecond line."},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("the write was refused on a disposable backlog: %s", textOf(result))
	}
	var got struct {
		Body string `json:"body"`
	}
	decodeStructured(t, result, &got)
	if !strings.Contains(got.Body, "Review complete: \"quoted\" paths, café, 日本語.\nSecond line.") {
		t.Errorf("the note is missing from the returned body: %q", got.Body)
	}
}

func TestStartupFailsLoudlyOnMisconfiguration(t *testing.T) {
	if beansPath == "" {
		t.Skip("pinned beans CLI not found")
	}
	cfg, data := backlog(t, sampleTasks())

	cases := []struct {
		name string
		args []string
		want string
	}{
		{"missing beans binary", []string{"--beans-bin=/nonexistent/beans", "--beans-config=" + cfg, "--beans-data=" + data}, "beans executable was not found"},
		{"missing config", []string{"--beans-bin=" + beansPath, "--beans-config=/nonexistent/.beans.yml", "--beans-data=" + data}, "beans config file was not found"},
		{"missing data dir", []string{"--beans-bin=" + beansPath, "--beans-config=" + cfg, "--beans-data=/nonexistent/.beans"}, "beans data directory was not found"},
		{"relative path", []string{"--beans-bin=beans", "--beans-config=" + cfg, "--beans-data=" + data}, "beans executable was not found"},
		{"no paths at all", nil, "path is not configured"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := exec.Command(binaryPath, append([]string{"--check"}, tc.args...)...).Output()
			if err == nil {
				t.Fatalf("--check succeeded with %v", tc.args)
			}
			var report struct {
				OK     bool   `json:"ok"`
				Error  string `json:"error"`
				Code   string `json:"error_code"`
				Remedy string `json:"remedy"`
			}
			if jsonErr := json.Unmarshal(out, &report); jsonErr != nil {
				t.Fatalf("--check did not produce a JSON report: %v\noutput: %s", jsonErr, out)
			}
			if report.OK {
				t.Error("report claims ok despite a failure")
			}
			if !strings.Contains(report.Error, tc.want) {
				t.Errorf("error = %q, want it to mention %q", report.Error, tc.want)
			}
			if report.Code != "backlog_unavailable" {
				t.Errorf("error code = %q, want backlog_unavailable", report.Code)
			}
			if report.Remedy == "" {
				t.Error("no remedy was offered")
			}
		})
	}
}

func TestCleanShutdown(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())
	cmd := exec.Command(binaryPath, "--beans-bin="+beansPath, "--beans-config="+cfg, "--beans-data="+data)
	var stderr strings.Builder
	cmd.Stderr = &stderr

	client := mcp.NewClient(&mcp.Implementation{Name: "shutdown-test", Version: "1.0.0"}, nil)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	session, err := client.Connect(c, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect failed: %v\n%s", err, stderr.String())
	}
	if _, err := session.ListTools(c, nil); err != nil {
		t.Fatal(err)
	}
	if err := session.Close(); err != nil {
		t.Errorf("close returned an error: %v", err)
	}
	if !strings.Contains(stderr.String(), "stopped cleanly") {
		t.Errorf("the server did not report a clean stop:\n%s", stderr.String())
	}
}

// ---- helpers ----

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func textOf(result *mcp.CallToolResult) string {
	var b strings.Builder
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			b.WriteString(text.Text)
		}
	}
	return b.String()
}

func decodeStructured(t *testing.T, result *mcp.CallToolResult, into any) {
	t.Helper()
	if result.StructuredContent == nil {
		t.Fatalf("no structured content in the result; text was: %s", textOf(result))
	}
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		t.Fatalf("could not decode structured content %s: %v", raw, err)
	}
}

func errorCode(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if result.StructuredContent == nil {
		return ""
	}
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(raw, &doc)
	return doc.Code
}

// TestBacklogIsolationFromUpwardDiscovery proves the adapter answers only from the
// configured backlog. A decoy .beans.yml and .beans directory are planted in the
// process working directory's parent, which is exactly where Beans would find them if
// the adapter let the CLI discover config by searching upward.
func TestBacklogIsolationFromUpwardDiscovery(t *testing.T) {
	cfg, data := backlog(t, sampleTasks())

	decoyRoot := t.TempDir()
	decoyData := filepath.Join(decoyRoot, ".beans")
	if err := os.MkdirAll(decoyData, 0o755); err != nil {
		t.Fatal(err)
	}
	decoyConfig := "beans:\n    path: .beans\n    prefix: other-\n    id_length: 4\n" +
		"    default_status: todo\n    default_type: task\n"
	if err := os.WriteFile(filepath.Join(decoyRoot, ".beans.yml"), []byte(decoyConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	// A task id that exists ONLY in the decoy backlog.
	if err := os.WriteFile(filepath.Join(decoyData, "other-777--secret.md"),
		[]byte(task("other-777", "Private task from an unrelated backlog", "todo", "Should never be readable.")), 0o644); err != nil {
		t.Fatal(err)
	}

	// Start the adapter with its working directory inside the decoy tree, and with
	// BEANS_PATH pointing at the decoy too: both of the ways discovery could go wrong.
	childDir := filepath.Join(decoyRoot, "child")
	if err := os.MkdirAll(childDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binaryPath,
		"--beans-bin="+beansPath, "--beans-config="+cfg, "--beans-data="+data)
	cmd.Dir = childDir
	cmd.Env = append(os.Environ(), "BEANS_PATH="+decoyData)
	var stderr strings.Builder
	cmd.Stderr = &stderr

	client := mcp.NewClient(&mcp.Implementation{Name: "isolation-test", Version: "1.0.0"}, nil)
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	session, err := client.Connect(c, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect failed: %v\n%s", err, stderr.String())
	}
	defer session.Close()

	// The decoy task must be invisible.
	result, err := session.CallTool(c, &mcp.CallToolParams{
		Name: "get_task", Arguments: map[string]any{"id": "other-777"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatalf("a task from the unrelated backlog was returned: %s", textOf(result))
	}
	if code := errorCode(t, result); code != "task_not_found" {
		t.Errorf("error code = %q, want task_not_found", code)
	}

	// And the configured backlog must still answer.
	result, err = session.CallTool(c, &mcp.CallToolParams{
		Name: "get_task", Arguments: map[string]any{"id": "task-101"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("the configured backlog stopped answering: %s", textOf(result))
	}

	// list_tasks must not include the decoy task either.
	result, err = session.CallTool(c, &mcp.CallToolParams{Name: "list_tasks"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(textOf(result), "other-777") {
		t.Errorf("list_tasks leaked the unrelated backlog: %s", textOf(result))
	}
}
