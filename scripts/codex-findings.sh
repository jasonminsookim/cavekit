#!/usr/bin/env bash
# codex-findings.sh — Sourceable utility for managing review findings (T-007)
# Source this file to get bp_findings_* functions.
#
# Usage:
#   source scripts/codex-findings.sh
#   bp_findings_init
#   bp_findings_append "P0" "src/main.go:L42" "Nil pointer dereference" "codex-tier-gate" 2
#   bp_findings_list_blocking
#   bp_findings_update_status "F-003" "FIXED"

# Guard against double-sourcing
[[ -n "${_BP_FINDINGS_LOADED:-}" ]] && { return 0 2>/dev/null || true; }
_BP_FINDINGS_LOADED=1

# ── bp_findings_path ───────────────────────────────────────────────────

bp_findings_path() {
  local root
  root="$(git rev-parse --show-toplevel 2>/dev/null || echo ".")"
  echo "${root}/context/impl/impl-review-findings.md"
}

# ── bp_findings_init ───────────────────────────────────────────────────

bp_findings_init() {
  local fpath
  fpath="$(bp_findings_path)"

  if [[ -f "$fpath" ]]; then
    # Migrate old format (without Source/Tier) if detected
    if grep -q '| Finding | Severity | File | Status | Task |' "$fpath" 2>/dev/null; then
      sed -i.bak \
        -e 's/| Finding | Severity | File | Status | Task |/| Finding | Severity | File | Status | Source | Tier | Task |/' \
        -e 's/|---------|----------|------|--------|------|/|---------|----------|------|--------|--------|------|------|/' \
        "$fpath"
      rm -f "${fpath}.bak"
    fi
    return 0
  fi

  mkdir -p "$(dirname "$fpath")"
  local now
  now="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

  cat > "$fpath" << EOF
---
created: "${now}"
last_edited: "${now}"
---

# Review Findings

| Finding | Severity | File | Status | Source | Tier | Task |
|---------|----------|------|--------|--------|------|------|
EOF
}

# ── bp_findings_next_id ────────────────────────────────────────────────

bp_findings_next_id() {
  local fpath
  fpath="$(bp_findings_path)"

  if [[ ! -f "$fpath" ]]; then
    echo "F-001"
    return 0
  fi

  local max_num=0
  local num
  while IFS= read -r id; do
    num="${id#F-}"
    num="$((10#$num))"
    if (( num > max_num )); then max_num=$num; fi
  done < <(grep -oE 'F-[0-9]+' "$fpath" | sort -u)

  printf "F-%03d\n" $(( max_num + 1 ))
}

# ── bp_findings_append ─────────────────────────────────────────────────
#   $1 — severity (P0-P3)
#   $2 — file[:line]
#   $3 — description
#   $4 — source tag (e.g. "codex-tier-gate")
#   $5 — tier number
#   $6 — (optional) task reference, defaults to "—"

bp_findings_append() {
  local severity="${1:?severity required (P0-P3)}"
  local file="${2:?file required}"
  local description="${3:?description required}"
  local source="${4:?source required}"
  local tier="${5:?tier number required}"
  local task="${6:-—}"

  local fpath
  fpath="$(bp_findings_path)"

  bp_findings_init

  local fid
  fid="$(bp_findings_next_id)"

  echo "| ${fid}: ${description} | ${severity} | ${file} | NEW | ${source} | ${tier} | ${task} |" >> "$fpath"

  echo "$fid"
}

# ── bp_findings_update_status ──────────────────────────────────────────
#   $1 — finding ID (e.g. F-003)
#   $2 — new status (e.g. FIXED, WONTFIX)

bp_findings_update_status() {
  local finding_id="${1:?finding_id required}"
  local new_status="${2:?new_status required}"

  local fpath
  fpath="$(bp_findings_path)"

  if [[ ! -f "$fpath" ]]; then
    echo "ERROR: findings file not found" >&2
    return 1
  fi

  if ! grep -q "${finding_id}:" "$fpath"; then
    echo "ERROR: finding ${finding_id} not found" >&2
    return 1
  fi

  # Replace Status column (4th pipe-delimited field) for the matching row
  perl -i -pe '
    if (/^\|\s*'"${finding_id}"':/) {
      my @f = split /\|/;
      if (scalar @f >= 5) {
        $f[4] = " '"${new_status}"' ";
        $_ = join("|", @f) . "\n";
      }
    }
  ' "$fpath"
}

# ── bp_findings_list_blocking ──────────────────────────────────────────
# Lists P0/P1 findings with status NEW (one per line: "ID|SEVERITY|FILE")

bp_findings_list_blocking() {
  local fpath
  fpath="$(bp_findings_path)"

  [[ ! -f "$fpath" ]] && return 0

  grep -E '^\|' "$fpath" \
    | grep -vF '| Finding' \
    | grep -vE '^\|[-]' \
    | while IFS='|' read -r _ finding severity file status rest; do
        severity="$(echo "$severity" | xargs)"
        status="$(echo "$status" | xargs)"
        if [[ "$severity" == "P0" || "$severity" == "P1" ]] && [[ "$status" == "NEW" ]]; then
          finding="$(echo "$finding" | xargs)"
          file="$(echo "$file" | xargs)"
          echo "${finding}|${severity}|${file}"
        fi
      done
}

# ── CLI mode ───────────────────────────────────────────────────────────

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  set -euo pipefail
  cmd="${1:-help}"
  shift || true
  case "$cmd" in
    init) bp_findings_init ;;
    next-id) bp_findings_next_id ;;
    append) bp_findings_append "$@" ;;
    update) bp_findings_update_status "$@" ;;
    blocking) bp_findings_list_blocking ;;
    path) bp_findings_path ;;
    help|--help|-h)
      echo "Usage: codex-findings.sh {init|next-id|append|update|blocking|path}"
      ;;
    *) echo "Unknown command: $cmd" >&2; exit 1 ;;
  esac
fi
