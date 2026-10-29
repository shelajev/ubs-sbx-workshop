// Command beans-mcp exposes the workshop's Beans backlog to sandboxed coding agents
// as a read-only MCP stdio server.
//
// It is registered on the HOST with:
//
//	sbx mcp add workshop-beans --command /path/to/beans-mcp \
//	  --args --beans-bin=/path/to/beans \
//	  --args --beans-config=/path/to/control/beans/.beans.yml \
//	  --args --beans-data=/path/to/control/beans/.beans \
//	  --dir /path/to/control
//
// Every path is explicit: the server never searches upward for a .beans.yml, so it
// cannot answer from the host user's unrelated backlog.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shelajev/ubs-sbx-workshop/mcp/beans/internal/backlog"
	"github.com/shelajev/ubs-sbx-workshop/mcp/beans/internal/server"
)

var (
	version = "dev"
	commit  = "unknown"
	buildAt = "unknown"
)

type options struct {
	beansBin    string
	beansConfig string
	beansData   string
	backlogName string
	timeout     time.Duration
	logLevel    string
	enableNotes bool
	showVersion bool
	check       bool
}

func main() {
	opts := parseFlags()

	// Diagnostics always go to stderr. stdout carries JSON-RPC only.
	logger := newLogger(opts.logLevel)
	server.Version = version

	if opts.showVersion {
		printVersion()
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := openBacklog(ctx, opts, logger)
	if err != nil {
		if opts.check {
			emitCheck(false, opts, nil, err)
			os.Exit(1)
		}
		logger.Error("startup failed", "error", err.Error())
		fmt.Fprintf(os.Stderr, "\nbeans-mcp could not start: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run `beans-mcp --check` with the same flags to see a machine-readable diagnosis.\n")
		os.Exit(1)
	}

	if opts.check {
		emitCheck(true, opts, client, nil)
		return
	}

	logger.Info("workshop-beans starting",
		"version", version, "commit", commit, "beans", client.BeansVersion(),
		"backlog", client.Name(), "note_writes", client.WriteArmed(), "pid", os.Getpid())

	mcpServer := server.New(server.Options{
		Client:         client,
		Logger:         logger,
		ExposeNoteTool: opts.enableNotes,
	})

	if err := mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("server stopped with an error", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("workshop-beans stopped cleanly")
}

func parseFlags() options {
	var opts options
	fs := flag.NewFlagSet("beans-mcp", flag.ExitOnError)
	fs.StringVar(&opts.beansBin, "beans-bin", envOr("WORKSHOP_BEANS_BIN", ""),
		"absolute path to the pinned beans executable")
	fs.StringVar(&opts.beansConfig, "beans-config", envOr("WORKSHOP_BEANS_CONFIG", ""),
		"absolute path to the workshop .beans.yml")
	fs.StringVar(&opts.beansData, "beans-data", envOr("WORKSHOP_BEANS_DATA", ""),
		"absolute path to the workshop .beans data directory")
	fs.StringVar(&opts.backlogName, "backlog-name", envOr("WORKSHOP_BACKLOG_NAME", "workshop"),
		"label reported to clients instead of the host path")
	fs.DurationVar(&opts.timeout, "timeout", 10*time.Second, "per-call timeout for the beans CLI")
	fs.StringVar(&opts.logLevel, "log-level", envOr("WORKSHOP_BEANS_LOG", "info"),
		"stderr log level: debug, info, warn or error")
	fs.BoolVar(&opts.enableNotes, "enable-presenter-note-tool", false,
		"register the presenter-only add_task_note tool (still refused unless the backlog is marked disposable)")
	fs.BoolVar(&opts.showVersion, "version", false, "print version metadata and exit")
	fs.BoolVar(&opts.check, "check", false, "validate configuration, print a JSON report and exit")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "beans-mcp %s — read-only MCP stdio server over a Beans backlog\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: beans-mcp --beans-bin PATH --beans-config PATH --beans-data PATH\n\n")
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[1:])

	// Absolute paths are required, but expanding a relative one here gives a clearer
	// error than letting the CLI resolve it against an unknown working directory.
	opts.beansBin = absOrRaw(opts.beansBin)
	opts.beansConfig = absOrRaw(opts.beansConfig)
	opts.beansData = absOrRaw(opts.beansData)
	return opts
}

