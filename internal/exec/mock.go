package exec

import (
	"context"
	"fmt"
	"strings"
)

// Call records a single command invocation.
type Call struct {
	Dir  string
	Name string
	Args []string
}

func (c Call) String() string {
	if c.Dir != "" {
		return fmt.Sprintf("[%s] %s %s", c.Dir, c.Name, strings.Join(c.Args, " "))
	}
	return fmt.Sprintf("%s %s", c.Name, strings.Join(c.Args, " "))
}

// MockExecutor records calls and returns pre-configured results for testing.
type MockExecutor struct {
	Calls    []Call
	handlers map[string]func(Call) (Result, error)
	// Default result when no handler matches.
	DefaultResult Result
}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		handlers: make(map[string]func(Call) (Result, error)),
	}
}

// OnCommand registers a handler for a command name.
func (m *MockExecutor) OnCommand(name string, handler func(Call) (Result, error)) {
	m.handlers[name] = handler
}

func (m *MockExecutor) Run(ctx context.Context, name string, args ...string) (Result, error) {
	return m.RunDir(ctx, "", name, args...)
}

func (m *MockExecutor) RunDir(ctx context.Context, dir string, name string, args ...string) (Result, error) {
	call := Call{Dir: dir, Name: name, Args: args}
	m.Calls = append(m.Calls, call)

	if handler, ok := m.handlers[name]; ok {
		return handler(call)
	}
	return m.DefaultResult, nil
}
