package performance

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// Dashboard Query Analyzer
// =============================================================================

// QueryTracer tracks query execution times and identifies slow queries
type QueryTracer struct {
	threshold    time.Duration
	slowQueries  []SlowQueryInfo
	queryHistory []QueryInfo
	logger       *slog.Logger
}

// SlowQueryInfo contains information about a slow query
type SlowQueryInfo struct {
	SQL        string
	Duration   time.Duration
	Threshold  time.Duration
	Timestamp  time.Time
	StackTrace string
}

// QueryInfo contains information about a query execution
type QueryInfo struct {
	SQL       string
	Duration  time.Duration
	Timestamp time.Time
	Error     error
}

// NewQueryTracer creates a new query tracer with the specified threshold
func NewQueryTracer(threshold time.Duration) *QueryTracer {
	return &QueryTracer{
		threshold:    threshold,
		slowQueries:  make([]SlowQueryInfo, 0),
		queryHistory: make([]QueryInfo, 0),
		logger:       slog.Default(),
	}
}

// NewQueryTracerWithLogger creates a new query tracer with a custom logger
func NewQueryTracerWithLogger(threshold time.Duration, logger *slog.Logger) *QueryTracer {
	return &QueryTracer{
		threshold:    threshold,
		slowQueries:  make([]SlowQueryInfo, 0),
		queryHistory: make([]QueryInfo, 0),
		logger:       logger,
	}
}

// TraceQueryStart starts tracing a query
func (t *QueryTracer) TraceQueryStart(ctx context.Context, info interface{}) context.Context {
	return context.WithValue(ctx, "query_start", time.Now())
}

// TraceQueryEnd ends tracing a query and records the result
func (t *QueryTracer) TraceQueryEnd(ctx context.Context, info interface{}, err error) {
	start := ctx.Value("query_start")
	if start == nil {
		return
	}

	duration := time.Since(start.(time.Time))
	queryInfo := QueryInfo{
		Duration:  duration,
		Timestamp: time.Now(),
		Error:     err,
	}

	t.queryHistory = append(t.queryHistory, queryInfo)

	// Check if query is slow
	if duration > t.threshold {
		slowQuery := SlowQueryInfo{
			Duration:   duration,
			Threshold:  t.threshold,
			Timestamp:  time.Now(),
			StackTrace: getStackTrace(),
		}
		t.slowQueries = append(t.slowQueries, slowQuery)

		t.logger.Warn("Slow query detected",
			slog.Duration("duration", duration),
			slog.Duration("threshold", t.threshold))
	}
}

// GetSlowQueries returns all slow queries recorded
func (t *QueryTracer) GetSlowQueries() []SlowQueryInfo {
	return t.slowQueries
}

// GetQueryHistory returns the query history
func (t *QueryTracer) GetQueryHistory() []QueryInfo {
	return t.queryHistory
}

// ClearHistory clears all recorded queries
func (t *QueryTracer) ClearHistory() {
	t.slowQueries = make([]SlowQueryInfo, 0)
	t.queryHistory = make([]QueryInfo, 0)
}

