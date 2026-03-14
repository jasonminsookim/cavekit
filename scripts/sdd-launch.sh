#!/bin/bash

# sdd-launch — One command to start the full SDD tmux dashboard + claude loop.
#
# Usage from any repo:
#   sdd-launch [OPTIONS]
#   sdd-launch --kill
#
# Creates a tmux session with:
#   Left (70%):    claude running /sdd-execute with your options
#   Top-right:     live progress dashboard
#   Bottom-right:  live activity feed
#
# All OPTIONS are forwarded to /sdd-execute (--filter, --adversarial, etc.)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SESSION_NAME="sdd-loop"
RIGHT_WIDTH_PCT=30

# ─── Kill mode ──────────────────────────────────────────────────────────────

if [[ "${1:-}" == "--kill" || "${1:-}" == "kill" ]]; then
  tmux kill-session -t "$SESSION_NAME" 2>/dev/null && echo "Session '$SESSION_NAME' killed." || echo "No active session."
  exit 0
fi

# ─── Preflight ──────────────────────────────────────────────────────────────

if ! command -v tmux &>/dev/null; then
  echo "tmux not found. Install: brew install tmux" >&2
  exit 1
fi

if ! command -v claude &>/dev/null; then
  echo "claude not found. Install Claude Code first." >&2
  exit 1
fi

# Kill existing session if one is running
if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
  echo "Existing '$SESSION_NAME' session found. Killing it..."
  tmux kill-session -t "$SESSION_NAME"
fi

# Collect all args to forward to /sdd-execute
SDD_ARGS="$*"

# ─── Detect real terminal size ──────────────────────────────────────────────

# stty gives the actual terminal size, tput can lie in subshells
if STTY_SIZE=$(stty size 2>/dev/null); then
  ROWS=$(echo "$STTY_SIZE" | awk '{print $1}')
  COLS=$(echo "$STTY_SIZE" | awk '{print $2}')
else
  # Fallback: use COLUMNS env var or default large
  COLS=${COLUMNS:-200}
  ROWS=${LINES:-50}
fi

WORK_DIR="$(pwd)"

# Calculate pane widths
RIGHT_WIDTH=$(( COLS * RIGHT_WIDTH_PCT / 100 ))
[[ "$RIGHT_WIDTH" -lt 35 ]] && RIGHT_WIDTH=35

echo "Starting SDD loop with dashboard (${COLS}x${ROWS})..."

# ─── Create tmux session ───────────────────────────────────────────────────

# Write a tiny launcher script that sends /sdd-execute after claude starts
LAUNCHER=$(mktemp /tmp/sdd-launch-XXXXXX.sh)
cat > "$LAUNCHER" <<LAUNCHER_EOF
#!/bin/bash
# Auto-cleanup
rm -f "$LAUNCHER"
# Start claude, then fall back to bash when it exits
cd "$WORK_DIR"
claude
LAUNCHER_EOF
chmod +x "$LAUNCHER"

# Create session — main pane runs claude
tmux new-session -d -s "$SESSION_NAME" -c "$WORK_DIR" -x "$COLS" -y "$ROWS" \
  "bash $LAUNCHER; exec bash"

# Wait a moment for the session to be ready, then send the command as keystrokes
(
  sleep 2
  tmux send-keys -t "$SESSION_NAME:0.0" "/sdd-execute $SDD_ARGS" Enter
) &

# Split right: progress pane (top-right)
tmux split-window -h -t "$SESSION_NAME" -l "$RIGHT_WIDTH" -c "$WORK_DIR" \
  "exec bash \"$SCRIPT_DIR/dashboard-progress.sh\""
tmux select-pane -T "sdd-progress"

# Split the right pane vertically: activity pane (bottom-right)
tmux split-window -v -t "$SESSION_NAME" -c "$WORK_DIR" \
  "exec bash \"$SCRIPT_DIR/dashboard-activity.sh\""
tmux select-pane -T "sdd-activity"

# Focus the main pane
tmux select-pane -t "$SESSION_NAME:0.0"

# ─── Attach ─────────────────────────────────────────────────────────────────

exec tmux attach-session -t "$SESSION_NAME"
