#!/bin/bash
# test.sh - End-to-end tests for install.sh.
#
# A local HTTP server stands in for GitHub releases so every code path of the
# installer can be exercised deterministically:
#   fresh install, already up to date, upgrade, downgrade, unknown installed
#   version, ITAMAE_VERSION pinning, ITAMAE_SKIP_RUN and ITAMAE_INSTALL_ARGS.
#
# Requirements: go, python3, sudo (passwordless) and network access for the
# final step, which runs a real "itamae install" of a small plugin.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
WORK="$(mktemp -d)"
MIRROR="${WORK}/mirror"
PORT="${PORT:-18080}"
BIN_DIR="${WORK}/bin"
export PATH="${BIN_DIR}:${PATH}"

cleanup() {
    if [ -n "${SERVER_PID:-}" ]; then
        kill "$SERVER_PID" 2>/dev/null || true
    fi
    rm -rf "$WORK"
}
trap cleanup EXIT

passed=0
failed=0

pass() { echo "✅ $1"; passed=$((passed + 1)); }
fail() { echo "❌ $1" >&2; failed=$((failed + 1)); }

assert_contains() {
    local haystack="$1" needle="$2" what="$3"
    if [[ "$haystack" == *"$needle"* ]]; then
        pass "$what"
    else
        fail "$what (expected output to contain: $needle)"
        echo "----- output -----"
        echo "$haystack"
        echo "------------------"
    fi
}

assert_version() {
    local want="$1" what="$2"
    local got
    got=$("${BIN_DIR}/itamae" version 2>/dev/null | head -n1 || true)
    if [[ "$got" == "itamae version ${want}" ]]; then
        pass "$what (installed ${want})"
    else
        fail "$what (expected ${want}, got: ${got})"
    fi
}

build_version() {
    local version="$1" out="$2"
    (cd "$ROOT" && go build -ldflags "-X github.com/yjmrobert/itamae/cmd.Version=${version} -X github.com/yjmrobert/itamae/cmd.GitCommit=test" -o "$out" .)
}

echo "=== Building fixture binaries"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$(uname -m)" in
    x86_64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
esac
ASSET="itamae-${OS}-${ARCH}"
mkdir -p "${MIRROR}/releases/download/v9.9.8" "${MIRROR}/releases/download/v9.9.9" "${MIRROR}/releases/latest/download" "${MIRROR}/api" "$BIN_DIR"
build_version v9.9.8 "${MIRROR}/releases/download/v9.9.8/${ASSET}"
build_version v9.9.9 "${MIRROR}/releases/download/v9.9.9/${ASSET}"
cp "${MIRROR}/releases/download/v9.9.9/${ASSET}" "${MIRROR}/releases/latest/download/${ASSET}"
echo '{"tag_name": "v9.9.9"}' > "${MIRROR}/api/latest.json"

echo "=== Starting local release mirror on port ${PORT}"
(cd "$MIRROR" && python3 -m http.server "$PORT" --bind 127.0.0.1 > "${WORK}/server.log" 2>&1) &
SERVER_PID=$!
for _ in $(seq 1 50); do
    curl -fs "http://127.0.0.1:${PORT}/api/latest.json" > /dev/null 2>&1 && break
    sleep 0.1
done
curl -fs "http://127.0.0.1:${PORT}/api/latest.json" > /dev/null

run_installer() {
    # Everything except network access to GitHub is redirected to the mirror.
    env ITAMAE_RELEASE_URL="http://127.0.0.1:${PORT}/releases" \
        ITAMAE_API_URL="http://127.0.0.1:${PORT}/api/latest.json" \
        ITAMAE_INSTALL_DIR="$BIN_DIR" \
        ITAMAE_SKIP_PREREQS=1 \
        "$@" bash "${ROOT}/install.sh" 2>&1
}

