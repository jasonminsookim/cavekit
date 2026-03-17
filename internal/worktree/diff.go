package worktree

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// DiffStats holds the summary of changes between a worktree branch and main.
type DiffStats struct {
	FilesChanged int
	Insertions   int
	Deletions    int
}

func (d DiffStats) String() string {
	return fmt.Sprintf("%d files, +%d/-%d", d.FilesChanged, d.Insertions, d.Deletions)
}

// DiffStat returns change statistics between main and the current HEAD.
func (m *Manager) DiffStat(ctx context.Context, wtPath string) (DiffStats, error) {
	res, err := m.exec.RunDir(ctx, wtPath, "git", "diff", "--stat", "main...HEAD")
	if err != nil {
		return DiffStats{}, fmt.Errorf("git diff --stat: %w", err)
	}
	if res.ExitCode != 0 {
		// No main branch or no commits — return empty stats
		return DiffStats{}, nil
	}
	return parseDiffStat(res.Stdout), nil
}

// Diff returns the raw diff output between main and HEAD.
func (m *Manager) Diff(ctx context.Context, wtPath string) (string, error) {
	res, err := m.exec.RunDir(ctx, wtPath, "git", "diff", "main...HEAD")
	if err != nil {
		return "", fmt.Errorf("git diff: %w", err)
	}
	if res.ExitCode != 0 {
		return "", nil
	}
	return res.Stdout, nil
}

// parseDiffStat parses the summary line from git diff --stat output.
// Example: " 5 files changed, 120 insertions(+), 30 deletions(-)"
func parseDiffStat(output string) DiffStats {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return DiffStats{}
	}

	// The summary is always the last line
	summary := lines[len(lines)-1]

	var stats DiffStats
	parts := strings.Split(summary, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		fields := strings.Fields(part)
		if len(fields) < 2 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		switch {
		case strings.Contains(part, "file"):
			stats.FilesChanged = n
		case strings.Contains(part, "insertion"):
			stats.Insertions = n
		case strings.Contains(part, "deletion"):
			stats.Deletions = n
		}
	}
	return stats
}
