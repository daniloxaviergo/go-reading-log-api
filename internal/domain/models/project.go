package models

import (
	"context"
	"time"
)

// Project represents a reading project domain model
type Project struct {
	ctx        context.Context
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	TotalPage  int        `json:"total_page"`
	StartedAt  *time.Time `json:"started_at"`
	Page       int        `json:"page"`
	Reinicia   bool       `json:"reinicia"`
	Progress   *float64   `json:"progress,omitempty"`
	Status     *string    `json:"status,omitempty"`
	LogsCount  *int       `json:"logs_count,omitempty"`
	DaysUnread *int       `json:"days_unreading,omitempty"`
	MedianDay  *time.Time `json:"median_day,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

// NewProject creates a new Project with context
func NewProject(ctx context.Context, id int64, name string, totalPage int, page int, reinicia bool) *Project {
	return &Project{
		ctx:       ctx,
		ID:        id,
		Name:      name,
		TotalPage: totalPage,
		Page:      page,
		Reinicia:  reinicia,
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
