package worktree

import (
	"context"
	"fmt"
)

// Push commits all changes and pushes the branch to remote.
func (m *Manager) Push(ctx context.Context, wtPath, message string) error {
	// Stage all changes
	res, err := m.exec.RunDir(ctx, wtPath, "git", "add", "-A")
	if err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("git add failed: %s", res.Stderr)
	}

	// Commit
	res, err = m.exec.RunDir(ctx, wtPath, "git", "commit", "-m", message)
	if err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	// Exit code 1 = nothing to commit, not an error
	if res.ExitCode != 0 && res.ExitCode != 1 {
		return fmt.Errorf("git commit failed: %s", res.Stderr)
	}

	// Push with --set-upstream
	res, err = m.exec.RunDir(ctx, wtPath, "git", "push", "--set-upstream", "origin", "HEAD")
	if err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("git push failed (exit %d): %s", res.ExitCode, res.Stderr)
	}

	return nil
}
