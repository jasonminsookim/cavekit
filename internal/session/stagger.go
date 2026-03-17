package session

import (
	"context"
	"time"
)

// StaggeredLauncher launches multiple instances with a delay between them.
type StaggeredLauncher struct {
	mgr         *Manager
	delay       time.Duration
	projectRoot string
}

// NewStaggeredLauncher creates a staggered launcher.
func NewStaggeredLauncher(mgr *Manager, projectRoot string, delay time.Duration) *StaggeredLauncher {
	if delay == 0 {
		delay = 5 * time.Second
	}
	return &StaggeredLauncher{
		mgr:         mgr,
		delay:       delay,
		projectRoot: projectRoot,
	}
}

// LaunchAll starts multiple instances with a delay between each.
// The first starts immediately; subsequent wait for the delay.
// Runs in background — returns immediately.
func (l *StaggeredLauncher) LaunchAll(ctx context.Context, instances []*Instance, frontierNames []string) {
	go func() {
		for i, inst := range instances {
			if i > 0 {
				select {
				case <-time.After(l.delay):
				case <-ctx.Done():
					return
				}
			}
			l.mgr.Start(ctx, inst, l.projectRoot, frontierNames[i], 3*time.Second)
		}
	}()
}
