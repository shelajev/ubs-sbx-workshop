// Package server exposes the workshop backlog as MCP tools over stdio.
//
// Protocol traffic uses stdout exclusively. Every diagnostic goes to stderr, because
// a stray stdout write corrupts the JSON-RPC stream and is one of the hardest
// failures to diagnose from inside a sandbox.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shelajev/ubs-sbx-workshop/mcp/beans/internal/backlog"
)

// Version is set at build time with -ldflags "-X .../internal/server.Version=...".
var Version = "dev"

// ServerName is the MCP implementation name clients see.
const ServerName = "workshop-beans"

// GetTaskArgs is the get_task input schema.
type GetTaskArgs struct {
	ID string `json:"id" jsonschema:"the workshop task id, for example task-102"`
}

// ListTasksArgs is the complete list_tasks input schema. There is no free-text query
// and no way to express an arbitrary filter: that is the point of the tool.
type ListTasksArgs struct {
	Status string `json:"status,omitempty" jsonschema:"optional status filter: todo, in-progress, draft, completed or scrapped"`
	Tag    string `json:"tag,omitempty" jsonschema:"optional tag filter, lowercase letters digits and dashes"`
	Limit  int    `json:"limit,omitempty" jsonschema:"maximum number of tasks to return, 1 to 50, default 20"`
}

// AddNoteArgs is the presenter-only write tool input.
type AddNoteArgs struct {
	ID   string `json:"id" jsonschema:"the workshop task id to annotate"`
	Note string `json:"note" jsonschema:"the note text to append, at most 2000 characters"`
}

// Options configures the MCP server.
type Options struct {
	Client *backlog.Client
	Logger *slog.Logger
	// ExposeNoteTool registers add_task_note. The client still refuses the call
	// unless the backlog is marked disposable, so registration alone grants nothing.
	ExposeNoteTool bool
}

// New builds the MCP server with the workshop's tool surface.
func New(opts Options) *mcp.Server {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	srv := mcp.NewServer(&mcp.Implementation{
		Name:    ServerName,
		Title:   "Workshop Beans backlog (read-only)",
		Version: Version,
	}, &mcp.ServerOptions{
		Instructions: "Read the assigned workshop task with get_task(id). " +
			"Record the returned version value in your report so a later read can " +
			"detect that the task changed. This backlog is read-only.",
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "get_task",
		Description: "Return one workshop task: title, body, acceptance criteria, status, " +
			"dependencies and a content version (etag) from the host backlog.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args GetTaskArgs) (*mcp.CallToolResult, any, error) {
		task, err := opts.Client.Get(ctx, args.ID)
		if err != nil {
			return toolError(logger, "get_task", err), nil, nil
		}
		logger.Info("get_task served", "id", task.ID, "version", task.Version)
		return textResult(summariseTask(task)), task, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "list_tasks",
		Description: "List workshop tasks with an optional status and tag filter. " +
			"The result is bounded; use get_task for the full body.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args ListTasksArgs) (*mcp.CallToolResult, any, error) {
		result, err := opts.Client.List(ctx, backlog.ListFilter{
			Status: args.Status, Tag: args.Tag, Limit: args.Limit,
		})
		if err != nil {
			return toolError(logger, "list_tasks", err), nil, nil
		}
		logger.Info("list_tasks served", "count", result.Count, "truncated", result.Truncated)
		return textResult(summariseList(result)), result, nil
	})

	if opts.ExposeNoteTool {
		mcp.AddTool(srv, &mcp.Tool{
			Name: "add_task_note",
			Description: "Presenter demonstration only: append a note to a task in a " +
				"disposable backlog. Refused unless the backlog carries the disposable marker.",
			Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(false)},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, args AddNoteArgs) (*mcp.CallToolResult, any, error) {
			task, err := opts.Client.AddNote(ctx, args.ID, args.Note)
			if err != nil {
				return toolError(logger, "add_task_note", err), nil, nil
			}
			logger.Warn("add_task_note applied", "id", task.ID, "version", task.Version)
			return textResult(summariseTask(task)), task, nil
		})
		logger.Warn("add_task_note tool registered",
			"armed", opts.Client.WriteArmed(), "reason", opts.Client.WriteRefusalReason())
	}

	return srv
}

func boolPtr(b bool) *bool { return &b }

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}
}

// toolError converts a typed backlog error into an MCP tool error. The agent sees a
// stable code and a safe message; the host path or CLI detail is logged to stderr only.
func toolError(logger *slog.Logger, tool string, err error) *mcp.CallToolResult {
	var typed *backlog.Error
	code := "backlog_error"
	message := "the backlog could not answer this request"
	if errors.As(err, &typed) {
		code = string(typed.Code)
		message = typed.Message
		logger.Error("tool call failed", "tool", tool, "code", code, "detail", typed.Detail)
	} else {
		logger.Error("tool call failed", "tool", tool, "error", err.Error())
	}
	payload, _ := json.Marshal(map[string]string{"code": code, "message": message})
	return &mcp.CallToolResult{
		IsError:           true,
		Content:           []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("%s: %s", code, message)}},
		StructuredContent: json.RawMessage(payload),
	}
}

func summariseTask(t *backlog.Task) string {
	out := fmt.Sprintf("%s  %s\nstatus: %s   type: %s   priority: %s   version: %s\nbacklog: %s\n",
		t.ID, t.Title, t.Status, orDash(t.Type), orDash(t.Priority), t.Version, t.Backlog)
	if len(t.BlockedBy) > 0 {
		out += fmt.Sprintf("blocked by: %v\n", t.BlockedBy)
	}
	if len(t.Blocking) > 0 {
		out += fmt.Sprintf("blocking: %v\n", t.Blocking)
	}
	out += "\n" + t.Body + "\n"
	return out
}

func summariseList(r *backlog.ListResult) string {
	out := fmt.Sprintf("%d task(s) in the %s backlog", r.Count, r.Backlog)
	if r.Filter.Status != "" || r.Filter.Tag != "" {
		out += fmt.Sprintf(" (status=%s tag=%s)", orDash(r.Filter.Status), orDash(r.Filter.Tag))
	}
	if r.Truncated {
		out += fmt.Sprintf("; truncated at limit %d", r.Limit)
	}
	out += "\n"
	for _, t := range r.Tasks {
		out += fmt.Sprintf("  %-12s %-12s %s\n", t.ID, t.Status, t.Title)
	}
	return out
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
