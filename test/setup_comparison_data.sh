#!/bin/bash
# Setup Test Data for JSON Response Comparison
# This script creates sample data in both Go and Rails databases for testing

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if PostgreSQL is running
check_postgres() {
    if ! psql -h localhost -U "${DB_USER:-postgres}" -c "SELECT 1" > /dev/null 2>&1; then
        log_error "PostgreSQL is not running or not accessible"
        log_info "Start PostgreSQL with: brew services start postgresql (macOS) or sudo service postgresql start (Linux)"
        exit 1
    fi
    log_success "PostgreSQL is running"
}

# Create test database
create_test_database() {
    local db_name="reading_log_comparison"
    
    log_info "Creating test database: $db_name"
    
    # Drop if exists
    psql -h localhost -U "${DB_USER:-postgres}" -c "DROP DATABASE IF EXISTS $db_name;" > /dev/null 2>&1 || true
    
    # Create database
    psql -h localhost -U "${DB_USER:-postgres}" -c "CREATE DATABASE $db_name;" > /dev/null 2>&1
    
    log_success "Test database created: $db_name"
}

# Create schema
create_schema() {
    local db_name="reading_log_comparison"
    
    log_info "Creating schema in $db_name"
    
    psql -h localhost -U "${DB_USER:-postgres}" -d "$db_name" << 'EOF'
-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    total_page INT NOT NULL DEFAULT 0,
    started_at TIMESTAMP WITH TIME ZONE,
    page INT NOT NULL DEFAULT 0,
    reinicia BOOLEAN NOT NULL DEFAULT false,
    progress VARCHAR(255),
    status VARCHAR(255),
    logs_count INT DEFAULT 0,
    days_unread INT DEFAULT 0,
    median_day VARCHAR(255),
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Logs table
CREATE TABLE IF NOT EXISTS logs (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    data TIMESTAMP WITHOUT TIME ZONE,
    start_page INT NOT NULL DEFAULT 0,
    end_page INT NOT NULL DEFAULT 0,
    wday INT NOT NULL DEFAULT 0,
    note TEXT,
    text TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS index_logs_on_project_id ON logs(project_id);
CREATE INDEX IF NOT EXISTS index_logs_on_project_id_and_data_desc ON logs(project_id, data DESC);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Watsons table
CREATE TABLE IF NOT EXISTS watsons (
    id BIGSERIAL PRIMARY KEY,
    start_at TIMESTAMP WITH TIME ZONE,
    end_at TIMESTAMP WITH TIME ZONE,
    minutes INT,
    external_id VARCHAR(255),
    log_id BIGINT,
    project_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS index_watsons_on_log_id ON watsons(log_id);
CREATE INDEX IF NOT EXISTS index_watsons_on_project_id ON watsons(project_id);
EOF

    log_success "Schema created"
}

# Insert sample data
insert_sample_data() {
    local db_name="reading_log_comparison"
    
    log_info "Inserting sample data"
    
    psql -h localhost -U "${DB_USER:-postgres}" -d "$db_name" << 'EOF'
-- Insert projects with various states
INSERT INTO projects (name, total_page, started_at, page, reinicia) VALUES
    ('The Pragmatic Programmer', 352, NOW() - INTERVAL '30 days', 200, false),
    ('Clean Code', 464, NOW() - INTERVAL '45 days', 464, false),
    ('Design Patterns', 395, NOW() - INTERVAL '15 days', 50, false),
    ('Refactoring', 448, NOW() - INTERVAL '60 days', 0, false),
    ('The Mythical Man-Month', 322, NOW() - INTERVAL '10 days', 100, false);

-- Insert logs for projects
INSERT INTO logs (project_id, data, start_page, end_page, wday, note) VALUES
    -- The Pragmatic Programmer (project_id=1)
    (1, NOW() - INTERVAL '1 day', 175, 200, EXTRACT(DOW FROM (NOW() - INTERVAL '1 day')), 'Finished chapter 8'),
    (1, NOW() - INTERVAL '3 days', 150, 175, EXTRACT(DOW FROM (NOW() - INTERVAL '3 days')), 'Great insights on debugging'),
    (1, NOW() - INTERVAL '5 days', 125, 150, EXTRACT(DOW FROM (NOW() - INTERVAL '5 days')), 'Morning reading session'),
    (1, NOW() - INTERVAL '7 days', 100, 125, EXTRACT(DOW FROM (NOW() - INTERVAL '7 days')), 'Started new chapter'),
    (1, NOW() - INTERVAL '10 days', 75, 100, EXTRACT(DOW FROM (NOW() - INTERVAL '10 days')), 'Evening reading'),
    (1, NOW() - INTERVAL '14 days', 50, 75, EXTRACT(DOW FROM (NOW() - INTERVAL '14 days')), 'Quick 25 pages'),
    (1, NOW() - INTERVAL '20 days', 25, 50, EXTRACT(DOW FROM (NOW() - INTERVAL '20 days')), 'Weekend reading'),
    (1, NOW() - INTERVAL '25 days', 0, 25, EXTRACT(DOW FROM (NOW() - INTERVAL '25 days')), 'Started the book'),

    -- Clean Code (project_id=2) - Finished
    (2, NOW() - INTERVAL '2 days', 440, 464, EXTRACT(DOW FROM (NOW() - INTERVAL '2 days')), 'Final chapter'),
    (2, NOW() - INTERVAL '5 days', 400, 440, EXTRACT(DOW FROM (NOW() - INTERVAL '5 days')), 'Getting close'),
    (2, NOW() - INTERVAL '10 days', 350, 400, EXTRACT(DOW FROM (NOW() - INTERVAL '10 days')), 'Great content'),
    (2, NOW() - INTERVAL '20 days', 250, 350, EXTRACT(DOW FROM (NOW() - INTERVAL '20 days')), 'Mid book'),
    (2, NOW() - INTERVAL '30 days', 150, 250, EXTRACT(DOW FROM (NOW() - INTERVAL '30 days')), 'Learning a lot'),
    (2, NOW() - INTERVAL '40 days', 0, 150, EXTRACT(DOW FROM (NOW() - INTERVAL '40 days')), 'Started Clean Code'),

    -- Design Patterns (project_id=3) - Recently started
    (3, NOW() - INTERVAL '2 days', 30, 50, EXTRACT(DOW FROM (NOW() - INTERVAL '2 days')), 'Observer pattern'),
    (3, NOW() - INTERVAL '5 days', 0, 30, EXTRACT(DOW FROM (NOW() - INTERVAL '5 days')), 'Started Design Patterns'),

    -- The Mythical Man-Month (project_id=5) - Active
    (5, NOW() - INTERVAL '1 day', 85, 100, EXTRACT(DOW FROM (NOW() - INTERVAL '1 day')), 'Brooks Law is real'),
    (5, NOW() - INTERVAL '4 days', 60, 85, EXTRACT(DOW FROM (NOW() - INTERVAL '4 days')), 'Interesting historical context'),
    (5, NOW() - INTERVAL '8 days', 30, 60, EXTRACT(DOW FROM (NOW() - INTERVAL '8 days')), 'Classic book'),
    (5, NOW() - INTERVAL '10 days', 0, 30, EXTRACT(DOW FROM (NOW() - INTERVAL '10 days')), 'Just started');

-- Insert a user
INSERT INTO users (name, email) VALUES
    ('Test User', 'test@example.com');
EOF

    log_success "Sample data inserted"
}

# Verify data
verify_data() {
    local db_name="reading_log_comparison"
    
    log_info "Verifying data"
    
    local project_count
    local log_count
    local user_count
    
    project_count=$(psql -h localhost -U "${DB_USER:-postgres}" -d "$db_name" -t -c "SELECT COUNT(*) FROM projects;" | tr -d ' ')
    log_count=$(psql -h localhost -U "${DB_USER:-postgres}" -d "$db_name" -t -c "SELECT COUNT(*) FROM logs;" | tr -d ' ')
    user_count=$(psql -h localhost -U "${DB_USER:-postgres}" -d "$db_name" -t -c "SELECT COUNT(*) FROM users;" | tr -d ' ')
    
    log_success "Data verification:"
    log_success "  - Projects: $project_count"
    log_success "  - Logs: $log_count"
    log_success "  - Users: $user_count"
}

# Update environment file
update_env_file() {
    local db_name="reading_log_comparison"
    
    log_info "Creating .env.comparison file"
    
    cat > .env.comparison << EOF
# Environment for JSON Response Comparison Testing
DB_HOST=localhost
DB_PORT=5432
DB_USER=${DB_USER:-postgres}
DB_PASS=${DB_PASS:-}
DB_DATABASE=$db_name
SERVER_PORT=3000
LOG_LEVEL=info
LOG_FORMAT=text
EOF

    log_success "Created .env.comparison file"
    log_info "Use with: export $(cat .env.comparison | xargs) && go run ./cmd"
}

# Main
main() {
    log_info "========================================"
    log_info "  Setup Comparison Test Data"
    log_info "========================================"
    echo ""

    check_postgres
    echo ""

    create_test_database
    echo ""

    create_schema
    echo ""

    insert_sample_data
    echo ""

    verify_data
    echo ""

    update_env_file
    echo ""

    log_info "========================================"
    log_success "  Setup Complete!"
    log_info "========================================"
    echo ""
    log_info "Next steps:"
    log_info "  1. Update .env or use .env.comparison"
    log_info "  2. Start Go API: make run"
    log_info "  3. Start Rails API: cd rails-app && rails s -p 3001"
    log_info "  4. Run comparison: make compare-responses"
    echo ""
}

main "$@"
