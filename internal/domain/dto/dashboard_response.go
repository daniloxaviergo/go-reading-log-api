package dto

import (
	"context"
	"fmt"
	"math"
)

// =============================================================================
// DashboardResponse - Main response structure for all dashboard endpoints
// =============================================================================

// DashboardResponse is the main response structure for all dashboard endpoints
// Supports all 8 dashboard endpoint types with optional Echart, Stats, and Logs
type DashboardResponse struct {
	ctx    context.Context
	Echart *EchartConfig `json:"echart,omitempty"`
	Stats  *StatsData    `json:"stats,omitempty"`
	Logs   []LogEntry    `json:"logs,omitempty"`
}

// NewDashboardResponse creates a new DashboardResponse with optional fields
func NewDashboardResponse() *DashboardResponse {
	return &DashboardResponse{
		Logs: make([]LogEntry, 0),
	}
}

// GetContext returns the embedded context
func (d *DashboardResponse) GetContext() context.Context {
	if d.ctx == nil {
		return context.Background()
	}
	return d.ctx
}

// SetContext sets the context for the DashboardResponse
func (d *DashboardResponse) SetContext(ctx context.Context) {
	d.ctx = ctx
}

// SetEchart sets the ECharts configuration
func (d *DashboardResponse) SetEchart(config *EchartConfig) *DashboardResponse {
	d.Echart = config
	return d
}

// SetStats sets the statistics data
func (d *DashboardResponse) SetStats(stats *StatsData) *DashboardResponse {
	d.Stats = stats
	return d
}

// AddLog adds a log entry to the response
func (d *DashboardResponse) AddLog(log LogEntry) *DashboardResponse {
	d.Logs = append(d.Logs, log)
	return d
}

// Validate validates the DashboardResponse structure
func (d *DashboardResponse) Validate() error {
	if d == nil {
		return fmt.Errorf("dashboard response is nil")
	}

	// Validate EchartConfig if present
	if d.Echart != nil {
		if err := d.Echart.Validate(); err != nil {
			return fmt.Errorf("echart config validation failed: %w", err)
		}
	}

	// Validate StatsData if present
	if d.Stats != nil {
		if err := d.Stats.Validate(); err != nil {
			return fmt.Errorf("stats data validation failed: %w", err)
		}
	}

	// Validate each log entry
	for i, log := range d.Logs {
		if err := log.Validate(); err != nil {
			return fmt.Errorf("log entry %d validation failed: %w", i, err)
		}
	}

	return nil
}

// =============================================================================
// EchartConfig - ECharts-style chart configurations
// =============================================================================

// EchartConfig holds ECharts-style chart configurations
type EchartConfig struct {
	ctx     context.Context
	Title   string                 `json:"title,omitempty"`
	Tooltip map[string]interface{} `json:"tooltip,omitempty"`
	Legend  *Legend                `json:"legend,omitempty"`
	Series  []Series               `json:"series,omitempty"`
	XAxis   *Axis                  `json:"xAxis,omitempty"`
	YAxis   *Axis                  `json:"yAxis,omitempty"`
	Grid    *Grid                  `json:"grid,omitempty"`
	Toolbox map[string]interface{} `json:"toolbox,omitempty"`
}

// NewEchartConfig creates a new EchartConfig with empty series
func NewEchartConfig() *EchartConfig {
	return &EchartConfig{
		Series: make([]Series, 0),
	}
}

// NewEchartConfigWithSeries creates a new EchartConfig with initial series
func NewEchartConfigWithSeries(series []Series) *EchartConfig {
	return &EchartConfig{
		Series: series,
	}
}

// GetContext returns the embedded context
func (e *EchartConfig) GetContext() context.Context {
	if e.ctx == nil {
		return context.Background()
	}
	return e.ctx
}

// SetContext sets the context for the EchartConfig
func (e *EchartConfig) SetContext(ctx context.Context) {
	e.ctx = ctx
}

// SetTitle sets the chart title
func (e *EchartConfig) SetTitle(title string) *EchartConfig {
	e.Title = title
	return e
}

// SetTooltip sets the tooltip configuration
func (e *EchartConfig) SetTooltip(tooltip map[string]interface{}) *EchartConfig {
	e.Tooltip = tooltip
	return e
}

// SetLegend sets the legend configuration
func (e *EchartConfig) SetLegend(legend *Legend) *EchartConfig {
	e.Legend = legend
	return e
}

// AddSeries adds a series to the chart
func (e *EchartConfig) AddSeries(series Series) *EchartConfig {
	e.Series = append(e.Series, series)
	return e
}

