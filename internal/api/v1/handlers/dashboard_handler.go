package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
)

// DashboardHandler handles HTTP requests for dashboard endpoints
type DashboardHandler struct {
	repo            repository.DashboardRepository
	userConfig      *service.UserConfigService
	projectsService dashboard.ProjectsServiceInterface
}

// NewDashboardHandler creates a new DashboardHandler with the given repository, config service, and projects service
func NewDashboardHandler(repo repository.DashboardRepository, userConfig *service.UserConfigService, projectsService dashboard.ProjectsServiceInterface) *DashboardHandler {
	return &DashboardHandler{
		repo:            repo,
		userConfig:      userConfig,
		projectsService: projectsService,
	}
}

// =============================================================================
// Dashboard Endpoints
// =============================================================================

// Day handles GET /v1/dashboard/day.json - Returns daily statistics with weekday breakdown
func (h *DashboardHandler) Day(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get date from query parameter or use today's date
	dateStr := r.URL.Query().Get("date")
	var targetDate time.Time
	var err error

	if dateStr != "" {
		targetDate, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			http.Error(w, `{"error": "invalid date format", "details": {"date": "must be in RFC3339 format"}}`, http.StatusBadRequest)
			return
		}
	} else {
		targetDate = time.Now()
	}

	// Get daily stats from repository
	stats, err := h.repo.GetDailyStats(ctx, targetDate)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Calculate derived fields using StatsData
	statsData := &dto.StatsData{
		TotalPages: stats.TotalPages,
	}

	// Calculate mean_day using V1::MeanLog algorithm (7-day intervals)
	// Get mean by weekday using the new repository method
	meanDay, meanDayErr := h.repo.GetMeanByWeekday(ctx, int(targetDate.Weekday()))
	if meanDayErr == nil && meanDay != nil {
		statsData.MeanDay = *meanDay
	} else {
		statsData.MeanDay = 0.0
	}

	// Calculate progress_geral (from all projects)
	// Uses sum of project Page field divided by sum of project total_page (capacity)
	// This matches the Rails calculation: (sum of project.page) / (sum of project.total_page) * 100
	aggregates, aggErr := h.repo.GetProjectAggregates(ctx)
	if aggErr == nil && len(aggregates) > 0 {
		var totalCapacity int
		var totalPage int
		for _, agg := range aggregates {
			totalCapacity += agg.TotalPage // Project's total_page field
			// Query project's Page field from database for accurate progress calculation
			var projectPage int
			pageQuery := `SELECT COALESCE(page, 0) FROM projects WHERE id = $1`
			pageErr := h.repo.GetPool().QueryRow(ctx, pageQuery, agg.ProjectID).Scan(&projectPage)
			if pageErr == nil {
				totalPage += projectPage
			}
		}
		if totalCapacity > 0 {
			statsData.ProgressGeral = math.Round(float64(totalPage)/float64(totalCapacity)*100*1000) / 1000
		} else {
			statsData.ProgressGeral = 0.0
		}
	} else {
		statsData.ProgressGeral = 0.0
	}

	// Calculate per_pages (ratio of last week to previous week)
	// For today's data, we use the current day's pages as "last week" and calculate ratio
	// Get previous period data for comparison
	prevStart := targetDate.AddDate(0, 0, -7)

	prevStats, prevErr := h.repo.GetDailyStats(ctx, prevStart)
	if prevErr == nil && prevStats.TotalPages > 0 {
		perPages := math.Round(float64(stats.TotalPages)/float64(prevStats.TotalPages)*1000) / 1000
		statsData.PerPages = &perPages
	} else {
		// Return null when no previous data available (e.g., new project or past logs)
		statsData.PerPages = nil
	}

	// Calculate spec_mean_day (predicted average)
	statsData.SpecMeanDay = math.Round(float64(statsData.MeanDay)*1.15*1000) / 1000

	// Calculate max_day (maximum pages in a single day for the target weekday)
	maxDay, maxErr := h.repo.GetMaxByWeekday(ctx, targetDate)
	if maxErr == nil {
		statsData.MaxDay = maxDay
	}

	// Calculate mean_geral (overall mean across all weekdays)
	meanGeral, meanErr := h.repo.GetOverallMean(ctx, targetDate)
	if meanErr == nil {
		statsData.MeanGeral = meanGeral
	}

	// Calculate per_mean_day (ratio of current mean to previous period mean)
	prevMean, prevMeanErr := h.repo.GetPreviousPeriodMean(ctx, targetDate)
	if prevMeanErr == nil && prevMean != nil && *prevMean > 0 {
		ratio := math.Round(float64(statsData.MeanDay)/float64(*prevMean)*1000) / 1000
		statsData.PerMeanDay = &ratio
	} else {
		statsData.PerMeanDay = nil
	}

	// Calculate per_spec_mean_day (ratio of speculative mean to previous period speculative mean)
	prevSpecMean, prevSpecMeanErr := h.repo.GetPreviousPeriodSpecMean(ctx, targetDate)
	if prevSpecMeanErr == nil && prevSpecMean != nil && *prevSpecMean > 0 {
		ratio := math.Round(float64(statsData.SpecMeanDay)/float64(*prevSpecMean)*1000) / 1000
		statsData.PerSpecMeanDay = &ratio
	} else {
		statsData.PerSpecMeanDay = nil
	}

	// Return flat JSON with stats key at root level (not JSON:API envelope)
	response := map[string]interface{}{
		"stats": statsData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(response)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}

