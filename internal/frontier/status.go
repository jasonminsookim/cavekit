package frontier

import (
	"os"
	"path/filepath"
)

// FrontierStatus classifies the overall status of a frontier.
type FrontierStatus int

const (
	FrontierAvailable  FrontierStatus = iota // Has incomplete tasks, no active worktree
	FrontierInProgress                       // Has an active worktree with Ralph Loop running
	FrontierDone                             // All tasks complete
)

func (s FrontierStatus) String() string {
	switch s {
	case FrontierDone:
		return "done"
	case FrontierInProgress:
		return "in-progress"
	default:
		return "available"
	}
}

// Icon returns a display icon for the frontier status.
func (s FrontierStatus) Icon() string {
	switch s {
	case FrontierDone:
		return "✓"
	case FrontierInProgress:
		return "⟳"
	default:
		return "·"
	}
}

// ClassifyFrontier determines the overall status of a frontier.
func ClassifyFrontier(f *Frontier, statuses TaskStatusMap, worktreePath string) FrontierStatus {
	summary := ComputeProgress(f, statuses)

	// All done?
	if summary.Done == summary.Total && summary.Total > 0 {
		return FrontierDone
	}

	// Check for active Ralph Loop in worktree
	if worktreePath != "" {
		ralphLoopPath := filepath.Join(worktreePath, ".claude", "ralph-loop.local.md")
		if _, err := os.Stat(ralphLoopPath); err == nil {
			return FrontierInProgress
		}
	}

	return FrontierAvailable
}