// SetXAxis sets the X axis configuration
func (e *EchartConfig) SetXAxis(axis *Axis) *EchartConfig {
	e.XAxis = axis
	return e
}

// SetYAxis sets the Y axis configuration
func (e *EchartConfig) SetYAxis(axis *Axis) *EchartConfig {
	e.YAxis = axis
	return e
}

// SetGrid sets the grid configuration
func (e *EchartConfig) SetGrid(grid *Grid) *EchartConfig {
	e.Grid = grid
	return e
}

// SetToolbox sets the toolbox configuration
func (e *EchartConfig) SetToolbox(toolbox map[string]interface{}) *EchartConfig {
	e.Toolbox = toolbox
	return e
}

// Validate validates the EchartConfig structure
func (e *EchartConfig) Validate() error {
	if e == nil {
		return fmt.Errorf("echart config is nil")
	}

	// Validate series
	if len(e.Series) == 0 {
		return fmt.Errorf("series cannot be empty")
	}

	for i, series := range e.Series {
		if err := series.Validate(); err != nil {
			return fmt.Errorf("series %d validation failed: %w", i, err)
		}
	}

	return nil
}

// =============================================================================
// Legend - Legend configuration for charts
// =============================================================================

// Legend configuration
type Legend struct {
	ctx  context.Context
	Show bool     `json:"show"`
	Data []string `json:"data,omitempty"`
	Top  string   `json:"top,omitempty"`
	Left string   `json:"left,omitempty"`
}

// NewLegend creates a new Legend configuration
func NewLegend(show bool, data []string) *Legend {
	return &Legend{
		Show: show,
		Data: data,
	}
}

// GetContext returns the embedded context
func (l *Legend) GetContext() context.Context {
	if l.ctx == nil {
		return context.Background()
	}
	return l.ctx
}

// SetContext sets the context for the Legend
func (l *Legend) SetContext(ctx context.Context) {
	l.ctx = ctx
}

// Validate validates the Legend structure
func (l *Legend) Validate() error {
	if l == nil {
		return fmt.Errorf("legend is nil")
	}
	return nil
}

// =============================================================================
// Series - Chart series definition
// =============================================================================

// Series represents a chart series
type Series struct {
	ctx       context.Context
	Name      string                 `json:"name"`
	Type      string                 `json:"type"` // line, bar, pie, gauge, radar
	Data      []interface{}          `json:"data"`
	ItemStyle map[string]interface{} `json:"itemStyle,omitempty"`
	LineStyle map[string]interface{} `json:"lineStyle,omitempty"`
	AreaStyle map[string]interface{} `json:"areaStyle,omitempty"`
}

// NewSeries creates a new Series with required fields
func NewSeries(name, chartType string, data []interface{}) *Series {
	return &Series{
		Name:      name,
		Type:      chartType,
		Data:      data,
		ItemStyle: make(map[string]interface{}),
	}
}

// GetContext returns the embedded context
func (s *Series) GetContext() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}

// SetContext sets the context for the Series
func (s *Series) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// SetItemStyle sets the item style configuration
func (s *Series) SetItemStyle(style map[string]interface{}) *Series {
	s.ItemStyle = style
	return s
}

// SetLineStyle sets the line style configuration
func (s *Series) SetLineStyle(style map[string]interface{}) *Series {
	s.LineStyle = style
	return s
}

// SetAreaStyle sets the area style configuration
func (s *Series) SetAreaStyle(style map[string]interface{}) *Series {
	s.AreaStyle = style
	return s
}

// Validate validates the Series structure
func (s *Series) Validate() error {
	if s == nil {
		return fmt.Errorf("series is nil")
	}

	if s.Name == "" {
		return fmt.Errorf("series name is required")
	}
	if s.Type == "" {
		return fmt.Errorf("series type is required")
	}
	if len(s.Data) == 0 {
		return fmt.Errorf("series data cannot be empty")
	}

	return nil
}

// =============================================================================
// Axis - Axis configuration for charts
// =============================================================================

// Axis configuration
type Axis struct {
	ctx         context.Context
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	BoundaryGap []bool `json:"boundaryGap,omitempty"`
}

// NewAxis creates a new Axis configuration
func NewAxis(axisType string) *Axis {
	return &Axis{
		Type: axisType,
	}
}

// GetContext returns the embedded context
func (a *Axis) GetContext() context.Context {
	if a.ctx == nil {
		return context.Background()
	}
	return a.ctx
}