// Projects handles GET /v1/dashboard/projects.json - Returns running projects in JSON:API format
// Response structure: { "data": [...], "stats": {...} }
// Each project: { "id": "123", "type": "projects", "attributes": {...} }
// Attributes use kebab-case field names to match Rails API
func (h *DashboardHandler) Projects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Call service method to get dashboard projects in JSON:API format
	response, err := h.projectsService.GetDashboardProjects(ctx)
	if err != nil {
		slog.Error("Failed to get dashboard projects", "error", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProjectsWithLogs handles GET /v1/dashboard/projects_with_logs.json - Returns all projects with eager-loaded logs and aggregate calculations
func (h *DashboardHandler) ProjectsWithLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all projects with logs from service
	aggregates, err := h.repo.GetProjectsWithLogs(ctx)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Create a map to store logs for each project
	logsMap := make(map[int64][]*dto.LogEntry)

	// Fetch logs for all projects sequentially
	for _, agg := range aggregates {
		// Get first 4 logs for this project, ordered by date DESC
		logs, err := h.repo.GetProjectLogs(ctx, agg.ProjectID, 4)
		if err != nil {
			fmt.Printf("Handler error fetching logs for project %d: %v\n", agg.ProjectID, err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
		logsMap[agg.ProjectID] = logs
	}

	// Build result with aggregate calculations and sort by progress descending
	var results []*dto.ProjectWithLogs

	for _, agg := range aggregates {
		logs := logsMap[agg.ProjectID]

		// Calculate total_pages from all project logs (not just first 4)
		totalPages, err := h.calculateTotalPages(ctx, agg.ProjectID)
		if err != nil {
			fmt.Printf("Handler error calculating total pages for project %d: %v\n", agg.ProjectID, err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Calculate pages from all project logs (not just first 4)
		pages, err := h.calculatePages(ctx, agg.ProjectID)
		if err != nil {
			fmt.Printf("Handler error calculating pages for project %d: %v\n", agg.ProjectID, err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Calculate progress_geral
		progress := h.calculateProgress(pages, totalPages)

		result := dto.NewProjectWithLogsFromPtrs(agg, logs, totalPages, pages, progress)
		results = append(results, result)
	}

	// Sort by progress descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Progress > results[j].Progress
	})

	// Wrap response in JSON:API envelope
	envelope := dto.NewJSONAPIEnvelopeWithArray([]dto.JSONAPIData{
		{
			Type:       "dashboard_projects_with_logs",
			ID:         strconv.FormatInt(time.Now().Unix(), 10),
			Attributes: results,
		},
	})

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(envelope)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}

// calculateTotalPages calculates the total pages for a project from all logs
func (h *DashboardHandler) calculateTotalPages(ctx context.Context, projectID int64) (int, error) {
	query := `
		SELECT COALESCE(SUM(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
		WHERE project_id = $1
	`

	var totalPages int
	err := h.repo.GetPool().QueryRow(ctx, query, projectID).Scan(&totalPages)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total pages: %w", err)
	}

	return totalPages, nil
}

// calculatePages calculates the pages read for a project from all logs
func (h *DashboardHandler) calculatePages(ctx context.Context, projectID int64) (int, error) {
	query := `
		SELECT COALESCE(SUM(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
		WHERE project_id = $1
	`

	var pages int
	err := h.repo.GetPool().QueryRow(ctx, query, projectID).Scan(&pages)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate pages: %w", err)
	}

	return pages, nil
}

// calculateProgress calculates the progress percentage
func (h *DashboardHandler) calculateProgress(pages, totalPages int) float64 {
	if totalPages <= 0 {
		return 0.0
	}
	progress := float64(pages) / float64(totalPages) * 100
	return math.Round(progress*1000) / 1000
}

// LastDays handles GET /v1/dashboard/last_days.json - Returns trend data for the last N days
func (h *DashboardHandler) LastDays(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get number of days from query parameter, default to 7
	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	// Get type parameter for different trend types
	trendType := r.URL.Query().Get("type")

	// Validate type parameter (must be 1-5)
	if trendType != "" {
		typeInt, err := strconv.Atoi(trendType)
		if err != nil || typeInt < 1 || typeInt > 5 {
			http.Error(w, `{"error": "invalid type parameter", "details": {"type": "must be between 1 and 5"}}`, http.StatusUnprocessableEntity)
			return
		}
	}

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	// Get logs data from repository
	logs, err := h.repo.GetLogsByDateRange(ctx, startDate, endDate)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Calculate average per day
	var avgPerDay float64
	if days > 0 && len(logs) > 0 {
		totalPages := 0
		for _, log := range logs {
			totalPages += log.ReadPages
		}
		avgPerDay = math.Round(float64(totalPages)/float64(days)*1000) / 1000
	}

	// Build response with logs array
	response := map[string]interface{}{
		"days":         days,
		"start_date":   startDate.Format(time.RFC3339),
		"end_date":     endDate.Format(time.RFC3339),
		"total_faults": len(logs), // Use log count as faults
		"avg_per_day":  avgPerDay,
		"type":         trendType,
		"logs":         logs,
	}

	// Wrap response in JSON:API envelope
	envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
		Type:       "dashboard_last_days",
		ID:         strconv.FormatInt(time.Now().Unix(), 10),
		Attributes: response,
	})

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(envelope)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}

// =============================================================================
// ECharts Endpoints
// =============================================================================

// Faults handles GET /v1/dashboard/echart/faults.json - Returns gauge chart config for faults
func (h *DashboardHandler) Faults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create FaultsService using dependency injection pattern
	userConfigProvider := dashboard.UserConfigProvider(h.userConfig)
	faultsService := dashboard.NewFaultsService(h.repo, userConfigProvider)

	// Get fault percentage from service
	percentage, err := faultsService.GetFaultsPercentage(ctx)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Create gauge chart using service method
	gauge := faultsService.CreateGaugeChart(percentage)

	// Wrap response in JSON:API envelope
	envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
		Type:       "dashboard_echart_faults",
		ID:         strconv.FormatInt(time.Now().Unix(), 10),
		Attributes: gauge,
	})

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(envelope)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}

