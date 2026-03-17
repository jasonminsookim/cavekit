package worktree

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverAll(t *testing.T) {
	// Create temp structure simulating sibling worktrees
	parent := t.TempDir()
	projectRoot := filepath.Join(parent, "myproject")
	os.MkdirAll(projectRoot, 0755)

	// Create blueprint worktrees
	wt1 := filepath.Join(parent, "myproject-blueprint-auth")
	os.MkdirAll(wt1, 0755)

	wt2 := filepath.Join(parent, "myproject-blueprint-payments")
	os.MkdirAll(wt2, 0755)
	// Add ralph loop marker to payments
	os.MkdirAll(filepath.Join(wt2, ".claude"), 0755)
	os.WriteFile(filepath.Join(wt2, ".claude", "ralph-loop.local.md"), []byte("active"), 0644)

	// Non-blueprint dir (should be excluded)
	os.MkdirAll(filepath.Join(parent, "myproject-other"), 0755)

	results, err := DiscoverAll(projectRoot)
	if err != nil {
		t.Fatalf("DiscoverAll: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	// Find auth and payments
	byName := map[string]DiscoveredWorktree{}
	for _, r := range results {
		byName[r.FrontierName] = r
	}

	auth, ok := byName["auth"]
	if !ok {
		t.Fatal("auth worktree not found")
	}
	if auth.Branch != "blueprint/auth" {
		t.Errorf("auth branch = %q", auth.Branch)
	}
	if auth.HasRalphLoop {
		t.Error("auth should not have ralph loop")
	}

	payments, ok := byName["payments"]
	if !ok {
		t.Fatal("payments worktree not found")
	}
	if !payments.HasRalphLoop {
		t.Error("payments should have ralph loop")
	}
}

func TestDiscoverAll_NoWorktrees(t *testing.T) {
	parent := t.TempDir()
	projectRoot := filepath.Join(parent, "myproject")
	os.MkdirAll(projectRoot, 0755)

	results, err := DiscoverAll(projectRoot)
	if err != nil {
		t.Fatalf("DiscoverAll: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0, got %d", len(results))
	}
}