// SetContext sets the context for the Axis
func (a *Axis) SetContext(ctx context.Context) {
	a.ctx = ctx
}

// SetXAxisName sets the X axis name (convenience method)
func (a *Axis) SetXAxisName(name string) *Axis {
	a.Name = name
	return a
}

// SetYAxisName sets the Y axis name (convenience method)
func (a *Axis) SetYAxisName(name string) *Axis {
	a.Name = name
	return a
}

// Validate validates the Axis structure
func (a *Axis) Validate() error {
	if a == nil {
		return fmt.Errorf("axis is nil")
	}
	if a.Type == "" {
		return fmt.Errorf("axis type is required")
	}
	return nil
}

// =============================================================================
// Grid - Grid configuration for charts
// =============================================================================

// Grid configuration
type Grid struct {
	ctx    context.Context
	Left   string `json:"left,omitempty"`
	Right  string `json:"right,omitempty"`
	Top    string `json:"top,omitempty"`
	Bottom string `json:"bottom,omitempty"`
}

// NewGrid creates a new Grid configuration
func NewGrid() *Grid {
	return &Grid{}
}

// GetContext returns the embedded context
func (g *Grid) GetContext() context.Context {
	if g.ctx == nil {
		return context.Background()
	}
	return g.ctx
}

// SetContext sets the context for the Grid
func (g *Grid) SetContext(ctx context.Context) {
	g.ctx = ctx
}

// Validate validates the Grid structure
func (g *Grid) Validate() error {
	if g == nil {
		return fmt.Errorf("grid is nil")
	}
	return nil
}

// =============================================================================
// StatsData - Aggregate statistics for dashboard views
// =============================================================================

// StatsData holds all statistical calculations for dashboard views
type StatsData struct {
	ctx               context.Context
	PreviousWeekPages int      `json:"previous_week_pages"`         // Sum of pages from previous 7 days
	LastWeekPages     int      `json:"last_week_pages"`             // Sum of pages from last 7 days
	PerPages          *float64 `json:"per_pages,omitempty"`         // Ratio: last_week / previous_week (nullable)
	MeanDay           float64  `json:"mean_day"`                    // Average pages per day (current weekday)
	SpecMeanDay       float64  `json:"spec_mean_day"`               // Predicted average for current weekday
	ProgressGeral     float64  `json:"progress_geral"`              // Overall completion percentage
	TotalPages        int      `json:"total_pages"`                 // Sum of all project total_page values
	Pages             int      `json:"pages"`                       // Sum of all project page values
	CountPages        int      `json:"count_pages"`                 // Sum of read_pages in period
	SpeculatePages    int      `json:"speculate_pages"`             // Config-based target for period
	MaxDay            *float64 `json:"max_day,omitempty"`           // Maximum pages in a single day (nullable)
	MeanGeral         *float64 `json:"mean_geral,omitempty"`        // General mean across all days (nullable)
	PerMeanDay        *float64 `json:"per_mean_day,omitempty"`      // Ratio for mean day (nullable)
	PerSpecMeanDay    *float64 `json:"per_spec_mean_day,omitempty"` // Ratio for speculative mean day (nullable)
}

// NewStatsData creates a new StatsData with zero values
func NewStatsData() *StatsData {
	return &StatsData{
		PerPages:       nil,
		MeanDay:        0.0,
		SpecMeanDay:    0.0,
		ProgressGeral:  0.0,
		MaxDay:         nil,
		MeanGeral:      nil,
		PerMeanDay:     nil,
		PerSpecMeanDay: nil,
	}
}

// GetContext returns the embedded context
func (s *StatsData) GetContext() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}

// SetContext sets the context for the StatsData
func (s *StatsData) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// SetPreviousWeekPages sets the previous week pages
func (s *StatsData) SetPreviousWeekPages(pages int) *StatsData {
	s.PreviousWeekPages = pages
	return s
}

// SetLastWeekPages sets the last week pages
func (s *StatsData) SetLastWeekPages(pages int) *StatsData {
	s.LastWeekPages = pages
	return s
}

// SetPerPages sets the per pages ratio
func (s *StatsData) SetPerPages(ratio *float64) *StatsData {
	s.PerPages = ratio
	return s
}

// SetMaxDay sets the maximum day pages
func (s *StatsData) SetMaxDay(pages *float64) *StatsData {
	s.MaxDay = pages
	return s
}

// SetMeanGeral sets the general mean
func (s *StatsData) SetMeanGeral(mean *float64) *StatsData {
	s.MeanGeral = mean
	return s
}

