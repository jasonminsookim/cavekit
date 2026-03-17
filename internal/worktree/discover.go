package worktree

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredWorktree represents a found Blueprint worktree.
type DiscoveredWorktree struct {
	Path         string // Full worktree path
	Branch       string // Branch name (blueprint/xxx)
	FrontierName string // Derived frontier name
	HasRalphLoop bool   // .claude/ralph-loop.local.md exists
}

// DiscoverAll finds all existing Blueprint worktrees for the given project.
// Scans {project_root}/../{project_name}-blueprint-* directories.
func DiscoverAll(projectRoot string) ([]DiscoveredWorktree, error) {
	projectName := filepath.Base(projectRoot)
	parentDir := filepath.Dir(projectRoot)
	pattern := projectName + "-blueprint-"

	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, err
	}

	var results []DiscoveredWorktree
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, pattern) {
			continue
		}

		frontierName := strings.TrimPrefix(name, pattern)
		wtPath := filepath.Join(parentDir, name)

		// Check for active Ralph Loop
		ralphLoopPath := filepath.Join(wtPath, ".claude", "ralph-loop.local.md")
		hasRalphLoop := fileExists(ralphLoopPath)

		results = append(results, DiscoveredWorktree{
			Path:         wtPath,
			Branch:       BranchName(frontierName),
			FrontierName: frontierName,
			HasRalphLoop: hasRalphLoop,
		})
	}

	return results, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
