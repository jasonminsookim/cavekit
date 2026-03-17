package frontier

import (
	"os"
	"path/filepath"
	"testing"
)

const testImplFile = `---
created: "2026-03-17"
---
# Implementation Tracking: Test

| Task | Status | Notes |
|------|--------|-------|
| T-001 | DONE | Completed module init |
| T-002 | DONE | Tmux session management |
| T-003 | IN PROGRESS | Working on capture |
| T-010 | BLOCKED | Waiting on T-009 |
`

func TestTrackStatus(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "impl-test.md"), []byte(testImplFile), 0644)

	statuses, err := TrackStatus(tmp)
	if err != nil {
		t.Fatalf("TrackStatus: %v", err)
	}

	if statuses["T-001"] != TaskDone {
		t.Errorf("T-001 = %v, want DONE", statuses["T-001"])
	}
	if statuses["T-002"] != TaskDone {
		t.Errorf("T-002 = %v, want DONE", statuses["T-002"])
	}
	if statuses["T-003"] != TaskInProgress {
		t.Errorf("T-003 = %v, want IN PROGRESS", statuses["T-003"])
	}
	if statuses["T-010"] != TaskBlocked {
		t.Errorf("T-010 = %v, want BLOCKED", statuses["T-010"])
	}
}

func TestTrackStatus_WordBoundary(t *testing.T) {
	// T-1 should NOT match T-10
	tmp := t.TempDir()
	content := `| Task | Status | Notes |
|------|--------|-------|
| T-10 | DONE | Task ten done |
`
	os.WriteFile(filepath.Join(tmp, "impl-test.md"), []byte(content), 0644)

	statuses, err := TrackStatus(tmp)
	if err != nil {
		t.Fatalf("TrackStatus: %v", err)
	}

	if statuses["T-10"] != TaskDone {
		t.Errorf("T-10 should be DONE, got %v", statuses["T-10"])
	}
	// T-1 should not exist in the map
	if _, exists := statuses["T-1"]; exists {
		t.Error("T-1 should not match T-10")
	}
}

func TestTrackStatus_MultipleDirs(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	os.WriteFile(filepath.Join(dir1, "impl-a.md"), []byte(`| T-001 | DONE | a |
`), 0644)
	os.WriteFile(filepath.Join(dir2, "impl-b.md"), []byte(`| T-002 | DONE | b |
`), 0644)

	statuses, err := TrackStatus(dir1, dir2)
	if err != nil {
		t.Fatalf("TrackStatus: %v", err)
	}

	if statuses["T-001"] != TaskDone {
		t.Errorf("T-001 from dir1 not found")
	}
	if statuses["T-002"] != TaskDone {
		t.Errorf("T-002 from dir2 not found")
	}
}

func TestComputeProgress(t *testing.T) {
	frontier := &Frontier{
		Tasks: []Task{
			{ID: "T-001"},
			{ID: "T-002"},
			{ID: "T-003"},
			{ID: "T-004"},
			{ID: "T-005"},
		},
	}

	statuses := TaskStatusMap{
		"T-001": TaskDone,
		"T-002": TaskDone,
		"T-003": TaskInProgress,
	}

	summary := ComputeProgress(frontier, statuses)
	if summary.Total != 5 {
		t.Errorf("Total = %d, want 5", summary.Total)
	}
	if summary.Done != 2 {
		t.Errorf("Done = %d, want 2", summary.Done)
	}
	if summary.InProgress != 1 {
		t.Errorf("InProgress = %d, want 1", summary.InProgress)
	}
	if summary.Remaining != 2 {
		t.Errorf("Remaining = %d, want 2", summary.Remaining)
	}
}

func TestTaskStatus_String(t *testing.T) {
	if TaskDone.String() != "DONE" {
		t.Errorf("TaskDone.String() = %q", TaskDone.String())
	}
	if TaskPending.String() != "PENDING" {
		t.Errorf("TaskPending.String() = %q", TaskPending.String())
	}
}