// SpeculateActual handles GET /v1/dashboard/echart/speculate_actual.json - Returns line chart for speculation vs actual
func (h *DashboardHandler) SpeculateActual(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get current date range (last 30 days)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// Get faults data for the period
	faults, err := h.repo.GetFaultsByDateRange(ctx, startDate, endDate)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Get prediction percentage from config
	predictionPct := h.userConfig.GetPredictionPct()
	if predictionPct == 0 {
		predictionPct = 0.15 // Default value
	}

	// Calculate predicted faults (used in daily calculation below)
	_ = int(float64(faults.FaultCount) * (1 + predictionPct))

	// Build line chart configuration
	lineChart := dto.NewEchartConfig().
		SetTitle("Speculated vs Actual Faults").
		SetTooltip(map[string]interface{}{
			"trigger": "axis",
		})

	legend := dto.NewLegend(true, []string{"Actual", "Speculated"})
	lineChart.SetLegend(legend)

	// Create axis configurations
	xAxis := dto.NewAxis("category")
	xAxis.Name = "Date"
	yAxis := dto.NewAxis("value")
	yAxis.Name = "Fault Count"

	lineChart.SetXAxis(xAxis)
	lineChart.SetYAxis(yAxis)

	// Add series for actual and predicted with 15 data points (one per day)
	actualData := make([]interface{}, 15)
	predictedData := make([]interface{}, 15)

	for i := 0; i < 15; i++ {
		dayDate := endDate.AddDate(0, 0, -i)
		dayStart := time.Date(dayDate.Year(), dayDate.Month(), dayDate.Day(), 0, 0, 0, 0, dayDate.Location())
		dayEnd := dayStart.AddDate(0, 0, 1).Add(-time.Second)

		// Get faults for this specific day
		dayFaults, _ := h.repo.GetFaultsByDateRange(ctx, dayStart, dayEnd)

		actualData[14-i] = dayFaults.FaultCount
		predictedData[14-i] = int(float64(dayFaults.FaultCount) * (1 + predictionPct))
	}

	lineChart.AddSeries(*dto.NewSeries("Actual", "line", actualData).
		SetLineStyle(map[string]interface{}{
			"width": 2,
			"type":  "solid",
		}))

	lineChart.AddSeries(*dto.NewSeries("Speculated", "line", predictedData).
		SetLineStyle(map[string]interface{}{
			"width": 2,
			"type":  "dashed",
		}))

	// Wrap response in JSON:API envelope
	envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
		Type:       "dashboard_echart_speculate_actual",
		ID:         strconv.FormatInt(time.Now().Unix(), 10),
		Attributes: lineChart,
	})

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(envelope)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}

