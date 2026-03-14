#!/bin/bash

# SDD Dashboard — Progress Pane (top-right)
# Shows frontier progress, task status, iteration count.
# Renders each frame into a buffer and flushes in one write — no flicker.

set -uo pipefail

R=$'\033[0m'
B=$'\033[1m'
D=$'\033[2m'
GR=$'\033[32m'
YL=$'\033[33m'
RD=$'\033[31m'
CY=$'\033[36m'
BGG=$'\033[42m'
BGY=$'\033[43m'
BGD=$'\033[100m'
EL=$'\033[K'

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

  emit "${B}${CY} SDD PROGRESS${R}"
  emit "${D}${hr}${R}"

  # ── Ralph Loop state ──
  local state_file=".claude/ralph-loop.local.md"
  local iteration="-" max_iter="-" started_at="" active=false

  if [[ -f "$state_file" ]]; then
    active=true
    iteration=$(grep -m1 '^iteration:' "$state_file" 2>/dev/null | awk '{print $2}' || echo "-")
    max_iter=$(grep -m1 '^max_iterations:' "$state_file" 2>/dev/null | awk '{print $2}' || echo "-")
    started_at=$(grep -m1 '^started_at:' "$state_file" 2>/dev/null | sed 's/started_at: *"\{0,1\}\([^"]*\)"\{0,1\}/\1/' || echo "")
  fi

  local elapsed="--:--"
  if [[ -n "$started_at" ]]; then
    local start_epoch
    start_epoch=$(date -jf "%Y-%m-%dT%H:%M:%SZ" "$started_at" +%s 2>/dev/null || date -d "$started_at" +%s 2>/dev/null || echo "")
    if [[ -n "$start_epoch" ]]; then
      local now_epoch=$(date +%s)
      local diff=$((now_epoch - start_epoch))
      elapsed="$((diff / 60))m $((diff % 60))s"
    fi
  fi

  if [[ "$active" == "true" ]]; then
    emit "  ${GR}ACTIVE${R}  iter ${B}${iteration}${R}/${max_iter}  ${D}${elapsed}${R}"
  else
    emit "  ${D}IDLE — no active loop${R}"
  fi
  emit ""

  # ── Find frontier ──
  local frontier=""
  for dir in context/frontiers context/plans; do
    [[ -d "$dir" ]] || continue
    frontier=$(find "$dir" -name "*frontier*" -type f 2>/dev/null | sort | head -1)
    [[ -n "$frontier" ]] && break
  done

  if [[ -z "$frontier" ]]; then
    emit "  ${D}No frontier found${R}"
    while ((line < max - 1)); do emit ""; done
    emit "${D}$(date +%H:%M:%S)${R}"
    printf '%b' "$buf"
    return
  fi

  # ── Parse tasks ──
  local total=0 done_count=0 in_progress=0 blocked=0
  local task_ids=()

  while IFS= read -r tline; do
    local tid
    tid=$(echo "$tline" | grep -oE 'T-[0-9]+' | head -1)
    [[ -n "$tid" ]] && task_ids+=("$tid") && total=$((total + 1))
  done < <(grep -E '^\s*-\s+T-[0-9]+' "$frontier" 2>/dev/null || true)

  declare -A task_status=()
  for impl_file in context/impl/impl-*.md; do
    [[ -f "$impl_file" ]] || continue
    while IFS= read -r tline; do
      local tid
      tid=$(echo "$tline" | grep -oE 'T-[0-9]+' | head -1)
      [[ -z "$tid" ]] && continue
      if echo "$tline" | grep -qi 'DONE'; then
        task_status["$tid"]="DONE"
      elif echo "$tline" | grep -qi 'IN.PROGRESS\|PARTIAL'; then
        task_status["$tid"]="WIP"
      elif echo "$tline" | grep -qi 'BLOCKED\|DEAD.END'; then
        task_status["$tid"]="BLOCKED"
      fi
    done < <(grep -E 'T-[0-9]+' "$impl_file" 2>/dev/null || true)
  done

  for tid in "${task_ids[@]+"${task_ids[@]}"}"; do
    case "${task_status[$tid]:-}" in
      DONE)    done_count=$((done_count + 1)) ;;
      WIP)     in_progress=$((in_progress + 1)) ;;
      BLOCKED) blocked=$((blocked + 1)) ;;
    esac
  done
  local remaining=$((total - done_count - in_progress - blocked))

  # ── Progress bar ──
  emit "${B} Tasks${R}"
  emit "${D}${hr}${R}"

  if [[ $total -gt 0 ]]; then
    local bar_w=$((cols - 8))
    [[ $bar_w -gt 36 ]] && bar_w=36
    local filled=$((done_count * bar_w / total))
    local wip_w=0
    [[ $in_progress -gt 0 ]] && wip_w=$((in_progress * bar_w / total))
    [[ $wip_w -lt 1 && $in_progress -gt 0 ]] && wip_w=1
    local empty=$((bar_w - filled - wip_w))
    [[ $empty -lt 0 ]] && empty=0

    local bar=""
    [[ $filled -gt 0 ]] && bar+=$(printf "${BGG}%*s${R}" "$filled" "")
    [[ $wip_w -gt 0 ]] && bar+=$(printf "${BGY}%*s${R}" "$wip_w" "")
    [[ $empty -gt 0 ]] && bar+=$(printf "${BGD}%*s${R}" "$empty" "")

    emit "  ${bar} $((done_count * 100 / total))%"
  fi

  emit ""
  emit "  ${GR}Done${R}        ${B}${done_count}${R}/${total}"
  [[ $in_progress -gt 0 ]] && emit "  ${YL}In Progress${R} ${B}${in_progress}${R}"
  [[ $blocked -gt 0 ]] && emit "  ${RD}Blocked${R}     ${B}${blocked}${R}"
  emit "  ${D}Remaining${R}   ${remaining}"
  emit ""

  # ── Tiers ──
  emit "${B} Tiers${R}"
  emit "${D}${hr}${R}"

  local cur_tier="" tier_total=0 tier_done=0
  while IFS= read -r tline; do
    if [[ "$tline" =~ ^##.*[Tt]ier ]]; then
      if [[ -n "$cur_tier" ]]; then
        local mk="${D}-${R}"
        [[ $tier_done -eq $tier_total && $tier_total -gt 0 ]] && mk="${GR}*${R}"
        [[ $tier_done -gt 0 && $tier_done -lt $tier_total ]] && mk="${YL}>${R}"
        emit "  ${mk} ${cur_tier}: ${tier_done}/${tier_total}"
      fi
      cur_tier=$(echo "$tline" | sed 's/^##[[:space:]]*//' | head -c $((cols - 14)))
      tier_total=0; tier_done=0
    elif [[ "$tline" =~ ^[[:space:]]*-[[:space:]]+T-[0-9]+ ]]; then
      local tid
      tid=$(echo "$tline" | grep -oE 'T-[0-9]+' | head -1)
      tier_total=$((tier_total + 1))
      [[ "${task_status[$tid]:-}" == "DONE" ]] && tier_done=$((tier_done + 1))
    fi
  done < "$frontier"

  if [[ -n "$cur_tier" ]]; then
    local mk="${D}-${R}"
    [[ $tier_done -eq $tier_total && $tier_total -gt 0 ]] && mk="${GR}*${R}"
    [[ $tier_done -gt 0 && $tier_done -lt $tier_total ]] && mk="${YL}>${R}"
    emit "  ${mk} ${cur_tier}: ${tier_done}/${tier_total}"
  fi
  emit ""

  # ── Dead ends ──
  local dead_ends=0
  for impl_file in context/impl/impl-*.md; do
    [[ -f "$impl_file" ]] || continue
    dead_ends=$((dead_ends + $(grep -ciE 'dead.end' "$impl_file" 2>/dev/null || echo 0)))
  done
  [[ $dead_ends -gt 0 ]] && emit "  ${RD}Dead ends: ${dead_ends}${R}" && emit ""

  # ── Adversarial ──
  local findings_file="context/impl/adversarial-findings.md"
  if [[ -f "$findings_file" ]]; then
    local critical=$(grep -ciE 'CRITICAL' "$findings_file" 2>/dev/null || echo 0)
    local high=$(grep -ciE '\bHIGH\b' "$findings_file" 2>/dev/null || echo 0)
    emit "${B} Adversarial${R}"
    emit "${D}${hr}${R}"
    [[ $critical -gt 0 ]] && emit "  ${RD}CRITICAL${R} ${critical}"
    [[ $high -gt 0 ]] && emit "  ${YL}HIGH${R}     ${high}"
    [[ $critical -eq 0 && $high -eq 0 ]] && emit "  ${GR}Clean${R}"
    emit ""
  fi

  # Blank remaining lines + footer
  while ((line < max - 1)); do emit ""; done
  emit "${D}$(date +%H:%M:%S)${R}"

  # Single atomic write
  printf '%b' "$buf"
}

while true; do
  render
  sleep 5
done
