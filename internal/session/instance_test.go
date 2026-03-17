package session

import "testing"

func TestStatus_String(t *testing.T) {
	tests := []struct {
		s    Status
		want string
	}{
		{StatusLoading, "Loading"},
		{StatusRunning, "Running"},
		{StatusReady, "Ready"},
		{StatusPaused, "Paused"},
		{StatusDone, "Done"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("Status(%d).String() = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestStatus_Icon(t *testing.T) {
	if got := StatusRunning.Icon(); got != "⟳" {
		t.Errorf("StatusRunning.Icon() = %q, want %q", got, "⟳")
	}
	if got := StatusDone.Icon(); got != "✓" {
		t.Errorf("StatusDone.Icon() = %q, want %q", got, "✓")
	}
}

func TestNewInstance(t *testing.T) {
	inst := NewInstance("auth", "/path/to/frontier.md", "claude")
	if inst.Title != "auth" {
		t.Errorf("Title = %q, want %q", inst.Title, "auth")
	}
	if inst.Status != StatusLoading {
		t.Errorf("Status = %v, want %v", inst.Status, StatusLoading)
	}
	if inst.Program != "claude" {
		t.Errorf("Program = %q, want %q", inst.Program, "claude")
	}
	if inst.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestInstance_IsActive(t *testing.T) {
	inst := NewInstance("test", "", "claude")

	inst.Status = StatusRunning
	if !inst.IsActive() {
		t.Error("Running should be active")
	}

	inst.Status = StatusPaused
	if inst.IsActive() {
		t.Error("Paused should not be active")
	}

	inst.Status = StatusDone
	if inst.IsActive() {
		t.Error("Done should not be active")
	}
}

func TestInstance_ProgressString(t *testing.T) {
	inst := NewInstance("auth", "", "claude")
	inst.Status = StatusRunning
	inst.TasksDone = 3
	inst.TasksTotal = 12

	got := inst.ProgressString()
	want := "⟳ auth 3/12"
	if got != want {
		t.Errorf("ProgressString() = %q, want %q", got, want)
	}
}

func TestInstance_ProgressString_Empty(t *testing.T) {
	inst := NewInstance("auth", "", "claude")
	if got := inst.ProgressString(); got != "" {
		t.Errorf("ProgressString() = %q, want empty", got)
	}
}