// SetPerMeanDay sets the per mean day ratio
func (s *StatsData) SetPerMeanDay(ratio *float64) *StatsData {
	s.PerMeanDay = ratio
	return s
}

// SetPerSpecMeanDay sets the per speculative mean day ratio
func (s *StatsData) SetPerSpecMeanDay(ratio *float64) *StatsData {
	s.PerSpecMeanDay = ratio
	return s
}

// SetMeanDay sets the mean day pages
func (s *StatsData) SetMeanDay(pages float64) *StatsData {
	s.MeanDay = pages
	return s
}

// SetSpecMeanDay sets the speculative mean day pages
func (s *StatsData) SetSpecMeanDay(pages float64) *StatsData {
	s.SpecMeanDay = pages
	return s
}

// SetProgressGeral sets the overall progress percentage
func (s *StatsData) SetProgressGeral(percentage float64) *StatsData {
	s.ProgressGeral = percentage
	return s
}

// SetTotalPages sets the total pages
func (s *StatsData) SetTotalPages(pages int) *StatsData {
	s.TotalPages = pages
	return s
}

// SetPages sets the current pages
func (s *StatsData) SetPages(pages int) *StatsData {
	s.Pages = pages
	return s
}

// SetCountPages sets the count pages
func (s *StatsData) SetCountPages(pages int) *StatsData {
	s.CountPages = pages
	return s
}

// SetSpeculatePages sets the speculative pages
func (s *StatsData) SetSpeculatePages(pages int) *StatsData {
	s.SpeculatePages = pages
	return s
}

// RoundToThreeDecimals rounds float64 values to 3 decimal places
func (s *StatsData) RoundToThreeDecimals() {
	if s.PerPages != nil {
		rounded := math.Round(*s.PerPages*1000) / 1000
		s.PerPages = &rounded
	}
	s.MeanDay = math.Round(s.MeanDay*1000) / 1000
	s.SpecMeanDay = math.Round(s.SpecMeanDay*1000) / 1000
	s.ProgressGeral = math.Round(s.ProgressGeral*1000) / 1000
	if s.MaxDay != nil {
		rounded := math.Round(*s.MaxDay*1000) / 1000
		s.MaxDay = &rounded
	}
	if s.MeanGeral != nil {
		rounded := math.Round(*s.MeanGeral*1000) / 1000
		s.MeanGeral = &rounded
	}
	if s.PerMeanDay != nil {
		rounded := math.Round(*s.PerMeanDay*1000) / 1000
		s.PerMeanDay = &rounded
	}
	if s.PerSpecMeanDay != nil {
		rounded := math.Round(*s.PerSpecMeanDay*1000) / 1000
		s.PerSpecMeanDay = &rounded
	}
}

// Validate validates the StatsData structure
func (s *StatsData) Validate() error {
	if s == nil {
		return fmt.Errorf("stats data is nil")
	}

	// Validate non-negative values
	if s.PreviousWeekPages < 0 {
		return fmt.Errorf("previous_week_pages cannot be negative")
	}
	if s.LastWeekPages < 0 {
		return fmt.Errorf("last_week_pages cannot be negative")
	}
	if s.TotalPages < 0 {
		return fmt.Errorf("total_pages cannot be negative")
	}
	if s.Pages < 0 {
		return fmt.Errorf("pages cannot be negative")
	}

	// Validate ProgressGeral percentage range (0-100)
	if s.ProgressGeral < 0 || s.ProgressGeral > 100 {
		return fmt.Errorf("progress_geral must be between 0 and 100")
	}

	// Validate non-negative floats (non-pointer fields)
	if s.MeanDay < 0 {
		return fmt.Errorf("mean_day cannot be negative")
	}
	if s.SpecMeanDay < 0 {
		return fmt.Errorf("spec_mean_day cannot be negative")
	}

	// Validate pointer fields when non-nil
	// PerPages: Remove 0-100 constraint, ratios can exceed 100%
	if s.PerPages != nil {
		if *s.PerPages < 0 {
			return fmt.Errorf("per_pages cannot be negative")
		}
	}

	// MaxDay: Must be non-negative when set
	if s.MaxDay != nil {
		if *s.MaxDay < 0 {
			return fmt.Errorf("max_day cannot be negative")
		}
	}

	// MeanGeral: Must be non-negative when set
	if s.MeanGeral != nil {
		if *s.MeanGeral < 0 {
			return fmt.Errorf("mean_geral cannot be negative")
		}
	}

	// PerMeanDay: Must be non-negative when set
	if s.PerMeanDay != nil {
		if *s.PerMeanDay < 0 {
			return fmt.Errorf("per_mean_day cannot be negative")
		}
	}

	// PerSpecMeanDay: Must be non-negative when set
	if s.PerSpecMeanDay != nil {
		if *s.PerSpecMeanDay < 0 {
			return fmt.Errorf("per_spec_mean_day cannot be negative")
		}
	}

	return nil
}

