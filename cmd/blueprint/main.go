package main

import (
	"context"
	"fmt"
	"os"
	osexec "os/exec"

	"github.com/julb/blueprint-monitor/internal/exec"
	"github.com/julb/blueprint-monitor/internal/session"
	"github.com/julb/blueprint-monitor/internal/tui"
	"github.com/julb/blueprint-monitor/internal/tmux"
	"github.com/julb/blueprint-monitor/internal/worktree"
)

const version = "v0.1.0"

func main() {
	cmd := "monitor"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "monitor", "":
		runMonitor()
	case "status":
		runStatus()
	case "kill":
		runKill()
	case "version":
		fmt.Println("blueprint-monitor", version)
	case "debug":
		runDebug()
	case "reset":
		runReset()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		fmt.Fprintln(os.Stderr, "usage: blueprint [monitor|status|kill|version|debug|reset]")
		os.Exit(1)
	}
}

func runMonitor() {
	// Preflight checks
	if err := preflight(); err != nil {
		fmt.Fprintf(os.Stderr, "preflight failed: %s\n", err)
		os.Exit(1)
	}

	// Determine project root
	cwd, _ := os.Getwd()
	executor := exec.NewRealExecutor()
	wtMgr := worktree.NewManager(executor)
	ctx := context.Background()
	root, err := wtMgr.ProjectRoot(ctx, cwd)
	if err != nil {
		root = cwd
	}

	// Parse program flag
	program := "claude"
	for i, arg := range os.Args {
		if (arg == "--program" || arg == "-p") && i+1 < len(os.Args) {
			program = os.Args[i+1]
		}
	}

	// Launch TUI
	if err := tui.Run(root, program); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %s\n", err)
		os.Exit(1)
	}
}

func runStatus() {
	executor := exec.NewRealExecutor()
	wtMgr := worktree.NewManager(executor)

	cwd, _ := os.Getwd()
	root, err := wtMgr.ProjectRoot(nil, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "not in a git repo: %s\n", err)
		os.Exit(1)
	}

	worktrees, err := worktree.DiscoverAll(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "discover worktrees: %s\n", err)
		os.Exit(1)
	}

	if len(worktrees) == 0 {
		fmt.Println("No Blueprint worktrees found.")
		return
	}

	for _, wt := range worktrees {
		icon := "·"
		if wt.HasRalphLoop {
			icon = "⟳"
		}
		fmt.Printf("%s %s: %s\n", icon, wt.FrontierName, wt.Path)
	}
}

func runKill() {
	executor := exec.NewRealExecutor()
	tmuxMgr := tmux.NewManager(executor)
	wtMgr := worktree.NewManager(executor)

	cwd, _ := os.Getwd()
	root, _ := wtMgr.ProjectRoot(nil, cwd)

	// Kill tmux sessions
	sessions, _ := tmuxMgr.ListSessions(nil)
	for _, s := range sessions {
		tmuxMgr.Kill(nil, s)
	}

	// Remove worktrees
	worktrees, _ := worktree.DiscoverAll(root)
	for _, wt := range worktrees {
		wtMgr.Remove(nil, root, wt.FrontierName)
	}

	fmt.Printf("Killed %d sessions, cleaned %d worktrees.\n", len(sessions), len(worktrees))
}

func runDebug() {
	store := session.NewStore("")
	fmt.Println("State file:", store.Path())
	fmt.Println("Version:", version)
}

func runReset() {
	store := session.NewStore("")
	os.Remove(store.Path())
	fmt.Println("State cleared.")
}

func preflight() error {
	if _, err := osexec.LookPath("tmux"); err != nil {
		return fmt.Errorf("tmux not installed")
	}
	if _, err := osexec.LookPath("git"); err != nil {
		return fmt.Errorf("git not installed")
	}
	return nil
}
