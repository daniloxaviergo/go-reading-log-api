package models

import (
	"context"
	"time"
)

// Log represents a reading log entry domain model
type Log struct {
	ctx       context.Context
	ID        int64      `json:"id"`
	ProjectID int64      `json:"project_id"`
	Data      *string    `json:"data"`
	StartPage int        `json:"start_page"`
	EndPage   int        `json:"end_page"`
	Wday      int        `json:"wday"`
	Note      *string    `json:"note"`
	Text      *string    `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// NewLog creates a new Log with context
func NewLog(ctx context.Context, id int64, projectID int64, startPage int, endPage int) *Log {
	return &Log{
		ctx:       ctx,
		ID:        id,
		ProjectID: projectID,
		StartPage: startPage,
		EndPage:   endPage,
	}
}

// GetContext returns the embedded context
func (l *Log) GetContext() context.Context {
	if l.ctx == nil {
		return context.Background()
	}
	return l.ctx
}

// SetContext sets the context for the Log
func (l *Log) SetContext(ctx context.Context) {
	l.ctx = ctx
}
