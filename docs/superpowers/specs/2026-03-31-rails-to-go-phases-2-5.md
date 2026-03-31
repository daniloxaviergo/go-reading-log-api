# Rails to Go Migration - Phases 2-5 Design

**Date**: 2026-03-31  
**Status**: Design Phase

---

## 1. Overview

This document outlines Phases 2-5 of the Rails to Go migration. Each phase builds on Phase 1 (core projects/logs API) to achieve full feature parity with the original Rails application.

---

## 2. Phase 2: Dashboard API

### 2.1 Purpose

Migrate all dashboard and echart endpoints for analytics and visualization.

### 2.2 Endpoints to Migrate

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/dashboard/day` | Daily dashboard statistics |
| GET | `/api/v1/dashboard/projects` | Project dashboard data |
| GET | `/api/v1/dashboard/last_days` | Last days dashboard data |
| GET | `/api/v1/dashboard/echart/faults` | Fault percentage (gauge chart) |
| GET | `/api/v1/dashboard/echart/speculate_actual` | Speculative vs actual data |
| GET | `/api/v1/dashboard/echart/faults_week_day` | Faults by weekday |
| GET | `/api/v1/dashboard/echart/mean_progress` | Mean progress over 30 days |
| GET | `/api/v1/dashboard/echart/day_week` | Day/week data |
| GET | `/api/v1/dashboard/echart/last_year_total` | Last year total by week |
| GET | `/api/v1/dashboard/echart/total` | Total data |
| GET | `/api/v1/dashboard/echart/books_per_month` | Books per month |
| GET | `/api/v1/dashboard/echart/monthly_trend` | Monthly trend line chart |

### 2.3 Analytics Logic

The dashboard endpoints calculate various metrics:

| Endpoint | Calculation |
|----------|-------------|
| faults | Fault count over last 30 days, percentage of max fault count |
| speculate_actual | Predicted vs actual reading progress |
| faults_week_day | Faults grouped by weekday (last 6 months) |
| mean_progress | Average daily reading over last 30 days |
| last_year_total | Total logs grouped by week (last year) |

### 2.4 Dependencies

- Phase 1 (projects/logs endpoints)
- Business logic from Rails `classes/v1/dashboard/` directory

---

## 3. Phase 3: Users Endpoint

### 3.1 Purpose

Add user management endpoint with optional authentication.

### 3.2 Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/users` | List all users |
| GET | `/api/v1/users/:id` | Get user by ID |
| POST | `/api/v1/users` | Create user |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |

### 3.3 User Schema

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | Primary key |
| name | string | User's name |
| email | string | User's email |

### 3.4 Authentication Options

| Option | Description |
|--------|-------------|
| A. JWT | Token-based auth with refresh tokens |
| B. Session | Cookie-based auth (matches Rails) |
| C. None | Open API (no auth) |

---

## 4. Phase 4: Background Jobs

### 4.1 Purpose

Restore background job processing previously handled by Sidekiq.

### 4.2 Job Types (from Rails)

| Job | Description |
|-----|-------------|
| Data sync | Sync with external services |
| Notifications | Send email/Slack notifications |
| Analytics | Pre-compute dashboard metrics |

### 4.3 Integration Options

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| NATS | Lightweight messaging | Simple, fast | Less ecosystem |
| NSQ | Distributed messaging | Scalable, reliable | More complex |
| RabbitMQ | Traditional broker | Mature, well-known | Heavier |
| AWS SQS | Cloud-based | Managed, scalable | AWS lock-in |

### 4.4 Implementation Notes

- Go worker processes jobs from queue
- Compatible with Sidekiq JSON format (optional)
- Retry logic with exponential backoff

---

## 5. Phase 5: Watson Tracking

### 5.1 Purpose

Add Watson time-tracking endpoint.

### 5.2 Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/watsons` | List all watsons |
| GET | `/api/v1/watsons/:id` | Get watson by ID |
| POST | `/api/v1/watsons` | Create watson |
| PUT | `/api/v1/watsons/:id` | Update watson |
| DELETE | `/api/v1/watsons/:id` | Delete watson |
| GET | `/api/v1/projects/:id/watsons` | Get watsons for project |

### 5.3 Watson Schema

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | Primary key |
| start_at | time.Time | Start time |
| end_at | time.Time | End time |
| minutes | int | Duration in minutes |
| external_id | string | External reference ID |
| log_id | int64 | Foreign key to logs |
| project_id | int64 | Foreign key to projects |

### 5.4 Dependencies

- Phase 1 (projects/logs)
- Phase 3 (if watsons are user-scoped)

---

## 7. CLI Tasks Migration

### 7.1 Purpose

Migrate Rails Rake tasks to Go CLI commands for data synchronization and batch processing.

### 7.2 Tasks to Migrate

| Task | Description |
|------|-------------|
| `dxupdate` | Import/update projects from `/leituras.txt` YAML file |
| `dxupdate_watson` | Import Watson frames from `/watson_frames` YAML file |
| `dxreading_day` | Output reading statistics by weekday |
| `dxlinear_data` | Generate linear regression dataset to file |
| `dxholtwinters` | Generate Holt-Winters forecasting dataset to file |
| `neural_dataset` | Generate neural network training dataset to CSV |

### 7.3 CLI Command Structure

```bash
# Data synchronization
go-reading-log-api dxupdate --file /leituras.txt
go-reading-log-api dxupdate-watson --file /watson_frames

# Data export
go-reading-log-api dxreading-day
go-reading-log-api dxlinear-data --output linear_data
go-reading-log-api dxholtwinters --output holtwinters
go-reading-log-api neural-dataset --output neural_dataset.csv
```

### 7.4 Key Logic

**dxupdate:**
- Load projects from YAML file
- Filter logs by last update timestamp
- Create new projects or update existing
- Only import new/updated logs

**dxupdate_watson:**
- Load Watson frames from YAML file
- Map Watson project names to database projects
- Link watsons to logs based on time proximity
- Skip duplicates

**dxreading_day:**
- Group logs by weekday
- Output statistics for current day

**dxlinear_data:**
- Generate datasets for linear regression
- Pages by day, number of days

**dxholtwinters:**
- Calculate reading frequency patterns
- Generate time-series data for forecasting

**neural_dataset:**
- Generate CSV for neural network training
- Format: `project_id, day_number, read_in_day`

### 7.5 Implementation Notes

- Use `github.com/spf13/cobra` for CLI framework
- Reuse repository layer from API
- Use `gopkg.in/gomail.v2` for email (if needed)
- File I/O with proper error handling

---

## 9. Migration Order Recommendation

```
- Phase 1: Core API (projects, logs) ✓
  ↓
- Phase 2: Dashboard API (analytics)
  ↓
- Phase 3: Users endpoint (optional auth)
  ↓
- Phase 4: Background jobs (message queue)
  ↓
- Phase 5: Watson tracking (CRUD)
  ↓
- Phase 6: CLI tasks (batch processing)
```

---

## 10. Success Criteria for Full Migration

Full migration complete when:
1. All Phase 1-5 endpoints return same JSON as Rails app
2. All tests pass
3. Performance meets or exceeds Rails baseline
4. Documentation complete
5. Zero Rails app dependencies (can decommission)

---

## 11. Notes

- Phases can be implemented in any order based on priority
- Dashboard analytics (Phase 2) is most complex due to calculations
- Background jobs (Phase 4) is optional if sync operations acceptable
- Users endpoint (Phase 3) can be skipped if not needed
- CLI tasks (Phase 6) are independent of API phases
