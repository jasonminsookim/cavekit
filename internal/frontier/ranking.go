package frontier

import (
	"fmt"
	"sort"
)

// RankedCandidate represents a frontier with its computed score.
type RankedCandidate struct {
	File     FrontierFile
	Score    int
	Status   FrontierStatus
	Selected bool
}

// RankAndSelect ranks frontier candidates and selects the best one.
// Returns error if filter is set and matches zero frontiers.
//
// Scoring: active worktree with Ralph Loop = 3, worktree exists or incomplete tasks = 2, base = 1.
// Ties break alphabetically (first in order wins, using > not >= for determinism).
func RankAndSelect(candidates []FrontierFile, statuses TaskStatusMap, filter string, worktreeChecker func(name string) (bool, bool)) ([]RankedCandidate, error) {
	// Apply filter if set
	if filter != "" {
		var filtered []FrontierFile
		for _, c := range candidates {
			if c.Name == filter {
				filtered = append(filtered, c)
			}
		}
		if len(filtered) == 0 {
			names := make([]string, len(candidates))
			for i, c := range candidates {
				names[i] = c.Name
			}
			return nil, fmt.Errorf("filter %q matches no frontiers. Available: %v", filter, names)
		}
		candidates = filtered
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	// Sort alphabetically first for deterministic tie-breaking
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})

	// Score each candidate
	ranked := make([]RankedCandidate, len(candidates))
	for i, c := range candidates {
		wtExists, hasLoop := false, false
		if worktreeChecker != nil {
			wtExists, hasLoop = worktreeChecker(c.Name)
		}

		score := 1 // base score
		if hasLoop {
			score = 3 // active Ralph Loop
		} else if wtExists {
			score = 2 // worktree exists
		} else {
			// Check if has incomplete tasks
			// (any candidate that made it here likely has tasks)
			score = 2
		}

		ranked[i] = RankedCandidate{
			File:   c,
			Score:  score,
			Status: classifyFromScore(score),
		}
	}

	// Select highest score (first alphabetically wins ties due to pre-sort + > comparison)
	bestIdx := 0
	for i := 1; i < len(ranked); i++ {
		if ranked[i].Score > ranked[bestIdx].Score {
			bestIdx = i
		}
	}
	ranked[bestIdx].Selected = true

	return ranked, nil
}

func classifyFromScore(score int) FrontierStatus {
	switch score {
	case 3:
		return FrontierInProgress
	case 2:
		return FrontierAvailable
	default:
		return FrontierAvailable
	}
}

// FormatCandidates returns a display-friendly string of candidates.
// The selected one is marked with →.
func FormatCandidates(ranked []RankedCandidate) string {
	var result string
	for _, r := range ranked {
		marker := "  "
		if r.Selected {
			marker = "→ "
		}
		result += fmt.Sprintf("%s%s (score: %d, status: %s)\n", marker, r.File.Name, r.Score, r.Status)
	}
	return result
}
