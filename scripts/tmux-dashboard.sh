#!/bin/bash

# SDD Tmux Dashboard
# Creates a |= layout: main pane left, two info panes stacked right.
# If not in tmux, starts a new tmux session wrapping the current shell.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SESSION_NAME="sdd-loop"
RIGHT_WIDTH_PCT=30
MIN_COLS=100

# ─── Kill mode ──────────────────────────────────────────────────────────────

if [[ "${1:-}" == "--kill" || "${1:-}" == "kill" ]]; then
  # Kill dashboard panes if they exist
  if [[ -n "${TMUX:-}" ]]; then
    for pane_title in "sdd-progress" "sdd-activity"; do
      pane_id=$(tmux list-panes -F '#{pane_id} #{pane_title}' 2>/dev/null | grep "$pane_title" | awk '{print $1}' | head -1)
      [[ -n "$pane_id" ]] && tmux kill-pane -t "$pane_id" 2>/dev/null || true
    done
    echo "Dashboard panes closed."
  else
    # Kill the whole sdd-loop session
    tmux kill-session -t "$SESSION_NAME" 2>/dev/null && echo "Session '$SESSION_NAME' killed." || echo "No active session."
  fi
  exit 0
fi

# ─── Ensure tmux ────────────────────────────────────────────────────────────

if [[ -z "${TMUX:-}" ]]; then
  # Not in tmux — start a new session with the dashboard
  if ! command -v tmux &>/dev/null; then
    echo "tmux not found. Install: brew install tmux" >&2
    exit 1
  fi

  echo "Starting tmux session '$SESSION_NAME'..."

  # Create session with main pane running user's shell
  tmux new-session -d -s "$SESSION_NAME" -x "$(tput cols)" -y "$(tput lines)"

  # Now split from inside that session
  tmux send-keys -t "$SESSION_NAME" "\"$SCRIPT_DIR/tmux-dashboard.sh\" --attach" Enter

  # Attach
  exec tmux attach-session -t "$SESSION_NAME"
fi

# ─── Attach mode (called from inside the tmux session we just created) ────

if [[ "${1:-}" == "--attach" ]]; then
  # We're inside tmux now, fall through to create panes
  :
fi

# ─── Check terminal size ───────────────────────────────────────────────────

COLS=$(tmux display-message -p '#{window_width}')
if [[ "$COLS" -lt "$MIN_COLS" ]]; then
  echo "Terminal too narrow (${COLS} cols, need ${MIN_COLS}+). Skipping dashboard." >&2
  exit 0
fi

# ─── Create panes ──────────────────────────────────────────────────────────

MAIN_PANE=$(tmux display-message -p '#{pane_id}')

# Calculate right pane width
RIGHT_WIDTH=$(( COLS * RIGHT_WIDTH_PCT / 100 ))
[[ "$RIGHT_WIDTH" -lt 35 ]] && RIGHT_WIDTH=35

# Split horizontally: main left, progress right
PROGRESS_PANE=$(tmux split-window -h -l "$RIGHT_WIDTH" -P -F '#{pane_id}' \
  "exec bash \"$SCRIPT_DIR/dashboard-progress.sh\"")

# Split the right pane vertically: progress top, activity bottom
ACTIVITY_PANE=$(tmux split-window -v -t "$PROGRESS_PANE" -P -F '#{pane_id}' \
  "exec bash \"$SCRIPT_DIR/dashboard-activity.sh\"")

# Name panes for identification
tmux select-pane -t "$PROGRESS_PANE" -T "sdd-progress"
tmux select-pane -t "$ACTIVITY_PANE" -T "sdd-activity"

# Focus back on main pane
tmux select-pane -t "$MAIN_PANE"

echo "Dashboard active. Kill with: $SCRIPT_DIR/tmux-dashboard.sh --kill"
