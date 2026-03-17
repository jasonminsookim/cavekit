package session

import (
	"path/filepath"

	"github.com/julb/blueprint-monitor/internal/frontier"
)

// UpdateProgress refreshes the progress fields on an instance
// by reading the frontier file and impl tracking files.
func UpdateProgress(inst *Instance) error {
	if inst.FrontierPath == "" {
		return nil
	}

	// Parse the frontier file
	f, err := frontier.Parse(inst.FrontierPath)
	if err != nil {
		return err
	}

	// Build impl directories to scan
	var implDirs []string
	if inst.WorktreePath != "" {
		implDirs = append(implDirs, filepath.Join(inst.WorktreePath, "context", "impl"))
	}

	// Track task statuses
	statuses, err := frontier.TrackStatus(implDirs...)
	if err != nil {
		return err
	}

	// Compute summary
	summary := frontier.ComputeProgress(f, statuses)
	inst.TasksDone = summary.Done
	inst.TasksTotal = summary.Total

	// Find current tier and task
	for _, task := range f.Tasks {
		status, exists := statuses[task.ID]
		if !exists || status == frontier.TaskPending {
			inst.CurrentTier = task.Tier
			inst.CurrentTask = task.ID
			break
		}
	}

	return nil
}
