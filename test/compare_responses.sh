#!/bin/bash
# JSON Response Comparison Script
# Compares JSON responses between Go and Rails APIs for all three endpoints
#
# Usage:
#   ./compare_responses.sh [options]
#
# Options:
#   -g, --go-url     Go API base URL (default: http://localhost:3000/api/v1)
#   -r, --rails-url  Rails API base URL (default: http://localhost:3001/api/v1)
#   -h, --help       Show help message
#
# Requirements:
#   - curl
#   - jq (version 1.6+ for deep comparison)
#   - Both Go and Rails APIs must be running

set -euo pipefail

# Configuration
GO_API_URL="${GO_API_URL:-http://localhost:3000/api/v1}"
RAILS_API_URL="${RAILS_API_URL:-http://localhost:3001/api/v1}"
TIMEOUT=10
TEMP_DIR=""
ENDPOINT_SUFFIX=""  # No .json suffix needed - APIs use plain /api/v1/... routes

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Statistics
TESTS_PASSED=0
TESTS_FAILED=0
ENDPOINTS_TESTED=0

# Cleanup function
cleanup() {
    if [[ -n "$TEMP_DIR" && -d "$TEMP_DIR" ]]; then
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
}

# Show help message
show_help() {
    cat << EOF
JSON Response Comparison Script

Compares JSON responses between Go and Rails APIs to verify identical structure
and values for all three endpoints.

Usage: $(basename "$0") [options]

Options:
    -g, --go-url GO_URL      Go API base URL (default: http://localhost:3000/api/v1)
    -r, --rails-url URL      Rails API base URL (default: http://localhost:3001/api/v1)
    -h, --help               Show this help message

Examples:
    ./compare_responses.sh
    ./compare_responses.sh -g http://localhost:8080/api/v1 -r http://localhost:3001/api/v1
    GO_API_URL=http://localhost:3000/api/v1 RAILS_API_URL=http://localhost:3001/api/v1 ./compare_responses.sh

Requirements:
    - curl
    - jq (version 1.6+ recommended)
    - Both Go and Rails APIs must be running

Endpoints tested:
    1. GET /api/v1/projects.json
    2. GET /api/v1/projects/{id}.json
    3. GET /api/v1/projects/{id}/logs.json

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -g|--go-url)
                GO_API_URL="$2"
                shift 2
                ;;
            -r|--rails-url)
                RAILS_API_URL="$2"
                shift 2
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Check required commands
check_requirements() {
    local missing=0

    for cmd in curl jq; do
        if ! command -v "$cmd" &> /dev/null; then
            log_error "$cmd is required but not installed"
            missing=1
        fi
    done

    # Check jq version for slurp support
    if command -v jq &> /dev/null; then
        jq_version=$(jq --version | sed 's/jq-//')
        jq_major=$(echo "$jq_version" | cut -d. -f1)
        jq_minor=$(echo "$jq_version" | cut -d. -f2)
        if [[ "$jq_major" -lt 1 || ( "$jq_major" -eq 1 && "$jq_minor" -lt 6 ) ]]; then
            log_warning "jq version $jq_version may not support all features. Recommended: 1.6+"
        fi
    fi

    if [[ $missing -eq 1 ]]; then
        exit 1
    fi
}

# Check if APIs are accessible
check_apis_accessible() {
    log_info "Checking API accessibility..."

    # Check Go API
    if ! curl -s --max-time "$TIMEOUT" "${GO_API_URL}/projects${ENDPOINT_SUFFIX}" > /dev/null 2>&1; then
        log_error "Go API not accessible at ${GO_API_URL}"
        log_info "Make sure the Go API is running with: make run"
        return 1
    fi

    # Check Rails API
    if ! curl -s --max-time "$TIMEOUT" "${RAILS_API_URL}/projects${ENDPOINT_SUFFIX}" > /dev/null 2>&1; then
        log_error "Rails API not accessible at ${RAILS_API_URL}"
        log_info "Make sure the Rails API is running on port 3001"
        return 1
    fi

    log_success "Both APIs are accessible"
    return 0
}

# Fetch JSON from an endpoint
fetch_json() {
    local url="$1"
    local output_file="$2"
    local suffix="${3:-}"

    # Create parent directory if it doesn't exist
    local parent_dir
    parent_dir=$(dirname "$output_file")
    mkdir -p "$parent_dir"

    curl -s --max-time "$TIMEOUT" -H "Accept: application/json" "${url}${suffix}" > "$output_file"
}

# Normalize JSON for comparison (sort keys, remove whitespace)
normalize_json() {
    local input_file="$1"
    local output_file="$2"

    jq -S '.' "$input_file" > "$output_file"
}

# Compare two JSON files for structural equality
compare_json_structures() {
    local file1="$1"
    local file2="$2"

    # Normalize both JSON files
    local norm1="$TEMP_DIR/normalized_1.json"
    local norm2="$TEMP_DIR/normalized_2.json"

    normalize_json "$file1" "$norm1"
    normalize_json "$file2" "$norm2"

    # Compare normalized JSON
    if diff -q "$norm1" "$norm2" > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Compare JSON values with tolerance for floating point numbers
compare_json_values() {
    local file1="$1"
    local file2="$2"
    local tolerance="${3:-0.01}"

    # Use jq to compare values with tolerance for floating point numbers
    jq -n --argjson tol "$tolerance" \
        --slurpfile j1 "$file1" \
        --slurpfile j2 "$file2" '
        def compare_values(val1; val2):
            if (val1 | type) == "object" and (val2 | type) == "object" then
                (val1 | keys) as $k1 |
                (val2 | keys) as $k2 |
                ($k1 - $k2 | length == 0) and ($k2 - $k1 | length == 0) and
                ($k1 | all(. as $k | compare_values(val1[$k]; val2[$k])))
            elif (val1 | type) == "array" and (val2 | type) == "array" then
                (val1 | length) == (val2 | length) and
                ([range(0; val1 | length)] | all(. as $i | compare_values(val1[$i]; val2[$i])))
            elif (val1 | type) == "number" and (val2 | type) == "number" then
                ((val1 - val2) | fabs) < $tolerance
            elif (val1 | type) == "string" and (val2 | type) == "string" then
                val1 == val2
            elif (val1 | type) == "null" and (val2 | type) == "null" then
                true
            elif (val1 | type) == "boolean" and (val2 | type) == "boolean" then
                val1 == val2
            else
                false
            end;

        compare_values($j1[0]; $j2[0])
    '
}

# Extract project ID from response for subsequent requests
get_first_project_id() {
    local response_file="$1"

    jq -r '.[0].id // empty' "$response_file"
}

# Get projects with logs from Rails API (for comparison reference)
get_rails_project_with_logs() {
    local project_id="$1"
    local output_file="$2"

    curl -s --max-time "$TIMEOUT" -H "Accept: application/json" \
        "${RAILS_API_URL}/projects/${project_id}${ENDPOINT_SUFFIX}" > "$output_file"
}

# Test the index endpoint (/api/v1/projects)
test_index_endpoint() {
    local name="Index Endpoint (GET /api/v1/projects)"
    local temp_file="$TEMP_DIR/index"
    local go_file="$temp_file/go.json"
    local rails_file="$temp_file/rails.json"

    log_info "Testing: $name"

    # Fetch from both APIs
    fetch_json "$GO_API_URL/projects" "$go_file" "$ENDPOINT_SUFFIX"
    fetch_json "$RAILS_API_URL/projects" "$rails_file" "$ENDPOINT_SUFFIX"

    # Check if we got valid JSON
    if ! jq empty "$go_file" 2>/dev/null; then
        log_error "Invalid JSON from Go API"
        return 1
    fi

    if ! jq empty "$rails_file" 2>/dev/null; then
        log_error "Invalid JSON from Rails API"
        return 1
    fi

    # Compare structures
    if compare_json_structures "$go_file" "$rails_file"; then
        log_success "$name - Structures match"
        ((TESTS_PASSED++))
    else
        log_error "$name - Structures differ"
        # Show a diff of normalized JSON
        local norm1="$TEMP_DIR/diff_norm1.json"
        local norm2="$TEMP_DIR/diff_norm2.json"
        normalize_json "$go_file" "$norm1"
        normalize_json "$rails_file" "$norm2"
        diff -u "$norm1" "$norm2" || true
        ((TESTS_FAILED++))
        return 1
    fi

    # Compare values with tolerance
    if compare_json_values "$go_file" "$rails_file" 0.01 | grep -q "true"; then
        log_success "$name - Values match within tolerance"
        ((TESTS_PASSED++))
    else
        log_error "$name - Values differ"
        ((TESTS_FAILED++))
        return 1
    fi

    ((ENDPOINTS_TESTED++))
    return 0
}

# Test the show endpoint (/api/v1/projects/:id)
test_show_endpoint() {
    local name="Show Endpoint (GET /api/v1/projects/:id)"
    local temp_file="$TEMP_DIR/show"
    local go_file="$temp_file/go.json"
    local rails_file="$temp_file/rails.json"

    log_info "Testing: $name"

    # Get first project ID from index response
    local project_id
    project_id=$(get_first_project_id "$TEMP_DIR/index/rails.json")

    if [[ -z "$project_id" ]]; then
        log_warning "No projects found in database, skipping show endpoint test"
        return 0
    fi

    # Fetch from both APIs
    fetch_json "$GO_API_URL/projects/$project_id" "$go_file" "$ENDPOINT_SUFFIX"
    fetch_json "$RAILS_API_URL/projects/$project_id" "$rails_file" "$ENDPOINT_SUFFIX"

    # Check if we got valid JSON
    if ! jq empty "$go_file" 2>/dev/null; then
        log_error "Invalid JSON from Go API"
        return 1
    fi

    if ! jq empty "$rails_file" 2>/dev/null; then
        log_error "Invalid JSON from Rails API"
        return 1
    fi

    # Compare structures
    if compare_json_structures "$go_file" "$rails_file"; then
        log_success "$name - Structures match"
        ((TESTS_PASSED++))
    else
        log_error "$name - Structures differ"
        local norm1="$TEMP_DIR/diff_norm3.json"
        local norm2="$TEMP_DIR/diff_norm4.json"
        normalize_json "$go_file" "$norm1"
        normalize_json "$rails_file" "$norm2"
        diff -u "$norm1" "$norm2" || true
        ((TESTS_FAILED++))
        return 1
    fi

    # Compare values with tolerance
    if compare_json_values "$go_file" "$rails_file" 0.01 | grep -q "true"; then
        log_success "$name - Values match within tolerance"
        ((TESTS_PASSED++))
    else
        log_error "$name - Values differ"
        ((TESTS_FAILED++))
        return 1
    fi

    ((ENDPOINTS_TESTED++))
    return 0
}

# Test the logs endpoint (/api/v1/projects/:id/logs)
test_logs_endpoint() {
    local name="Logs Endpoint (GET /api/v1/projects/:id/logs)"
    local temp_file="$TEMP_DIR/logs"
    local go_file="$temp_file/go.json"
    local rails_file="$temp_file/rails.json"

    log_info "Testing: $name"

    # Get first project ID
    local project_id
    project_id=$(get_first_project_id "$TEMP_DIR/index/rails.json")

    if [[ -z "$project_id" ]]; then
        log_warning "No projects found in database, skipping logs endpoint test"
        return 0
    fi

    # Fetch from both APIs
    fetch_json "$GO_API_URL/projects/$project_id/logs" "$go_file" "$ENDPOINT_SUFFIX"
    fetch_json "$RAILS_API_URL/projects/$project_id/logs" "$rails_file" "$ENDPOINT_SUFFIX"

    # Check if we got valid JSON
    if ! jq empty "$go_file" 2>/dev/null; then
        log_error "Invalid JSON from Go API"
        return 1
    fi

    if ! jq empty "$rails_file" 2>/dev/null; then
        log_error "Invalid JSON from Rails API"
        return 1
    fi

    # Compare structures
    if compare_json_structures "$go_file" "$rails_file"; then
        log_success "$name - Structures match"
        ((TESTS_PASSED++))
    else
        log_error "$name - Structures differ"
        local norm1="$TEMP_DIR/diff_norm5.json"
        local norm2="$TEMP_DIR/diff_norm6.json"
        normalize_json "$go_file" "$norm1"
        normalize_json "$rails_file" "$norm2"
        diff -u "$norm1" "$norm2" || true
        ((TESTS_FAILED++))
        return 1
    fi

    # Compare values with tolerance
    if compare_json_values "$go_file" "$rails_file" 0.01 | grep -q "true"; then
        log_success "$name - Values match within tolerance"
        ((TESTS_PASSED++))
    else
        log_error "$name - Values differ"
        ((TESTS_FAILED++))
        return 1
    fi

    ((ENDPOINTS_TESTED++))
    return 0
}

# Test edge cases
test_edge_cases() {
    local name="Edge Cases"
    local temp_file="$TEMP_DIR/edges"
    local go_file="$temp_file/go.json"
    local rails_file="$temp_file/rails.json"

    log_info "Testing: $name"

    # Test 1: Empty logs scenario
    log_info "  Testing empty logs scenario..."

    # Create a project without logs in test database
    local project_id
    project_id=$(get_first_project_id "$TEMP_DIR/index/rails.json")

    if [[ -n "$project_id" ]]; then
        # Get project data from both APIs
        fetch_json "$GO_API_URL/projects/$project_id" "$go_file"
        fetch_json "$RAILS_API_URL/projects/$project_id" "$rails_file"

        # Check logs_count is present and consistent
        local go_logs_count
        local rails_logs_count
        go_logs_count=$(jq '.logs_count // 0' "$go_file")
        rails_logs_count=$(jq '.logs_count // 0' "$rails_file")

        if [[ "$go_logs_count" == "$rails_logs_count" ]]; then
            log_success "  logs_count field consistent: $go_logs_count"
            ((TESTS_PASSED++))
        else
            log_error "  logs_count mismatch: Go=$go_logs_count, Rails=$rails_logs_count"
            ((TESTS_FAILED++))
        fi

        # Check progress field
        local go_progress
        local rails_progress
        go_progress=$(jq '.progress // "null"' "$go_file")
        rails_progress=$(jq '.progress // "null"' "$rails_file")

        if [[ "$go_progress" == "$rails_progress" ]]; then
            log_success "  progress field consistent: $go_progress"
            ((TESTS_PASSED++))
        else
            log_error "  progress mismatch: Go=$go_progress, Rails=$rails_progress"
            ((TESTS_FAILED++))
        fi
    fi

    # Test 2: Null date handling
    log_info "  Testing null date handling..."

    # Check if started_at field is properly handled (null vs omitted)
    if [[ -n "$project_id" ]]; then
        local go_started_at
        local rails_started_at
        go_started_at=$(jq '.started_at' "$go_file")
        rails_started_at=$(jq '.started_at' "$rails_file")

        # Both should be either null or a date string
        local go_type=$(jq -r '.started_at | type' "$go_file")
        local rails_type=$(jq -r '.started_at | type' "$rails_file")

        if [[ "$go_type" == "$rails_type" ]]; then
            log_success "  started_at type consistent: $go_type"
            ((TESTS_PASSED++))
        else
            log_error "  started_at type mismatch: Go=$go_type, Rails=$rails_type"
            ((TESTS_FAILED++))
        fi
    fi

    ((ENDPOINTS_TESTED++))
    return 0
}

# Generate summary report
generate_report() {
    echo ""
    echo "========================================"
    echo "       JSON Response Comparison Report"
    echo "========================================"
    echo ""
    echo "Endpoints tested: $ENDPOINTS_TESTED"
    echo "Tests passed:     $TESTS_PASSED"
    echo "Tests failed:     $TESTS_FAILED"
    echo ""
    echo "API URLs:"
    echo "  Go API:      $GO_API_URL"
    echo "  Rails API:   $RAILS_API_URL"
    echo ""

    if [[ $TESTS_FAILED -eq 0 ]]; then
        echo -e "${GREEN}All tests passed!${NC}"
        return 0
    else
        echo -e "${RED}Some tests failed. See details above.${NC}"
        return 1
    fi
}

# Main execution
main() {
    log_info "JSON Response Comparison Script"
    log_info "==============================="
    echo ""

    # Parse command line arguments
    parse_args "$@"

    # Check requirements
    check_requirements

    # Check if APIs are accessible
    if ! check_apis_accessible; then
        exit 1
    fi

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    log_info "Using temporary directory: $TEMP_DIR"
    echo ""

    # Run tests
    log_info "Running comparison tests..."
    echo ""

    # Run each test and collect results
    local all_passed=true

    if ! test_index_endpoint; then
        all_passed=false
    fi

    if ! test_show_endpoint; then
        all_passed=false
    fi

    if ! test_logs_endpoint; then
        all_passed=false
    fi

    if ! test_edge_cases; then
        all_passed=false
    fi

    # Generate report
    generate_report

    if [[ "$all_passed" == "true" ]]; then
        exit 0
    else
        exit 1
    fi
}

# Run main function
main "$@"