func absOrRaw(p string) string {
	if p == "" || filepath.IsAbs(p) {
		return p
	}
	if abs, err := filepath.Abs(p); err == nil {
		return abs
	}
	return p
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
}

func openBacklog(ctx context.Context, opts options, logger *slog.Logger) (*backlog.Client, error) {
	return backlog.Open(ctx, backlog.Config{
		BeansBin:        opts.beansBin,
		ConfigFile:      opts.beansConfig,
		DataDir:         opts.beansData,
		Name:            opts.backlogName,
		Timeout:         opts.timeout,
		AllowNoteWrites: opts.enableNotes,
		Logger:          logger,
	})
}

func printVersion() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(map[string]string{
		"name":     server.ServerName,
		"version":  version,
		"commit":   commit,
		"built_at": buildAt,
		"go":       runtime.Version(),
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"mcp_sdk":  "github.com/modelcontextprotocol/go-sdk v1.8.0",
		"tool_surface": "get_task, list_tasks" +
			" (add_task_note only with --enable-presenter-note-tool and a disposable backlog)",
	})
}

// checkReport is the machine-readable --check output. It reports paths the operator
// supplied, never secrets.
type checkReport struct {
	OK          bool     `json:"ok"`
	Version     string   `json:"version"`
	BeansBin    string   `json:"beans_bin"`
	BeansConfig string   `json:"beans_config"`
	BeansData   string   `json:"beans_data"`
	BacklogName string   `json:"backlog_name"`
	BeansVer    string   `json:"beans_version,omitempty"`
	TaskCount   int      `json:"task_count,omitempty"`
	SampleIDs   []string `json:"sample_task_ids,omitempty"`
	NoteWrites  string   `json:"note_writes"`
	ErrorCode   string   `json:"error_code,omitempty"`
	Error       string   `json:"error,omitempty"`
	Remedy      string   `json:"remedy,omitempty"`
}

func emitCheck(ok bool, opts options, client *backlog.Client, err error) {
	report := checkReport{
		OK: ok, Version: version,
		BeansBin: opts.beansBin, BeansConfig: opts.beansConfig,
		BeansData: opts.beansData, BacklogName: opts.backlogName,
		NoteWrites: "disabled",
	}
	if client != nil {
		report.BeansVer = client.BeansVersion()
		if client.WriteArmed() {
			report.NoteWrites = "enabled (disposable backlog)"
		} else if opts.enableNotes {
			report.NoteWrites = "refused: " + client.WriteRefusalReason()
		}
		ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
		defer cancel()
		if list, listErr := client.List(ctx, backlog.ListFilter{Limit: backlog.MaxListLimit}); listErr == nil {
			report.TaskCount = list.Count
			for _, t := range list.Tasks {
				if len(report.SampleIDs) == 5 {
					break
				}
				report.SampleIDs = append(report.SampleIDs, t.ID)
			}
		} else {
			report.OK = false
			report.Error = listErr.Error()
			report.Remedy = "Check that the data directory contains bean files and that `beans list --json` works with the same --config and --beans-path."
		}
	}
	if err != nil {
		var typed *backlog.Error
		if errors.As(err, &typed) {
			report.ErrorCode = string(typed.Code)
			report.Error = typed.Message
			report.Remedy = remedyFor(typed.Code)
		} else {
			report.Error = err.Error()
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(report)
}

func remedyFor(code backlog.Code) string {
	switch code {
	case backlog.CodeBacklogUnavail:
		return "Verify --beans-bin, --beans-config and --beans-data are absolute paths that exist, and that the beans binary is executable. Run scripts/backlog-init.sh to create the workshop backlog."
	case backlog.CodeBacklogTimeout:
		return "The beans CLI did not answer in time. Try running it directly with the same --config and --beans-path."
	default:
		return ""
	}
}
