package frontier

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClassifyFrontier_Done(t *testing.T) {
	f := &Frontier{
		Tasks: []Task{{ID: "T-001"}, {ID: "T-002"}},
	}
	statuses := TaskStatusMap{
		"T-001": TaskDone,
		"T-002": TaskDone,
	}

	status := ClassifyFrontier(f, statuses, "")
	if status != FrontierDone {
		t.Errorf("should be Done, got %v", status)
	}
}

func TestClassifyFrontier_InProgress(t *testing.T) {
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, ".claude"), 0755)
	os.WriteFile(filepath.Join(tmp, ".claude", "ralph-loop.local.md"), []byte("active"), 0644)

	f := &Frontier{
		Tasks: []Task{{ID: "T-001"}, {ID: "T-002"}},
	}
	statuses := TaskStatusMap{
		"T-001": TaskDone,
	}

	status := ClassifyFrontier(f, statuses, tmp)
	if status != FrontierInProgress {
		t.Errorf("should be InProgress, got %v", status)
	}
}

func TestClassifyFrontier_Available(t *testing.T) {
	f := &Frontier{
		Tasks: []Task{{ID: "T-001"}, {ID: "T-002"}},
	}
	statuses := TaskStatusMap{
		"T-001": TaskDone,
	}

	status := ClassifyFrontier(f, statuses, "")
	if status != FrontierAvailable {
		t.Errorf("should be Available, got %v", status)
	}
}

func TestFrontierStatus_String(t *testing.T) {
	if FrontierDone.String() != "done" {
		t.Errorf("FrontierDone.String() = %q", FrontierDone.String())
	}
	if FrontierInProgress.String() != "in-progress" {
		t.Errorf("FrontierInProgress.String() = %q", FrontierInProgress.String())
	}
	if FrontierAvailable.String() != "available" {
		t.Errorf("FrontierAvailable.String() = %q", FrontierAvailable.String())
	}
}

func TestFrontierStatus_Icon(t *testing.T) {
	if FrontierDone.Icon() != "✓" {
		t.Errorf("FrontierDone.Icon() = %q", FrontierDone.Icon())
	}
	if FrontierInProgress.Icon() != "⟳" {
		t.Errorf("FrontierInProgress.Icon() = %q", FrontierInProgress.Icon())
	}
	if FrontierAvailable.Icon() != "·" {
		t.Errorf("FrontierAvailable.Icon() = %q", FrontierAvailable.Icon())
	}
}
