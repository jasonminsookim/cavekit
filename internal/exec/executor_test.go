package exec

import (
	"context"
	"testing"
)

func TestRealExecutor_Echo(t *testing.T) {
	e := NewRealExecutor()
	res, err := e.Run(context.Background(), "echo", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stdout != "hello\n" {
		t.Errorf("got stdout=%q, want %q", res.Stdout, "hello\n")
	}
	if res.ExitCode != 0 {
		t.Errorf("got exit code %d, want 0", res.ExitCode)
	}
}

func TestRealExecutor_NonZeroExit(t *testing.T) {
	e := NewRealExecutor()
	res, err := e.Run(context.Background(), "sh", "-c", "exit 42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 42 {
		t.Errorf("got exit code %d, want 42", res.ExitCode)
	}
}

func TestRealExecutor_RunDir(t *testing.T) {
	e := NewRealExecutor()
	res, err := e.RunDir(context.Background(), "/tmp", "pwd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stdout == "" {
		t.Error("expected non-empty stdout from pwd")
	}
}

func TestMockExecutor_RecordsCalls(t *testing.T) {
	m := NewMockExecutor()
	m.DefaultResult = Result{Stdout: "ok\n"}

	res, err := m.Run(context.Background(), "tmux", "has-session", "-t", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stdout != "ok\n" {
		t.Errorf("got %q, want %q", res.Stdout, "ok\n")
	}
	if len(m.Calls) != 1 {
		t.Fatalf("got %d calls, want 1", len(m.Calls))
	}
	if m.Calls[0].Name != "tmux" {
		t.Errorf("got name=%q, want %q", m.Calls[0].Name, "tmux")
	}
}

func TestMockExecutor_Handler(t *testing.T) {
	m := NewMockExecutor()
	m.OnCommand("git", func(c Call) (Result, error) {
		return Result{Stdout: "main\n"}, nil
	})

	res, _ := m.Run(context.Background(), "git", "branch", "--show-current")
	if res.Stdout != "main\n" {
		t.Errorf("got %q, want %q", res.Stdout, "main\n")
	}
}
