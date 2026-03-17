package frontier

import (
	"strings"
	"testing"
)

func TestRankAndSelect_SingleCandidate(t *testing.T) {
	candidates := []FrontierFile{
		{Path: "/path/build-site.md", Name: "build"},
	}

	ranked, err := RankAndSelect(candidates, nil, "", nil)
	if err != nil {
		t.Fatalf("RankAndSelect: %v", err)
	}
	if len(ranked) != 1 {
		t.Fatalf("got %d, want 1", len(ranked))
	}
	if !ranked[0].Selected {
		t.Error("single candidate should be selected")
	}
}

func TestRankAndSelect_ActiveLoopWins(t *testing.T) {
	candidates := []FrontierFile{
		{Path: "/a", Name: "alpha"},
		{Path: "/b", Name: "beta"},
	}

	checker := func(name string) (bool, bool) {
		if name == "beta" {
			return true, true // worktree exists, has loop
		}
		return true, false // worktree exists, no loop
	}

	ranked, err := RankAndSelect(candidates, nil, "", checker)
	if err != nil {
		t.Fatalf("RankAndSelect: %v", err)
	}

	var selected string
	for _, r := range ranked {
		if r.Selected {
			selected = r.File.Name
		}
	}
	if selected != "beta" {
		t.Errorf("beta should win (score 3 vs 2), got %q", selected)
	}
}

func TestRankAndSelect_AlphabeticalTieBreak(t *testing.T) {
	candidates := []FrontierFile{
		{Path: "/c", Name: "charlie"},
		{Path: "/a", Name: "alpha"},
		{Path: "/b", Name: "beta"},
	}

	// All same score (no worktree checker)
	ranked, err := RankAndSelect(candidates, nil, "", nil)
	if err != nil {
		t.Fatalf("RankAndSelect: %v", err)
	}

	var selected string
	for _, r := range ranked {
		if r.Selected {
			selected = r.File.Name
		}
	}
	if selected != "alpha" {
		t.Errorf("alpha should win alphabetical tie, got %q", selected)
	}
}

func TestRankAndSelect_FilterMatch(t *testing.T) {
	candidates := []FrontierFile{
		{Path: "/a", Name: "auth"},
		{Path: "/b", Name: "payments"},
	}

	ranked, err := RankAndSelect(candidates, nil, "payments", nil)
	if err != nil {
		t.Fatalf("RankAndSelect: %v", err)
	}
	if len(ranked) != 1 {
		t.Fatalf("filter should return 1, got %d", len(ranked))
	}
	if ranked[0].File.Name != "payments" {
		t.Errorf("got %q, want payments", ranked[0].File.Name)
	}
}

func TestRankAndSelect_FilterNoMatch(t *testing.T) {
	candidates := []FrontierFile{
		{Path: "/a", Name: "auth"},
	}

	_, err := RankAndSelect(candidates, nil, "missing", nil)
	if err == nil {
		t.Error("should error on filter matching zero")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Errorf("error should mention filter, got: %v", err)
	}
}

func TestFormatCandidates(t *testing.T) {
	ranked := []RankedCandidate{
		{File: FrontierFile{Name: "auth"}, Score: 3, Status: FrontierInProgress, Selected: true},
		{File: FrontierFile{Name: "payments"}, Score: 2, Status: FrontierAvailable},
	}

	output := FormatCandidates(ranked)
	if !strings.Contains(output, "→ auth") {
		t.Errorf("selected should be marked with →: %s", output)
	}
	if strings.Contains(output, "→ payments") {
		t.Error("non-selected should not have →")
	}
}