// WeekdayFaults handles GET /v1/dashboard/echart/faults_week_day.json - Returns radar chart for weekday faults
func (h *DashboardHandler) WeekdayFaults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create WeekdayFaultsService using dependency injection pattern
	userConfigProvider := dashboard.UserConfigProvider(h.userConfig)
	weekdayFaultsService := dashboard.NewWeekdayFaultsService(h.repo, userConfigProvider)

	// Get weekday faults data (uses 6-month date range via service)
	weekdayFaults, err := weekdayFaultsService.GetWeekdayFaults(ctx)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Validate output against acceptance criteria
	if err := weekdayFaultsService.ValidateOutput(weekdayFaults.Faults); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		http.Error(w, `{"error": "validation failed"}`, http.StatusBadRequest)
		return
	}

	// Build radar chart using service helper method
	radarChart := weekdayFaultsService.CreateRadarChart(weekdayFaults.Faults)

	// Wrap response in JSON:API envelope
	envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
		Type:       "dashboard_echart_weekday_faults",
		ID:         strconv.FormatInt(time.Now().Unix(), 10),
		Attributes: radarChart,
	})

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(envelope)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}

// MeanProgress handles GET /v1/dashboard/echart/mean_progress.json - Returns line chart for mean progress
func (h *DashboardHandler) MeanProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create MeanProgressService using dependency injection pattern
	userConfigProvider := dashboard.UserConfigProvider(h.userConfig)
	meanProgressService := dashboard.NewMeanProgressService(h.repo, userConfigProvider)

	// Generate chart configuration - service always returns 30 data points
	chart, err := meanProgressService.GenerateChartConfig(ctx)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Wrap response in JSON:API envelope
	envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
		Type:       "dashboard_echart_mean_progress",
		ID:         strconv.FormatInt(time.Now().Unix(), 10),
		Attributes: chart,
	})

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(envelope)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}

// YearlyTotal handles GET /v1/dashboard/echart/last_year_total.json - Returns yearly trend chart
func (h *DashboardHandler) YearlyTotal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get current year
	now := time.Now()
	currentYear := now.Year()

	// Calculate date range for last 52 weeks
	endDate := time.Date(currentYear, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate := endDate.AddDate(0, 0, -364) // 52 weeks = 364 days

	// Get all logs for the last 52 weeks
	logs, err := h.repo.GetLogsByDateRange(ctx, startDate, endDate)
	if err != nil {
		fmt.Printf("Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Aggregate logs by week (52 weeks)
	weeklyData := make([]int, 52)
	weekLabels := make([]string, 52)

	for i := 0; i < 52; i++ {
		weekStart := startDate.AddDate(0, 0, i*7)
		weekEnd := weekStart.AddDate(0, 0, 6)
		weekLabels[i] = fmt.Sprintf("%s - %s",
			weekStart.Format("2006-W06"),
			weekEnd.Format("01/02"))

		// Count logs for this week
		count := 0
		for _, log := range logs {
			// Parse log date string
			if logTime, err := time.Parse(time.RFC3339, log.Data); err == nil {
				if !logTime.Before(weekStart) && !logTime.After(weekEnd) {
					count++
				}
			}
		}
		weeklyData[i] = count
	}

	// Build line chart configuration
	lineChart := dto.NewEchartConfig().
		SetTitle("Yearly Total Faults").
		SetTooltip(map[string]interface{}{
			"trigger": "axis",
		})

	legend := dto.NewLegend(true, []string{"Faults"})
	lineChart.SetLegend(legend)

	// Create axis configurations
	xAxis := dto.NewAxis("category")
	xAxis.Name = "Date"
	yAxis := dto.NewAxis("value")
	yAxis.Name = "Fault Count"

	lineChart.SetXAxis(xAxis)
	lineChart.SetYAxis(yAxis)

	// Convert weekly data to interface{} slice
	data := make([]interface{}, 52)
	for i, count := range weeklyData {
		data[i] = float64(count)
	}

	// Add series for weekly faults
	lineChart.AddSeries(*dto.NewSeries("Faults", "line", data).
		SetLineStyle(map[string]interface{}{
			"width": 2,
			"type":  "solid",
		}))

	// Wrap response in JSON:API envelope
	envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
		Type:       "dashboard_echart_yearly_total",
		ID:         strconv.FormatInt(time.Now().Unix(), 10),
		Attributes: lineChart,
	})

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
	// Debug: print the raw JSON
	var debugBuf []byte
	debugBuf, _ = json.Marshal(envelope)
	fmt.Printf("DEBUG: Raw JSON: %s\n", string(debugBuf))
}
