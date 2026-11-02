#!/bin/bash
# run.sh - Exercise a real itamae installation inside a disposable container.
#
# Expected environment:
#   CATEGORY     core | essentials | unverified
#   PKG_MANAGER  nala | apt-get   (which APT front end itamae should end up using)
#   SRC          path to the repository checkout (default /src)
#   REMOVE       set to 0 to skip the removal pass (default 1)
#
# The itamae binary under test must already be on PATH. The script fails on the
# first unexpected condition and prints a summary of what was verified.

set -euo pipefail

CATEGORY="${CATEGORY:?CATEGORY is required}"
PKG_MANAGER="${PKG_MANAGER:-apt-get}"
SRC="${SRC:-/src}"
REMOVE="${REMOVE:-1}"
SCRIPT_DIR="${SRC}/itamae/scripts/${CATEGORY}"

# Packages the container itself (or later test steps) depend on; they are
# installed and verified but never removed.
PROTECTED_FROM_REMOVAL="apt-transport-https ca-certificates curl gnupg python3-full wget"

# Scripts installed by itamae land in these user-level directories; itamae puts
# them on PATH for its own checks, so the harness must see them too.
export PATH="$HOME/.local/bin:$HOME/.cargo/bin:$HOME/.atuin/bin:$PATH"

failures=()

section() {
    echo ""
    echo "================================================================"
    echo "  $1"
    echo "================================================================"
}

fail() {
    echo "❌ $1" >&2
    failures+=("$1")
}

section "Environment"
# shellcheck disable=SC1091
. /etc/os-release
echo "Distribution : ${PRETTY_NAME}"
echo "Category     : ${CATEGORY}"
echo "Pkg manager  : ${PKG_MANAGER}"
echo "User         : $(id)"
echo "itamae       : $(command -v itamae) ($(itamae version | head -n1))"

section "Preparing package manager permutation"
case "$PKG_MANAGER" in
    nala)
        if ! command -v nala &> /dev/null; then
            sudo apt-get update
            sudo apt-get install -y nala
        fi
        command -v nala
        ;;
    apt-get)
        if command -v nala &> /dev/null; then
            sudo apt-get purge -y nala
        fi
        ! command -v nala &> /dev/null
        ;;
    *)
        echo "Unknown PKG_MANAGER: $PKG_MANAGER" >&2
        exit 2
        ;;
esac

# Plugin IDs of the category, and the subset whose check() must pass afterwards
# (manual plugins only print instructions and are reported as skipped).
mapfile -t ALL_IDS < <(itamae list --category "$CATEGORY" --ids)
VERIFIABLE_IDS=()
for id in "${ALL_IDS[@]}"; do
    if ! grep -q "^# INSTALL_METHOD: manual" "${SCRIPT_DIR}/${id}.sh"; then
        VERIFIABLE_IDS+=("$id")
    fi
done
echo "Plugins (${#ALL_IDS[@]}): ${ALL_IDS[*]}"

section "Before: nothing from the category should be installed yet"
PRESENT_BEFORE=()
for id in "${VERIFIABLE_IDS[@]}"; do
    if bash "${SCRIPT_DIR}/${id}.sh" check &> /dev/null; then
        PRESENT_BEFORE+=("$id")
    fi
done
echo "Already present (base image): ${PRESENT_BEFORE[*]:-none}"

section "itamae install --category ${CATEGORY}"
INSTALL_ARGS=(install --category "$CATEGORY" --yes --plain
    --input "GIT_USER_NAME=Itamae CI" --input "GIT_USER_EMAIL=ci@example.com")
if [ "$CATEGORY" = "unverified" ]; then
    INSTALL_ARGS+=(--all)
fi
start=$(date +%s)
if itamae "${INSTALL_ARGS[@]}"; then
    echo "✅ itamae install exited 0 in $(( $(date +%s) - start ))s"
else
    fail "itamae install exited non-zero"
fi

section "itamae check --category ${CATEGORY}"
if itamae check --category "$CATEGORY" --plugins "$(IFS=,; echo "${VERIFIABLE_IDS[*]}")"; then
    echo "✅ every installable tool passes its check()"
else
    fail "itamae check reported missing tools"
fi

section "Package manager actually used"
LOG=$(ls -t /tmp/itamae-logs/*.log 2>/dev/null | head -n1 || true)
if [ -z "$LOG" ]; then
    fail "no debug log was written to /tmp/itamae-logs"
else
    echo "Debug log: $LOG ($(wc -l < "$LOG") lines)"
    PM_RE="\$ sudo env DEBIAN_FRONTEND=noninteractive ${PKG_MANAGER} \(install\|update\)"
    if grep -q "$PM_RE" "$LOG"; then
        echo "✅ ${PKG_MANAGER} was used for package operations"
    elif grep -q "\$ sudo env DEBIAN_FRONTEND=noninteractive \(nala\|apt-get\) install" "$LOG"; then
        fail "expected ${PKG_MANAGER} to be used, log shows: $(grep -o '\$ sudo env DEBIAN_FRONTEND=noninteractive [a-z-]* \(install\|update\)' "$LOG" | sort -u | tr '\n' ' ')"
    else
        echo "ℹ no APT operations in this category"
    fi
    if grep -q "ERROR:" "$LOG"; then
        echo "Errors recorded in the debug log:"
        grep "ERROR:" "$LOG" | head -n 20
    fi
fi

if [ "$REMOVE" = "1" ]; then
    section "Removal pass (reverse order)"
    for (( i=${#ALL_IDS[@]}-1; i>=0; i-- )); do
        id="${ALL_IDS[$i]}"
        if [[ " ${PROTECTED_FROM_REMOVAL} " == *" ${id} "* ]]; then
            echo "⏭  ${id}: protected, not removed"
            continue
        fi
        echo "--- remove ${id}"
        if ! bash "${SCRIPT_DIR}/${id}.sh" remove; then
            fail "${id}: remove exited non-zero"
            continue
        fi
        if [[ " ${PRESENT_BEFORE[*]} " == *" ${id} "* ]]; then
            echo "ℹ  ${id}: was part of the base image, not asserting removal"
            continue
        fi
        if bash "${SCRIPT_DIR}/${id}.sh" check &> /dev/null; then
            fail "${id}: check() still passes after remove"
        else
            echo "✅ ${id}: removed"
        fi
    done
fi

section "Summary"
if [ ${#failures[@]} -eq 0 ]; then
    echo "✅ ${CATEGORY} on ${PRETTY_NAME} with ${PKG_MANAGER}: all checks passed"
    exit 0
fi
echo "❌ ${#failures[@]} failure(s):"
for f in "${failures[@]}"; do
    echo "   - $f"
done
exit 1
