#!/bin/bash
# Test Redis Rate Limiter
# This script tests the rate limiting functionality of the API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="http://localhost:8080"
API_KEY="dev_api_key_123"
TEST_ENDPOINT="/api/v1/countries"
RATE_LIMIT_ENDPOINT="/api/v1/rate-limit/info"
MAX_REQUESTS=10
DELAY=0.1

# Logging function
log() { echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"; }
warn() { echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"; }
error() { echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"; }
info() { echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"; }

# Check if API is running
check_api() {
    log "Checking if API is running..."
    if ! curl -s -f "${API_BASE_URL}/health" > /dev/null; then
        error "API is not running at ${API_BASE_URL}"
        error "Please start the API first: make run"
        exit 1
    fi
    log "API is running ✓"
}

# Test rate limit info endpoint
test_rate_limit_info() {
    log "Testing rate limit info endpoint..."
    
    response=$(curl -s -H "Authorization: Bearer ${API_KEY}" \
        "${API_BASE_URL}${RATE_LIMIT_ENDPOINT}")
    
    if [[ $? -eq 0 ]]; then
        log "Rate limit info endpoint working ✓"
        echo "Response: $response" | jq '.' 2>/dev/null || echo "Response: $response"
    else
        error "Rate limit info endpoint failed"
        exit 1
    fi
}

# Test rate limiting by making multiple requests
test_rate_limiting() {
    log "Testing rate limiting with ${MAX_REQUESTS} requests..."
    
    success_count=0
    rate_limited_count=0
    
    for i in $(seq 1 $MAX_REQUESTS); do
        response=$(curl -s -w "%{http_code}" -H "Authorization: Bearer ${API_KEY}" \
            "${API_BASE_URL}${TEST_ENDPOINT}")
        
        http_code="${response: -3}"
        response_body="${response%???}"
        
        if [[ $http_code -eq 200 ]]; then
            success_count=$((success_count + 1))
            info "Request $i: SUCCESS (200)"
        elif [[ $http_code -eq 429 ]]; then
            rate_limited_count=$((rate_limited_count + 1))
            warn "Request $i: RATE LIMITED (429)"
            echo "Rate limit response: $response_body" | jq '.' 2>/dev/null || echo "Rate limit response: $response_body"
        else
            error "Request $i: UNEXPECTED ($http_code)"
            echo "Response: $response_body"
        fi
        
        sleep $DELAY
    done
    
    log "Rate limiting test completed:"
    log "  - Successful requests: $success_count"
    log "  - Rate limited requests: $rate_limited_count"
    
    if [[ $rate_limited_count -gt 0 ]]; then
        log "✓ Rate limiting is working correctly"
    else
        warn "Rate limiting may not be working (no 429 responses)"
    fi
}

# Test rate limit headers
test_rate_limit_headers() {
    log "Testing rate limit headers..."
    
    headers=$(curl -s -I -H "Authorization: Bearer ${API_KEY}" \
        "${API_BASE_URL}${TEST_ENDPOINT}")
    
    echo "$headers" | grep -E "(X-RateLimit|RateLimit)" || warn "No rate limit headers found"
    
    if echo "$headers" | grep -q "X-RateLimit-Limit"; then
        log "✓ Rate limit headers are present"
    else
        warn "Rate limit headers may not be configured"
    fi
}

# Test rate limit configuration
test_rate_limit_config() {
    log "Testing rate limit configuration endpoint..."
    
    response=$(curl -s -H "Authorization: Bearer ${API_KEY}" \
        "${API_BASE_URL}/api/v1/rate-limit/config")
    
    if [[ $? -eq 0 ]]; then
        log "Rate limit configuration endpoint working ✓"
        echo "Configuration:" | jq '.' 2>/dev/null || echo "Configuration: $response"
    else
        error "Rate limit configuration endpoint failed"
    fi
}

# Test rate limit stats
test_rate_limit_stats() {
    log "Testing rate limit statistics endpoint..."
    
    response=$(curl -s -H "Authorization: Bearer ${API_KEY}" \
        "${API_BASE_URL}/api/v1/rate-limit/stats")
    
    if [[ $? -eq 0 ]]; then
        log "Rate limit statistics endpoint working ✓"
        echo "Statistics:" | jq '.' 2>/dev/null || echo "Statistics: $response"
    else
        error "Rate limit statistics endpoint failed"
    fi
}

# Main test function
run_tests() {
    log "Starting Redis Rate Limiter tests..."
    
    check_api
    test_rate_limit_info
    test_rate_limit_config
    test_rate_limit_stats
    test_rate_limit_headers
    test_rate_limiting
    
    log "🎉 All rate limiter tests completed!"
    log "Check the output above for any warnings or errors."
}

# Handle script arguments
case "${1:-}" in
    test) run_tests ;;
    info) check_api; test_rate_limit_info ;;
    config) check_api; test_rate_limit_config ;;
    stats) check_api; test_rate_limit_stats ;;
    headers) check_api; test_rate_limit_headers ;;
    limit) check_api; test_rate_limiting ;;
    help|--help|-h) 
        echo "Usage: $0 [command]"
        echo "Commands:"
        echo "  test    - Run all tests"
        echo "  info    - Test rate limit info endpoint"
        echo "  config  - Test rate limit config endpoint"
        echo "  stats   - Test rate limit stats endpoint"
        echo "  headers - Test rate limit headers"
        echo "  limit   - Test rate limiting behavior"
        echo "  help    - Show this help message"
        ;;
    *) run_tests ;;
esac