echo ""
echo "=== 1. Fresh install of the latest release"
out=$(run_installer ITAMAE_SKIP_RUN=1)
assert_contains "$out" "itamae not found - starting installation" "fresh install path taken"
assert_contains "$out" "Latest available version: v9.9.9" "latest version resolved from the API"
assert_version v9.9.9 "latest release installed"

echo ""
echo "=== 2. Re-run when already on the latest version"
out=$(run_installer ITAMAE_SKIP_RUN=1)
assert_contains "$out" "Already running the target version (v9.9.9)" "up-to-date path taken"
assert_contains "$out" "ITAMAE_SKIP_RUN is set" "ITAMAE_SKIP_RUN honoured"
assert_version v9.9.9 "binary untouched"

echo ""
echo "=== 3. Downgrade with ITAMAE_VERSION"
out=$(run_installer ITAMAE_SKIP_RUN=1 ITAMAE_VERSION=v9.9.8)
assert_contains "$out" "downgrading" "downgrade path taken"
assert_version v9.9.8 "pinned version installed"

echo ""
echo "=== 4. Upgrade from an older installed version"
out=$(run_installer ITAMAE_SKIP_RUN=1)
assert_contains "$out" "Newer version available: v9.9.8 → v9.9.9" "upgrade path taken"
assert_version v9.9.9 "upgraded to latest"

echo ""
echo "=== 5. Installed binary with an unparsable version is replaced"
printf '#!/bin/sh\necho "itamae but broken"\n' > "${BIN_DIR}/itamae"
chmod +x "${BIN_DIR}/itamae"
out=$(run_installer ITAMAE_SKIP_RUN=1)
assert_contains "$out" "Cannot determine installed version; forcing update" "unknown version forces update"
assert_version v9.9.9 "broken binary replaced"

echo ""
echo "=== 6. Download failure leaves the existing binary in place"
out=$(run_installer ITAMAE_SKIP_RUN=1 ITAMAE_VERSION=v0.0.404 || true)
assert_contains "$out" "Failed to download itamae" "download failure reported"
assert_version v9.9.9 "previous binary restored after failed download"

echo ""
echo "=== 7. Non-interactive run without ITAMAE_INSTALL_ARGS prints guidance"
out=$(run_installer < /dev/null)
assert_contains "$out" "Run 'itamae install' yourself" "installer explains how to run itamae without a terminal"

echo ""
echo "=== 8. ITAMAE_INSTALL_ARGS runs a real headless install"
out=$(run_installer ITAMAE_INSTALL_ARGS="--category essentials --plugins stow --yes --plain" < /dev/null)
assert_contains "$out" "Running: itamae install --category essentials --plugins stow --yes --plain" "install args forwarded"
assert_contains "$out" "All packages installed successfully" "stow installed through itamae"
if command -v stow > /dev/null; then
    pass "stow is on PATH after the run"
else
    fail "stow should be installed"
fi

echo ""
echo "=== 9. Unsupported architecture is rejected"
out=$(env ITAMAE_SKIP_PREREQS=1 bash -c 'uname() { if [ "$1" = -m ]; then echo mips; else command uname "$@"; fi; }; export -f uname; bash '"${ROOT}/install.sh" 2>&1 || true)
assert_contains "$out" "Unsupported architecture: mips" "unsupported architecture rejected"

echo ""
echo "=== 10. Prerequisite step installs git (and nala where available)"
out=$(env ITAMAE_RELEASE_URL="http://127.0.0.1:${PORT}/releases" ITAMAE_API_URL="http://127.0.0.1:${PORT}/api/latest.json" ITAMAE_INSTALL_DIR="$BIN_DIR" ITAMAE_SKIP_RUN=1 bash "${ROOT}/install.sh" 2>&1)
assert_contains "$out" "Prerequisites installed" "prerequisite installation completed"
if command -v git > /dev/null; then
    pass "git available after prerequisites"
else
    fail "git should be installed by the prerequisite step"
fi

echo ""
echo "=== Results: ${passed} passed, ${failed} failed"
[ "$failed" -eq 0 ]