// =============================================================================
// ProjectAggregateResponse - Response DTO for project aggregates with progress
// =============================================================================

// ProjectAggregateResponse represents a project aggregate with calculated progress
type ProjectAggregateResponse struct {
	ctx         context.Context
	ProjectID   int64   `json:"project_id"`
	ProjectName string  `json:"project_name"`
	TotalPages  int     `json:"total_pages"`
	LogCount    int     `json:"log_count"`
	Progress    float64 `json:"progress"`
}

// NewProjectAggregateResponse creates a new ProjectAggregateResponse
func NewProjectAggregateResponse(projectID int64, projectName string, totalPages, logCount int, progress float64) *ProjectAggregateResponse {
	return &ProjectAggregateResponse{
		ProjectID:   projectID,
		ProjectName: projectName,
		TotalPages:  totalPages,
		LogCount:    logCount,
		Progress:    math.Round(progress*1000) / 1000,
	}
}

// GetContext returns the embedded context
func (p *ProjectAggregateResponse) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the ProjectAggregateResponse
func (p *ProjectAggregateResponse) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// Validate validates the ProjectAggregateResponse structure
func (p *ProjectAggregateResponse) Validate() error {
	if p == nil {
		return fmt.Errorf("project aggregate response is nil")
	}
	if p.ProjectID <= 0 {
		return fmt.Errorf("project_id must be positive")
	}
	if p.TotalPages < 0 {
		return fmt.Errorf("total_pages cannot be negative")
	}
	if p.LogCount < 0 {
		return fmt.Errorf("log_count cannot be negative")
	}
	if p.Progress < 0 || p.Progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}
	return nil
}

// =============================================================================
// LogEntry - Eager-loaded log entry with project data
// =============================================================================

// LogEntry represents a log entry with eager-loaded project data
type LogEntry struct {
	ctx       context.Context
	ID        int64    `json:"id"`
	ProjectID int64    `json:"project_id"`
	Data      string   `json:"data"` // RFC3339 formatted timestamp
	StartPage int      `json:"start_page"`
	EndPage   int      `json:"end_page"`
	Note      *string  `json:"note,omitempty"`
	Project   *Project `json:"project"`    // Eager-loaded project data
	ReadPages int      `json:"read_pages"` // Calculated: end_page - start_page
}

// NewLogEntry creates a new LogEntry with calculated read_pages
func NewLogEntry(id int64, data string, startPage, endPage int, note *string, project *Project) *LogEntry {
	return &LogEntry{
		ID:        id,
		Data:      data,
		StartPage: startPage,
		EndPage:   endPage,
		Note:      note,
		Project:   project,
		ReadPages: endPage - startPage,
	}
}

// GetContext returns the embedded context
func (l *LogEntry) GetContext() context.Context {
	if l.ctx == nil {
		return context.Background()
	}
	return l.ctx
}

// SetContext sets the context for the LogEntry
func (l *LogEntry) SetContext(ctx context.Context) {
	l.ctx = ctx
}

// Validate validates the LogEntry structure
func (l *LogEntry) Validate() error {
	if l == nil {
		return fmt.Errorf("log entry is nil")
	}

	// Validate required fields
	if l.ID <= 0 {
		return fmt.Errorf("id must be positive")
	}
	if l.StartPage < 0 {
		return fmt.Errorf("start_page cannot be negative")
	}
	if l.EndPage < 0 {
		return fmt.Errorf("end_page cannot be negative")
	}
	if l.ReadPages < 0 {
		return fmt.Errorf("read_pages cannot be negative (end_page must >= start_page)")
	}

	// Validate project data if present
	if l.Project != nil {
		if err := l.Project.Validate(); err != nil {
			return fmt.Errorf("project validation failed: %w", err)
		}
	}

	return nil
}

// =============================================================================
// Project - Minimal project data for eager loading in LogEntry
// =============================================================================

// Project represents minimal project data for eager loading
// Note: This is a simplified version without circular references
type Project struct {
	ctx       context.Context
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	TotalPage int    `json:"total_page"`
	Page      int    `json:"page"`
}

