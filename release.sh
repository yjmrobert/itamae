#!/bin/bash
# release.sh - Automated release script with conventional commit version bumping

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
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

# Check if we're on the master branch
check_branch() {
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    if [ "$CURRENT_BRANCH" != "master" ] && [ "$CURRENT_BRANCH" != "main" ]; then
        print_error "Not on master/main branch (currently on: $CURRENT_BRANCH)"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Check for uncommitted changes
check_clean_working_tree() {
    if ! git diff-index --quiet HEAD --; then
        print_error "You have uncommitted changes. Please commit or stash them first."
        git status --short
        exit 1
    fi
}

# Get the latest tag
get_latest_tag() {
    # Get the latest tag, fallback to v0.0.0 if none exists
    LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    echo "$LATEST_TAG"
}

# Parse semantic version
parse_version() {
    local version=$1
    # Remove 'v' prefix if present
    version=${version#v}
    
    # Split into major.minor.patch
    IFS='.' read -r MAJOR MINOR PATCH <<< "$version"
    
    # Remove any pre-release or metadata suffixes from patch
    PATCH=${PATCH%%-*}
    PATCH=${PATCH%%+*}
    
    echo "$MAJOR" "$MINOR" "$PATCH"
}

# Bump version based on type
bump_version() {
    local bump_type=$1
    local major=$2
    local minor=$3
    local patch=$4
    
    case $bump_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            print_error "Invalid bump type: $bump_type"
            exit 1
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# Analyze commits since last tag using conventional commits
analyze_commits() {
    local last_tag=$1
    local has_breaking=false
    local has_feat=false
    local has_fix=false
    
    print_info "Analyzing commits since ${last_tag}..."
    echo ""
    
    # Get commits since last tag
    if [ "$last_tag" = "v0.0.0" ]; then
        COMMITS=$(git log --pretty=format:"%s" --no-merges)
    else
        COMMITS=$(git log ${last_tag}..HEAD --pretty=format:"%s" --no-merges)
    fi
    
    if [ -z "$COMMITS" ]; then
        print_warning "No commits found since ${last_tag}"
        return 1
    fi
    
    # Count commit types
    local breaking_count=0
    local feat_count=0
    local fix_count=0
    local other_count=0
    
    while IFS= read -r commit; do
        echo "  - $commit"
        
        # Check for breaking changes (BREAKING CHANGE in body or ! after type)
        if echo "$commit" | grep -qiE "(BREAKING CHANGE|^[a-z]+\!:)"; then
            has_breaking=true
            breaking_count=$((breaking_count + 1))
        # Check for features
        elif echo "$commit" | grep -qiE "^feat(\([^)]*\))?:"; then
            has_feat=true
            feat_count=$((feat_count + 1))
        # Check for fixes
        elif echo "$commit" | grep -qiE "^fix(\([^)]*\))?:"; then
            has_fix=true
            fix_count=$((fix_count + 1))
        else
            other_count=$((other_count + 1))
        fi
    done <<< "$COMMITS"
    
    echo ""
    print_info "Commit summary:"
    echo "  Breaking changes: $breaking_count"
    echo "  Features:         $feat_count"
    echo "  Fixes:            $fix_count"
    echo "  Other:            $other_count"
    echo ""
    
    # Determine bump type
    if [ "$has_breaking" = true ]; then
        echo "major"
    elif [ "$has_feat" = true ]; then
        echo "minor"
    elif [ "$has_fix" = true ]; then
        echo "patch"
    else
        echo "patch"  # Default to patch for other changes
    fi
}

# Main release process
main() {
    print_header "Itamae Release Script"
    
    # Pre-flight checks
    print_info "Running pre-flight checks..."
    check_branch
    check_clean_working_tree
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
    LATEST_TAG=$(get_latest_tag)
    print_info "Latest tag: ${LATEST_TAG}"
    
    # Parse current version
    read -r MAJOR MINOR PATCH <<< "$(parse_version "$LATEST_TAG")"
    print_info "Current version: ${MAJOR}.${MINOR}.${PATCH}"
    
    # Analyze commits
    BUMP_TYPE=$(analyze_commits "$LATEST_TAG")
    if [ $? -ne 0 ]; then
        print_error "No commits to release"
        exit 1
    fi
    
    # Calculate new version
    NEW_VERSION=$(bump_version "$BUMP_TYPE" "$MAJOR" "$MINOR" "$PATCH")
    print_success "Determined bump type: ${BUMP_TYPE}"
    print_success "New version: ${NEW_VERSION}"
    
    # Confirm with user
    echo ""
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
    
    # Generate commit message for tag
    TAG_MESSAGE="Release ${NEW_VERSION}"
    
    print_info "Creating git tag..."
    git tag -a "${NEW_VERSION}" -m "${TAG_MESSAGE}"
    print_success "Tag created: ${NEW_VERSION}"
    
    # Push to origin
    print_info "Pushing to origin..."
    git push origin "${CURRENT_BRANCH:-master}"
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
