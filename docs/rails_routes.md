# Rails API Routes Report

Generated: 2026-04-08
Source: `docker exec reading-log-rails-api bundle exec rails routes`

## Overview

This document provides a comprehensive report of all Rails API routes defined in the application.

## Route Summary

| HTTP Verb | Path | Controller#Action |
|-----------|------|-------------------|
| GET | /v1/projects/:project_id/logs | v1/logs#index |
| GET | /v1/projects | v1/projects#index |
| GET | /v1/projects/:id | v1/projects#show |
| GET | /v1/users | v1/users#index |
| POST | /v1/users | v1/users#create |
| GET | /v1/users/:id | v1/users#show |
| PATCH | /v1/users/:id | v1/users#update |
| PUT | /v1/users/:id | v1/users#update |
| DELETE | /v1/users/:id | v1/users#destroy |
| GET | /v1/dashboard/day | v1/dashboard/day#index |
| GET | /v1/dashboard/projects | v1/dashboard/projects#index |
| GET | /v1/dashboard/last_days | v1/dashboard/last_days#index |
| GET | /v1/dashboard/echart/speculate_actual | v1/dashboard/echart/speculate_actual#index |
| GET | /v1/dashboard/echart/faults_week_day | v1/dashboard/echart/faults_week_day#index |
| GET | /v1/dashboard/echart/mean_progress | v1/dashboard/echart/mean_progress#index |
| GET | /v1/dashboard/echart/day_week | v1/dashboard/echart/day_week#index |
| GET | /v1/dashboard/echart/faults | v1/dashboard/echart/faults#index |
| GET | /v1/dashboard/echart/total | v1/dashboard/echart/total#index |
| GET | /v1/dashboard/echart/last_year_total | v1/dashboard/echart/last_year_total#index |
| GET | /v1/dashboard/echart/books_per_month | v1/dashboard/echart/books_per_month#index |
| GET | /v1/dashboard/echart/monthly_trend | v1/dashboard/echart/monthly_trend#index |

## API Endpoints by Namespace

### v1 Namespace

#### Projects Endpoints

| HTTP Verb | Path | Controller#Action | Description |
|-----------|------|-------------------|-------------|
| GET | /v1/projects/:project_id/logs | v1/logs#index | Get logs for a specific project |
| GET | /v1/projects | v1/projects#index | List all projects |
| GET | /v1/projects/:id | v1/projects#show | Get a specific project by ID |

#### Users Endpoints

| HTTP Verb | Path | Controller#Action | Description |
|-----------|------|-------------------|-------------|
| GET | /v1/users | v1/users#index | List all users |
| POST | /v1/users | v1/users#create | Create a new user |
| GET | /v1/users/:id | v1/users#show | Get a specific user by ID |
| PATCH | /v1/users/:id | v1/users#update | Update a user (partial) |
| PUT | /v1/users/:id | v1/users#update | Update a user (full) |
| DELETE | /v1/users/:id | v1/users#destroy | Delete a user |

### Dashboard Endpoints

#### Day Dashboard

| HTTP Verb | Path | Controller#Action | Description |
|-----------|------|-------------------|-------------|
| GET | /v1/dashboard/day | v1/dashboard/day#index | Get day-wise dashboard data |

#### Projects Dashboard

| HTTP Verb | Path | Controller#Action | Description |
|-----------|------|-------------------|-------------|
| GET | /v1/dashboard/projects | v1/dashboard/projects#index | Get dashboard data for projects |

#### Last Days Dashboard

| HTTP Verb | Path | Controller#Action | Description |
|-----------|------|-------------------|-------------|
| GET | /v1/dashboard/last_days | v1/dashboard/last_days#index | Get data for last days |

### ECharts Dashboard Endpoints

| HTTP Verb | Path | Controller#Action | Description |
|-----------|------|-------------------|-------------|
| GET | /v1/dashboard/echart/speculate_actual | v1/dashboard/echart/speculate_actual#index | Speculate vs actual chart data |
| GET | /v1/dashboard/echart/faults_week_day | v1/dashboard/echart/faults_week_day#index | Faults by week day chart data |
| GET | /v1/dashboard/echart/mean_progress | v1/dashboard/echart/mean_progress#index | Mean progress chart data |
| GET | /v1/dashboard/echart/day_week | v1/dashboard/echart/day_week#index | Day of week chart data |
| GET | /v1/dashboard/echart/faults | v1/dashboard/echart/faults#index | Faults chart data |
| GET | /v1/dashboard/echart/total | v1/dashboard/echart/total#index | Total chart data |
| GET | /v1/dashboard/echart/last_year_total | v1/dashboard/echart/last_year_total#index | Last year total chart data |
| GET | /v1/dashboard/echart/books_per_month | v1/dashboard/echart/books_per_month#index | Books per month chart data |
| GET | /v1/dashboard/echart/monthly_trend | v1/dashboard/echart/monthly_trend#index | Monthly trend chart data |

## Controller Summary

| Controller | Actions |
|------------|---------|
| v1/logs | index |
| v1/projects | index, show |
| v1/users | index, create, show, update, destroy |
| v1/dashboard/day | index |
| v1/dashboard/projects | index |
| v1/dashboard/last_days | index |
| v1/dashboard/echart/speculate_actual | index |
| v1/dashboard/echart/faults_week_day | index |
| v1/dashboard/echart/mean_progress | index |
| v1/dashboard/echart/day_week | index |
| v1/dashboard/echart/faults | index |
| v1/dashboard/echart/total | index |
| v1/dashboard/echart/last_year_total | index |
| v1/dashboard/echart/books_per_month | index |
| v1/dashboard/echart/monthly_trend | index |

## Route Count

- **Total Routes**: 20
- **v1/projects**: 3 routes
- **v1/users**: 5 routes
- **v1/dashboard**: 9 routes
- **v1/dashboard/echart**: 9 routes

## Notes

- All routes are under the `/v1` API namespace
- The `/v1/dashboard/echart` endpoints provide data for various ECharts visualizations
- User endpoints support full CRUD operations (Create, Read, Update, Delete)
- Logs are accessed through project associations (nested route)