// NewProject creates a new Project instance
func NewProject(id int64, name string, totalPage, page int) *Project {
	return &Project{
		ctx:       context.Background(),
		ID:        id,
		Name:      name,
		TotalPage: totalPage,
		Page:      page,
	}
}

// GetContext returns the embedded context
func (p *Project) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the Project
func (p *Project) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// Validate validates the Project structure
func (p *Project) Validate() error {
	if p == nil {
		return fmt.Errorf("project is nil")
	}
	if p.ID <= 0 {
		return fmt.Errorf("project id must be positive")
	}
	if p.TotalPage <= 0 {
		return fmt.Errorf("total_page must be positive")
	}
	if p.Page < 0 {
		return fmt.Errorf("page cannot be negative")
	}
	if p.Page > p.TotalPage {
		return fmt.Errorf("page cannot exceed total_page")
	}
	return nil
}

// =============================================================================
// Legacy Types - Included for compatibility with existing repository interface
// =============================================================================

// DailyStats represents daily page statistics with weekday breakdown
type DailyStats struct {
	ctx        context.Context
	TotalPages int `json:"total_pages"`
	LogCount   int `json:"log_count"`
}

// NewDailyStats creates a new DailyStats response
func NewDailyStats(totalPages, logCount int) *DailyStats {
	return &DailyStats{
		TotalPages: totalPages,
		LogCount:   logCount,
	}
}

// GetContext returns the embedded context
func (d *DailyStats) GetContext() context.Context {
	if d.ctx == nil {
		return context.Background()
	}
	return d.ctx
}

// SetContext sets the context for the DailyStats
func (d *DailyStats) SetContext(ctx context.Context) {
	d.ctx = ctx
}

// ProjectAggregate represents project-level aggregation data
type ProjectAggregate struct {
	ctx         context.Context
	ProjectID   int64  `json:"project_id"`
	ProjectName string `json:"project_name"`
	TotalPages  int    `json:"total_pages"` // Sum of read pages from logs
	LogCount    int    `json:"log_count"`
	TotalPage   int    `json:"total_page"` // Project's total_page field (for progress calculation)
}

// NewProjectAggregate creates a new ProjectAggregate response
func NewProjectAggregate(projectID int64, projectName string, totalPages, logCount int) *ProjectAggregate {
	return &ProjectAggregate{
		ctx:         context.Background(),
		ProjectID:   projectID,
		ProjectName: projectName,
		TotalPages:  totalPages,
		LogCount:    logCount,
	}
}

// GetContext returns the embedded context
func (p *ProjectAggregate) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the ProjectAggregate
func (p *ProjectAggregate) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// FaultStats represents fault counting statistics
type FaultStats struct {
	ctx        context.Context
	FaultCount int `json:"fault_count"`
}

// NewFaultStats creates a new FaultStats response
func NewFaultStats(faultCount int) *FaultStats {
	return &FaultStats{
		FaultCount: faultCount,
	}
}

// GetContext returns the embedded context
func (f *FaultStats) GetContext() context.Context {
	if f.ctx == nil {
		return context.Background()
	}
	return f.ctx
}

// SetContext sets the context for the FaultStats
func (f *FaultStats) SetContext(ctx context.Context) {
	f.ctx = ctx
}

// WeekdayFaults represents fault distribution by weekday
type WeekdayFaults struct {
	ctx    context.Context
	Faults map[int]int `json:"faults"`
}

// NewWeekdayFaults creates a new WeekdayFaults response
func NewWeekdayFaults(faults map[int]int) *WeekdayFaults {
	return &WeekdayFaults{
		Faults: faults,
	}
}

// GetContext returns the embedded context
func (w *WeekdayFaults) GetContext() context.Context {
	if w.ctx == nil {
		return context.Background()
	}
	return w.ctx
}

// SetContext sets the context for the WeekdayFaults
func (w *WeekdayFaults) SetContext(ctx context.Context) {
	w.ctx = ctx
}

// =============================================================================
// ProjectWithLogs - Project with eager-loaded logs and aggregate calculations
// =============================================================================

// ProjectWithLogs represents a project with its eager-loaded logs and aggregate calculations
type ProjectWithLogs struct {
	ctx        context.Context
	Project    *ProjectAggregateResponse `json:"project"`
	Logs       []LogEntry                `json:"logs"`
	TotalPages int                       `json:"total_pages"`
	Pages      int                       `json:"pages"`
	Progress   float64                   `json:"progress_geral"`
}

