package dto

import "context"

// ProjectResponse represents the JSON response for Project entities
// Matches Rails ProjectSerializer output
type ProjectResponse struct {
	ctx        context.Context
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	StartedAt  *string        `json:"started_at"`
	Progress   *float64       `json:"progress"`
	TotalPage  int            `json:"total_page"`
	Page       int            `json:"page"`
	Status     *string        `json:"status"`
	LogsCount  *int           `json:"logs_count"`
	DaysUnread *int           `json:"days_unreading"`
	MedianDay  *string        `json:"median_day"`
	FinishedAt *string        `json:"finished_at"`
	Logs       []*LogResponse `json:"logs,omitempty"`
}

// NewProjectResponse creates a new ProjectResponse with context
func NewProjectResponse(
	id int64,
	name string,
	startedAt *string,
	totalPage int,
	page int,
) *ProjectResponse {
	return &ProjectResponse{
		ctx:       context.Background(),
		ID:        id,
		Name:      name,
		StartedAt: startedAt,
		TotalPage: totalPage,
		Page:      page,
	}
}

// GetContext returns the embedded context
func (p *ProjectResponse) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the ProjectResponse
func (p *ProjectResponse) SetContext(ctx context.Context) {
	p.ctx = ctx
}
