#!/usr/bin/env bash
# =============================================================================
# Bulk rename script: dockerscout -> dockerscout
#
# - Replaces text in file contents (ASCII/text files only)
# - Renames files and directories that contain dockerscout in their path
# - Handles common case variants:
#     dockerscout  -> dockerscout
#     DockerScout  -> DockerScout
#     DOCKERSCOUT  -> DOCKERSCOUT
#
# Usage:
#   ./scripts/rename-dockerscout-to-dockerscout.sh               # dry-run
#   ./scripts/rename-dockerscout-to-dockerscout.sh --apply       # apply changes
#   ./scripts/rename-dockerscout-to-dockerscout.sh --root .      # custom root
# =============================================================================

set -euo pipefail

ROOT="$(pwd)"
APPLY=0
SELF_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/$(basename "${BASH_SOURCE[0]}")"

usage() {
  cat <<'EOF'
rename-dockerscout-to-dockerscout.sh

Options:
  --apply         Apply changes (default is dry-run)
  --dry-run       Preview changes only (default)
  --root <path>   Project root to process (default: current directory)
  -h, --help      Show help
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply)
      APPLY=1
      shift
      ;;
    --dry-run)
      APPLY=0
      shift
      ;;
    --root)
      ROOT="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ ! -d "$ROOT" ]]; then
  echo "Root directory not found: $ROOT" >&2
  exit 1
fi

ROOT="$(cd "$ROOT" && pwd)"

echo "=== dockerscout -> dockerscout migration ==="
echo "Root: $ROOT"
if [[ "$APPLY" -eq 1 ]]; then
  echo "Mode: APPLY"
else
  echo "Mode: DRY-RUN"
fi

declare -i CONTENT_UPDATES=0
declare -i PATH_RENAMES=0

is_text_file() {
  local file="$1"
  grep -Iq . "$file"
}

new_name_from_old() {
  local old="$1"
  local new="$old"
  new="${new//dockerscout/dockerscout}"
  new="${new//DockerScout/DockerScout}"
  new="${new//DOCKERSCOUT/DOCKERSCOUT}"
  printf '%s' "$new"
}

replace_in_file() {
  local file="$1"
  local tmp
  tmp="$(mktemp)"

  perl -0pe 's/dockerscout/dockerscout/g; s/DockerScout/DockerScout/g; s/DOCKERSCOUT/DOCKERSCOUT/g' "$file" > "$tmp"

  if ! cmp -s "$file" "$tmp"; then
    if [[ "$APPLY" -eq 1 ]]; then
      mv "$tmp" "$file"
      echo "UPDATED content: $file"
    else
      rm -f "$tmp"
      echo "WOULD update content: $file"
    fi
    CONTENT_UPDATES+=1
  else
    rm -f "$tmp"
  fi
}

rename_path_if_needed() {
  local old_path="$1"

  if [[ "$old_path" == "$SELF_PATH" ]]; then
    return 0
  fi

  local new_path
  new_path="$(new_name_from_old "$old_path")"

  if [[ "$old_path" == "$new_path" ]]; then
    return 0
  fi

  if [[ "$APPLY" -eq 1 ]]; then
    if [[ -e "$new_path" ]]; then
      echo "SKIP rename (target exists): $old_path -> $new_path" >&2
      return 0
    fi
    mv "$old_path" "$new_path"
    echo "RENAMED path: $old_path -> $new_path"
  else
    echo "WOULD rename path: $old_path -> $new_path"
  fi

  PATH_RENAMES+=1
}

echo ""
echo "[1/2] Scanning and updating file contents..."
while IFS= read -r -d '' file; do
  case "$file" in
    */.git/*|*/node_modules/*|*/.idea/*|*/.vscode/*)
      continue
      ;;
  esac

  if is_text_file "$file"; then
    replace_in_file "$file"
  fi
done < <(find "$ROOT" -type f -print0)

echo ""
echo "[2/2] Renaming files/directories (deepest first)..."
while IFS= read -r -d '' path; do
  case "$path" in
    */.git/*)
      continue
      ;;
  esac
  rename_path_if_needed "$path"
done < <(find "$ROOT" -depth \( -name '*dockerscout*' -o -name '*DockerScout*' -o -name '*DOCKERSCOUT*' \) -print0)

echo ""
echo "=== Summary ==="
echo "Content updates: $CONTENT_UPDATES"
echo "Path renames:    $PATH_RENAMES"

if [[ "$APPLY" -eq 0 ]]; then
  echo ""
  echo "Dry-run complete. Re-run with --apply to perform changes."
fi