// NewProjectWithLogs creates a new ProjectWithLogs instance
func NewProjectWithLogs(project *ProjectAggregateResponse, logs []LogEntry, totalPages, pages int, progress float64) *ProjectWithLogs {
	return &ProjectWithLogs{
		Project:    project,
		Logs:       logs,
		TotalPages: totalPages,
		Pages:      pages,
		Progress:   math.Round(progress*1000) / 1000,
	}
}

// NewProjectWithLogsFromPtrs creates a new ProjectWithLogs instance from pointer slice
func NewProjectWithLogsFromPtrs(project *ProjectAggregateResponse, logs []*LogEntry, totalPages, pages int, progress float64) *ProjectWithLogs {
	// Convert []*LogEntry to []LogEntry
	logSlice := make([]LogEntry, len(logs))
	for i, log := range logs {
		logSlice[i] = *log
	}
	return &ProjectWithLogs{
		Project:    project,
		Logs:       logSlice,
		TotalPages: totalPages,
		Pages:      pages,
		Progress:   math.Round(progress*1000) / 1000,
	}
}

// GetContext returns the embedded context
func (p *ProjectWithLogs) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the ProjectWithLogs
func (p *ProjectWithLogs) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// Validate validates the ProjectWithLogs structure
func (p *ProjectWithLogs) Validate() error {
	if p == nil {
		return fmt.Errorf("project with logs is nil")
	}

	// Validate project data if present
	if p.Project != nil {
		if err := p.Project.Validate(); err != nil {
			return fmt.Errorf("project validation failed: %w", err)
		}
	}

	// Validate each log entry
	for i, log := range p.Logs {
		if err := log.Validate(); err != nil {
			return fmt.Errorf("log entry %d validation failed: %w", i, err)
		}
	}

	// Validate non-negative values
	if p.TotalPages < 0 {
		return fmt.Errorf("total_pages cannot be negative")
	}
	if p.Pages < 0 {
		return fmt.Errorf("pages cannot be negative")
	}
	if p.Progress < 0 || p.Progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	return nil
}

// =============================================================================
// DashboardProjectsResponse - JSON:API format for /v1/dashboard/projects.json
// =============================================================================

// DashboardProjectsResponse is the root response structure for the dashboard projects endpoint
// Matches Rails JSON:API format: { "data": [...], "stats": {...} }
type DashboardProjectsResponse struct {
	ctx   context.Context
	Data  []DashboardProjectItem `json:"data"`
	Stats *DashboardStats        `json:"stats"`
}

// NewDashboardProjectsResponse creates a new DashboardProjectsResponse
func NewDashboardProjectsResponse() *DashboardProjectsResponse {
	return &DashboardProjectsResponse{
		Data:  make([]DashboardProjectItem, 0),
		Stats: NewDashboardStats(),
	}
}

// GetContext returns the embedded context
func (d *DashboardProjectsResponse) GetContext() context.Context {
	if d.ctx == nil {
		return context.Background()
	}
	return d.ctx
}

// SetContext sets the context for the DashboardProjectsResponse
func (d *DashboardProjectsResponse) SetContext(ctx context.Context) {
	d.ctx = ctx
}

// AddProject adds a project item to the data array
func (d *DashboardProjectsResponse) AddProject(project DashboardProjectItem) *DashboardProjectsResponse {
	d.Data = append(d.Data, project)
	return d
}

// SetStats sets the stats object
func (d *DashboardProjectsResponse) SetStats(stats *DashboardStats) *DashboardProjectsResponse {
	d.Stats = stats
	return d
}

// DashboardProjectItem represents a single project in JSON:API format
// Structure: { "id": "123", "type": "projects", "attributes": {...} }
type DashboardProjectItem struct {
	ctx        context.Context
	ID         string                      `json:"id"`         // String ID (e.g., "446")
	Type       string                      `json:"type"`       // Always "projects"
	Attributes *DashboardProjectAttributes `json:"attributes"` // Flattened attributes
}

// NewDashboardProjectItem creates a new DashboardProjectItem
func NewDashboardProjectItem(id string, attributes *DashboardProjectAttributes) *DashboardProjectItem {
	return &DashboardProjectItem{
		ID:         id,
		Type:       "projects",
		Attributes: attributes,
	}
}

