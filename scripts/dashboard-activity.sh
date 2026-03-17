#!/bin/bash

# SDD Dashboard — Activity Pane (bottom-right)
# Shows recent loop iterations and git commits.
# Renders each frame into a buffer and flushes in one write — no flicker.

set -uo pipefail

R=$'\033[0m'
B=$'\033[1m'
D=$'\033[2m'
GR=$'\033[32m'
YL=$'\033[33m'
RD=$'\033[31m'
BL=$'\033[34m'
EL=$'\033[K'
TASK_ID_PATTERN='T-([A-Za-z0-9]+-)*[A-Za-z0-9]+'

tput civis 2>/dev/null
trap 'tput cnorm 2>/dev/null' EXIT
printf '\033[2J'

render() {
  local cols=$(tput cols 2>/dev/null || echo 40)
  local rows=$(tput lines 2>/dev/null || echo 24)
  local hr=$(printf '%.0s─' $(seq 1 "$cols"))
  local line=0
  local max=$((rows - 1))
  local buf=$'\033[H'

  emit() { ((line < max)) && buf+="${1}${EL}"$'\n' && line=$((line + 1)); }

  emit "${B}${BL} ACTIVITY${R}"
  emit "${BL}${D}${hr}${R}"

  # ── Recent iterations ──
  local log_file="context/impl/loop-log.md"

  if [[ -f "$log_file" ]]; then
    local max_entries=$(( (rows - 12) / 2 ))
    [[ $max_entries -lt 3 ]] && max_entries=3
    [[ $max_entries -gt 8 ]] && max_entries=8

    local all_iters=()
    local iter_buf="" in_iter=false

    while IFS= read -r fline; do
      if [[ "$fline" =~ ^###[[:space:]]+Iteration ]]; then
        [[ -n "$iter_buf" ]] && all_iters+=("$iter_buf")
        iter_buf="$fline"
        in_iter=true
      elif [[ "$in_iter" == "true" ]]; then
        if [[ "$fline" =~ ^###[[:space:]] ]] && ! [[ "$fline" =~ ^###[[:space:]]+Iteration ]]; then
          in_iter=false
        else
          iter_buf+=$'\n'"$fline"
        fi
      fi
    done < "$log_file"
    [[ -n "$iter_buf" ]] && all_iters+=("$iter_buf")

    local start_idx=$(( ${#all_iters[@]} - max_entries ))
    [[ $start_idx -lt 0 ]] && start_idx=0

    if [[ ${#all_iters[@]} -gt 0 ]]; then
      for ((i=start_idx; i<${#all_iters[@]}; i++)); do
        local block="${all_iters[$i]}"
        local iter_num=$(echo "$block" | grep -oE 'Iteration [0-9]+' | head -1 | awk '{print $2}')
        local task=$(echo "$block" | grep -m1 'Task:' | sed 's/.*Task:[[:space:]]*//' | head -c $((cols - 10)))
        local status=$(echo "$block" | grep -m1 'Status:' | sed 's/.*Status:[[:space:]]*//' | awk '{print $1}')

        local icon="${D}?${R}"
        case "$status" in
          DONE)    icon="${GR}*${R}" ;;
          PARTIAL) icon="${YL}~${R}" ;;
          BLOCKED) icon="${RD}x${R}" ;;
        esac

        emit "  ${icon} ${D}#${iter_num}${R} ${task}"
      done
    else
      emit "  ${D}No iterations yet${R}"
    fi
  else
    emit "  ${D}Waiting for first iteration...${R}"
  fi

  emit ""

  # ── Git commits (pre-fetched before rendering) ──
  emit "${B} Git${R}"
  emit "${BL}${D}${hr}${R}"

  local git_output=""
  if git rev-parse --git-dir &>/dev/null; then
    git_output=$(git log --oneline -5 2>/dev/null || true)
  fi

  if [[ -n "$git_output" ]]; then
    while IFS= read -r gline; do
      local hash=${gline%% *}
      local msg=${gline#* }
      emit "  ${D}${hash}${R} ${msg:0:$((cols - 12))}"
    done <<< "$git_output"
  else
    emit "  ${D}No commits${R}"
  fi

  emit ""

  # ── Current task ──
  local current_task=""
  for impl_file in context/impl/impl-*.md; do
    [[ -f "$impl_file" ]] || continue
    local wip=$(grep -iE 'IN.PROGRESS|PARTIAL' "$impl_file" 2>/dev/null | grep -oE "$TASK_ID_PATTERN" | tail -1)
    [[ -n "$wip" ]] && current_task="$wip"
  done

  [[ -n "$current_task" ]] && emit "  ${YL}>${R} Working on ${B}${current_task}${R}" && emit ""

  if [[ ! -f ".claude/ralph-loop.local.md" ]]; then
    emit "  ${D}Loop not active${R}"
  fi

  # Blank remaining lines
  while ((line < max - 1)); do emit ""; done
  emit "${D}$(date +%H:%M:%S)${R}"

  # Single atomic write
  printf '%b' "$buf"
}

while true; do
  render
  sleep 3
done
