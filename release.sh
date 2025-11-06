#!/bin/bash
# release.sh - Automated release script with conventional commit version bumping

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# Main release process
main() {
    print_header "Itamae Release Script"
    
    # Check if we're on master/main branch
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    if [ "$CURRENT_BRANCH" != "master" ] && [ "$CURRENT_BRANCH" != "main" ]; then
        print_error "Not on master/main branch (currently on: $CURRENT_BRANCH)"
        exit 1
    fi
    
    # Check for uncommitted changes
    print_info "Running pre-flight checks..."
    if ! git diff-index --quiet HEAD --; then
        print_error "You have uncommitted changes. Please commit or stash them first."
        git status --short
        exit 1
    fi
    print_success "Working tree is clean"
    
    # Run tests
    print_header "Running Tests"
    print_info "Formatting code..."
    go fmt ./...
    print_success "Code formatted"
    
    print_info "Running vet..."
    go vet ./...
    print_success "Vet passed"
    
    print_info "Running tests..."
    go test ./...
    print_success "All tests passed"
    
    # Get version information
    print_header "Version Analysis"
    
    # Get the latest tag or detect if no tags exist
    LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    
    if [ -z "$LATEST_TAG" ]; then
        print_info "No existing tags found"
        CURRENT_VERSION="0.0.0"
        COMMIT_RANGE="HEAD"
    else
        print_info "Latest tag: ${LATEST_TAG}"
        CURRENT_VERSION="${LATEST_TAG#v}"
        COMMIT_RANGE="${LATEST_TAG}..HEAD"
    fi
    
    print_info "Current version: ${CURRENT_VERSION}"
    
    # Analyze commits
    if [ -z "$LATEST_TAG" ]; then
        print_info "Analyzing all commits..."
    else
        print_info "Analyzing commits since ${LATEST_TAG}..."
    fi
    echo ""
    
    # Get commits
    COMMITS=$(git log ${COMMIT_RANGE} --pretty=format:"%s" --no-merges 2>/dev/null || echo "")
    
    if [ -z "$COMMITS" ]; then
        print_error "No commits found for analysis"
        exit 1
    fi
    
    # Display commits
    echo "$COMMITS" | while IFS= read -r commit; do
        echo "  - $commit"
    done
    echo ""
    
    # Count commit types
    BREAKING_COUNT=$(echo "$COMMITS" | grep -cE "^[a-z]+(\(.+\))?!:" || true)
    BREAKING_COUNT=$((BREAKING_COUNT + $(echo "$COMMITS" | grep -ci "BREAKING CHANGE:" || true)))
    FEAT_COUNT=$(echo "$COMMITS" | grep -cE "^feat(\(.+\))?:" || true)
    FIX_COUNT=$(echo "$COMMITS" | grep -cE "^fix(\(.+\))?:" || true)
    OTHER_COUNT=$(echo "$COMMITS" | grep -cvE "^(feat|fix)(\(.+\))?:" || true)
    
    print_info "Commit summary:"
    echo "  Breaking changes: $BREAKING_COUNT"
    echo "  Features:         $FEAT_COUNT"
    echo "  Fixes:            $FIX_COUNT"
    echo "  Other:            $OTHER_COUNT"
    echo ""
    
    # Determine version bump
    IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
    
    if [ $BREAKING_COUNT -gt 0 ]; then
        BUMP_TYPE="major"
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
    elif [ $FEAT_COUNT -gt 0 ]; then
        BUMP_TYPE="minor"
        MINOR=$((MINOR + 1))
        PATCH=0
    elif [ $FIX_COUNT -gt 0 ]; then
        BUMP_TYPE="patch"
        PATCH=$((PATCH + 1))
    else
        print_warning "No conventional commits found (feat/fix/BREAKING CHANGE)"
        print_info "Defaulting to patch bump"
        BUMP_TYPE="patch"
        PATCH=$((PATCH + 1))
    fi
    
    NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"
    
    print_success "Determined bump type: ${BUMP_TYPE}"
    print_success "New version: ${NEW_VERSION}"
    echo ""
    
    # Confirm with user
    print_warning "This will create and push tag ${NEW_VERSION}"
    read -p "Continue with release? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Release cancelled"
        exit 0
    fi
    
    # Build with new version
    print_header "Building Release"
    print_info "Building ${NEW_VERSION}..."
    VERSION=${NEW_VERSION} ./build.sh
    
    # Test the built binary
    print_info "Testing built binary..."
    ./bin/itamae version
    print_success "Binary built and tested successfully"
    
    # Create and push tag
    print_header "Creating Release"
    
    print_info "Creating git tag..."
    git tag -a "${NEW_VERSION}" -m "Release ${NEW_VERSION}"
    print_success "Tag created: ${NEW_VERSION}"
    
    # Push to origin
    print_info "Pushing to origin..."
    git push origin "${CURRENT_BRANCH}"
    print_success "Pushed branch to origin"
    
    print_info "Pushing tag to origin..."
    git push origin "${NEW_VERSION}"
    print_success "Pushed tag to origin"
    
    # Success message
    print_header "Release Complete!"
    echo ""
    print_success "Released ${NEW_VERSION}"
    echo ""
    print_info "GitHub Actions will now build the release artifacts."
    print_info "View the release at: https://github.com/yjmrobert/itamae/releases/tag/${NEW_VERSION}"
    echo ""
}

# Handle script interruption
trap 'print_error "Release interrupted"; exit 130' INT TERM

# Run main function
main "$@"
