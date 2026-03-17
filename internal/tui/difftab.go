package tui

import (
	"context"
	"strings"

	"github.com/julb/blueprint-monitor/internal/worktree"
)

// DiffTab renders git diff output for the selected instance.
type DiffTab struct {
	wtMgr      *worktree.Manager
	rawDiff    string
	stats      worktree.DiffStats
	scrollPos  int
}

// NewDiffTab creates a diff tab.
func NewDiffTab(wtMgr *worktree.Manager) *DiffTab {
	return &DiffTab{wtMgr: wtMgr}
}

// Refresh updates the diff content for the given worktree.
func (d *DiffTab) Refresh(ctx context.Context, wtPath string) {
	if wtPath == "" {
		d.rawDiff = ""
		d.stats = worktree.DiffStats{}
		return
	}

	stats, err := d.wtMgr.DiffStat(ctx, wtPath)
	if err == nil {
		d.stats = stats
	}

	diff, err := d.wtMgr.Diff(ctx, wtPath)
	if err == nil {
		d.rawDiff = diff
	}
}

// Content returns the styled diff output.
func (d *DiffTab) Content() string {
	if d.rawDiff == "" {
		return "No diff available."
	}

	header := DiffHeaderStyle.Render(d.stats.String()) + "\n\n"

	// Apply basic diff coloring
	lines := strings.Split(d.rawDiff, "\n")
	var styled []string
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			styled = append(styled, DiffAddStyle.Render(line))
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			styled = append(styled, DiffRemoveStyle.Render(line))
		case strings.HasPrefix(line, "@@"):
			styled = append(styled, DiffHeaderStyle.Render(line))
		default:
			styled = append(styled, line)
		}
	}

	// Apply scroll position
	if d.scrollPos > 0 && d.scrollPos < len(styled) {
		styled = styled[d.scrollPos:]
	} else if d.scrollPos >= len(styled) {
		d.scrollPos = max(0, len(styled)-1)
		if len(styled) > 0 {
			styled = styled[len(styled)-1:]
		}
	}

	return header + strings.Join(styled, "\n")
}

// ScrollDown moves the scroll position down.
func (d *DiffTab) ScrollDown(n int) {
	d.scrollPos += n
}

// ScrollUp moves the scroll position up.
func (d *DiffTab) ScrollUp(n int) {
	d.scrollPos -= n
	if d.scrollPos < 0 {
		d.scrollPos = 0
	}
}