// Report generates a report of slow queries
func (t *QueryTracer) Report() string {
	var sb strings.Builder

	sb.WriteString("# Query Performance Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("**Slow Query Threshold:** %v\n\n", t.threshold))

	if len(t.slowQueries) == 0 {
		sb.WriteString("No slow queries detected.\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("## Slow Queries (%d total)\n\n", len(t.slowQueries)))

	for i, sq := range t.slowQueries {
		sb.WriteString(fmt.Sprintf("### Query %d\n\n", i+1))
		sb.WriteString(fmt.Sprintf("**Duration:** %v (threshold: %v)\n\n", sq.Duration, sq.Threshold))
		sb.WriteString(fmt.Sprintf("**Time:** %s\n\n", sq.Timestamp.Format(time.RFC3339)))
		sb.WriteString("#### Stack Trace\n\n")
		sb.WriteString("```\n")
		sb.WriteString(sq.StackTrace)
		sb.WriteString("\n```\n\n")
	}

	return sb.String()
}

// truncateSQL truncates SQL string for display
func truncateSQL(sql string, maxLen int) string {
	if len(sql) <= maxLen {
		return sql
	}
	return sql[:maxLen] + "...[truncated]"
}

// getStackTrace returns the current stack trace
func getStackTrace() string {
	buf := make([]byte, 4096)
	runtime.Stack(buf, false)
	return string(buf)
}

// =============================================================================
// Dashboard Query Analyzer Utilities
// =============================================================================

// AnalyzeDashboardQueries analyzes queries executed during a benchmark
type DashboardQueryAnalyzer struct {
	tracer     *QueryTracer
	endpoint   string
	startTime  time.Time
	endTime    time.Time
	totalCalls int
	successful int
	failed     int
}

// NewDashboardQueryAnalyzer creates a new analyzer for a specific endpoint
func NewDashboardQueryAnalyzer(endpoint string) *DashboardQueryAnalyzer {
	return &DashboardQueryAnalyzer{
		tracer:    NewQueryTracer(10 * time.Millisecond),
		endpoint:  endpoint,
		startTime: time.Now(),
	}
}

// Start begins the analysis
func (a *DashboardQueryAnalyzer) Start() {
	a.startTime = time.Now()
}

// End ends the analysis and generates a report
func (a *DashboardQueryAnalyzer) End() {
	a.endTime = time.Now()
}

// RecordCall records a query call result
func (a *DashboardQueryAnalyzer) RecordCall(success bool, duration time.Duration) {
	a.totalCalls++
	if success {
		a.successful++
	} else {
		a.failed++
	}
}

// GetReport generates a comprehensive report
func (a *DashboardQueryAnalyzer) GetReport() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Dashboard Query Analysis Report\n\n"))
	sb.WriteString(fmt.Sprintf("**Endpoint:** %s\n\n", a.endpoint))
	sb.WriteString(fmt.Sprintf("**Period:** %s to %s\n\n",
		a.startTime.Format(time.RFC3339),
		a.endTime.Format(time.RFC3339)))

	// Summary statistics
	total := a.totalCalls
	if total > 0 {
		successRate := float64(a.successful) / float64(total) * 100
		sb.WriteString("## Summary\n\n")
		sb.WriteString(fmt.Sprintf("- **Total Calls:** %d\n", total))
		sb.WriteString(fmt.Sprintf("- **Successful:** %d (%.1f%%)\n", a.successful, successRate))
		sb.WriteString(fmt.Sprintf("- **Failed:** %d (%.1f%%)\n\n", a.failed, 100-successRate))
	}

	// Slow query analysis
	sb.WriteString("## Slow Query Analysis\n\n")
	slowQueries := a.tracer.GetSlowQueries()
	if len(slowQueries) == 0 {
		sb.WriteString("No slow queries detected.\n\n")
	} else {
		sb.WriteString(fmt.Sprintf("Found %d slow queries (>10ms):\n\n", len(slowQueries)))
		for i, sq := range slowQueries {
			sb.WriteString(fmt.Sprintf("%d. **%v**\n", i+1, sq.Duration))
		}
	}

	return sb.String()
}

// =============================================================================
// Integration with Dashboard Benchmarks
// =============================================================================

// BenchmarkWithQueryAnalysis runs a benchmark with query analysis enabled
func BenchmarkWithQueryAnalysis(b *testing.B, endpoint string, fn func()) {
	analyzer := NewDashboardQueryAnalyzer(endpoint)
	analyzer.Start()

	defer func() {
		analyzer.End()
		b.Log(analyzer.GetReport())
	}()

	// Run the benchmark
	fn()

	// Record results
	analyzer.RecordCall(true, 0) // Duration will be set in the actual benchmark
}

// QueryAnalysisMiddleware wraps a handler to analyze queries
func QueryAnalysisMiddleware(next http.Handler, analyzer *DashboardQueryAnalyzer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		analyzer.RecordCall(true, duration)
	})
}
