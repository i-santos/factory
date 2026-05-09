#!/usr/bin/env bash
set -euo pipefail

script_name="provision-initiative-worktree.sh"

resolve_factory_skill_script() {
  local repo_root="$1"
  local candidates=()

  if [[ -n "${FACTORY_SKILL_ROOT:-}" ]]; then
    candidates+=("$FACTORY_SKILL_ROOT/scripts/$script_name")
  fi

  candidates+=(
    "$repo_root/.agents/skills/factory/scripts/$script_name"
    "$repo_root/.codex/skills/factory/scripts/$script_name"
    "${CODEX_HOME:-$HOME/.codex}/skills/factory/scripts/$script_name"
    "$HOME/.agents/skills/factory/scripts/$script_name"
  )

  local candidate
  for candidate in "${candidates[@]}"; do
    if [[ -x "$candidate" ]]; then
      printf '%s\n' "$candidate"
      return 0
    fi
  done

  return 1
}

repo_root="$(git rev-parse --show-toplevel)"

if canonical_script="$(resolve_factory_skill_script "$repo_root")"; then
  exec "$canonical_script" "$@"
fi

cat >&2 <<EOF
Unable to locate executable factory skill script: scripts/$script_name

Set FACTORY_SKILL_ROOT to the installed factory skill directory, or install the
factory skill in one of the supported locations:
- \$CODEX_HOME/skills/factory
- \$HOME/.codex/skills/factory
- \$HOME/.agents/skills/factory
- <repo>/.agents/skills/factory
EOF
exit 1
