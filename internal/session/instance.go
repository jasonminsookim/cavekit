// Package session defines the instance model that ties together a tmux session,
// git worktree, and frontier into a single manageable unit.
package session

import "time"

// Status represents the current state of an instance.
type Status int

const (
	StatusLoading Status = iota
	StatusRunning
	StatusReady
	StatusPaused
	StatusDone
)

func (s Status) String() string {
	switch s {
	case StatusLoading:
		return "Loading"
	case StatusRunning:
		return "Running"
	case StatusReady:
		return "Ready"
	case StatusPaused:
		return "Paused"
	case StatusDone:
		return "Done"
	default:
		return "Unknown"
	}
}

// Icon returns a display icon for the status.
func (s Status) Icon() string {
	switch s {
	case StatusLoading:
		return "◌"
	case StatusRunning:
		return "⟳"
	case StatusReady:
		return "●"
	case StatusPaused:
		return "⏸"
	case StatusDone:
		return "✓"
	default:
		return "?"
	}
}

// Instance represents one Claude Code agent working on one frontier.
type Instance struct {
	Title        string    `json:"title"`
	FrontierPath string    `json:"frontier_path"`
	WorktreePath string    `json:"worktree_path"`
	TmuxSession  string    `json:"tmux_session"`
	Status       Status    `json:"status"`
	Program      string    `json:"program"`
	CreatedAt    time.Time `json:"created_at"`

	// Progress fields (updated periodically from frontier tracking).
	TasksDone  int    `json:"tasks_done"`
	TasksTotal int    `json:"tasks_total"`
	CurrentTier int   `json:"current_tier"`
	CurrentTask string `json:"current_task"`

	// Diff stats (updated periodically from worktree).
	BranchName  string `json:"branch_name"`
	DiffAdded   int    `json:"diff_added"`
	DiffRemoved int    `json:"diff_removed"`
}

// NewInstance creates a new instance with default values.
func NewInstance(title, frontierPath, program string) *Instance {
	return &Instance{
		Title:        title,
		FrontierPath: frontierPath,
		Status:       StatusLoading,
		Program:      program,
		CreatedAt:    time.Now(),
	}
}

// IsActive returns true if the instance is running or ready (not paused/done).
func (i *Instance) IsActive() bool {
	return i.Status == StatusRunning || i.Status == StatusReady || i.Status == StatusLoading
}

// ProgressString returns a compact progress display.
func (i *Instance) ProgressString() string {
	if i.TasksTotal == 0 {
		return ""
	}
	return i.Status.Icon() + " " + i.Title + " " + formatFraction(i.TasksDone, i.TasksTotal)
}

func formatFraction(a, b int) string {
	return itoa(a) + "/" + itoa(b)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