// GetContext returns the embedded context
func (p *DashboardProjectItem) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the DashboardProjectItem
func (p *DashboardProjectItem) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// DashboardProjectAttributes contains the flattened project attributes
// Uses kebab-case field names to match Rails API
type DashboardProjectAttributes struct {
	ctx           context.Context
	Name          string  `json:"name"`                 // project_name
	StartedAt     *string `json:"started-at,omitempty"` // Calculated from earliest log
	Progress      float64 `json:"progress"`             // (pages / total_page) * 100
	TotalPage     int     `json:"total-page"`           // total_pages
	Page          int     `json:"page"`                 // pages
	Status        string  `json:"status"`               // "stopped" or calculated
	LogsCount     int     `json:"logs-count"`           // log_count
	DaysUnreading int     `json:"days-unreading"`       // Calculated from latest log
}

// NewDashboardProjectAttributes creates a new DashboardProjectAttributes
func NewDashboardProjectAttributes() *DashboardProjectAttributes {
	return &DashboardProjectAttributes{
		Status: "stopped", // Default status
	}
}

// GetContext returns the embedded context
func (a *DashboardProjectAttributes) GetContext() context.Context {
	if a.ctx == nil {
		return context.Background()
	}
	return a.ctx
}

// SetContext sets the context for the DashboardProjectAttributes
func (a *DashboardProjectAttributes) SetContext(ctx context.Context) {
	a.ctx = ctx
}

// SetStartedAt sets the started-at field with string pointer
func (a *DashboardProjectAttributes) SetStartedAt(date string) *DashboardProjectAttributes {
	a.StartedAt = &date
	return a
}

// DashboardStats contains the simplified stats object
// Only keeps: progress_geral, total_pages, pages (matches Rails format)
type DashboardStats struct {
	ctx           context.Context
	ProgressGeral float64 `json:"progress_geral"` // Overall completion percentage
	TotalPages    int     `json:"total_pages"`    // Sum of all project total_page values
	Pages         int     `json:"pages"`          // Sum of all project page values
}

// NewDashboardStats creates a new DashboardStats with zero values
func NewDashboardStats() *DashboardStats {
	return &DashboardStats{
		ProgressGeral: 0.0,
		TotalPages:    0,
		Pages:         0,
	}
}

// GetContext returns the embedded context
func (s *DashboardStats) GetContext() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}

// SetContext sets the context for the DashboardStats
func (s *DashboardStats) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// SetProgressGeral sets the progress_geral field
func (s *DashboardStats) SetProgressGeral(percentage float64) *DashboardStats {
	s.ProgressGeral = percentage
	return s
}

// SetTotalPages sets the total_pages field
func (s *DashboardStats) SetTotalPages(pages int) *DashboardStats {
	s.TotalPages = pages
	return s
}

// SetPages sets the pages field
func (s *DashboardStats) SetPages(pages int) *DashboardStats {
	s.Pages = pages
	return s
}

// =============================================================================
// ProgressDay - Daily progress entry with color for visual map
// =============================================================================

// ProgressDay represents a single day's progress with calculated value and color.
// Used for mean progress visualization with ECharts visual map.
type ProgressDay struct {
	ctx        context.Context
	Date       string  `json:"date"`        // RFC3339 formatted date
	DailyPages float64 `json:"daily_pages"` // Pages read on this day
	Progress   float64 `json:"progress"`    // Calculated progress percentage
	Color      string  `json:"color"`       // Color for visual map
}

// NewProgressDay creates a new ProgressDay instance.
func NewProgressDay(date string, dailyPages, progress float64, color string) *ProgressDay {
	return &ProgressDay{
		Date:       date,
		DailyPages: math.Round(dailyPages*1000) / 1000,
		Progress:   math.Round(progress*1000) / 1000,
		Color:      color,
	}
}

// GetContext returns the embedded context.
func (p *ProgressDay) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the ProgressDay.
func (p *ProgressDay) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// Validate validates the ProgressDay structure.
func (p *ProgressDay) Validate() error {
	if p == nil {
		return fmt.Errorf("progress day is nil")
	}
	if p.Date == "" {
		return fmt.Errorf("date is required")
	}
	if p.DailyPages < 0 {
		return fmt.Errorf("daily_pages cannot be negative")
	}
	if p.Progress < -100 || p.Progress > 100 {
		return fmt.Errorf("progress must be between -100 and 100")
	}
	if p.Color == "" {
		return fmt.Errorf("color is required")
	}
	return nil
}

// ValidateSlice validates a slice of ProgressDay entries.
func ValidateProgressDays(days []*ProgressDay) error {
	if days == nil {
		return fmt.Errorf("progress days slice is nil")
	}
	for i, day := range days {
		if err := day.Validate(); err != nil {
			return fmt.Errorf("progress day %d validation failed: %w", i, err)
		}
	}
	return nil
}
