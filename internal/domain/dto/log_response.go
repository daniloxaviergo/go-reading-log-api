package dto

import (
	"context"
)

// LogResponse represents the JSON response for Log entities
// Matches Rails LogSerializer output (without nested project)
type LogResponse struct {
	ctx       context.Context
	ID        int64            `json:"id"`
	Data      *string          `json:"data"`
	StartPage int              `json:"start_page"`
	EndPage   int              `json:"end_page"`
	Note      *string          `json:"note"`
	Project   *ProjectResponse `json:"project,omitempty"`
}

// NewLogResponse creates a new LogResponse with context
func NewLogResponse(id int64, data *string, startPage int, endPage int) *LogResponse {
	return &LogResponse{
		ctx:       context.Background(),
		ID:        id,
		Data:      data,
		StartPage: startPage,
		EndPage:   endPage,
	}
}

// GetContext returns the embedded context
func (l *LogResponse) GetContext() context.Context {
	if l.ctx == nil {
		return context.Background()
	}
	return l.ctx
}

// SetContext sets the context for the LogResponse
func (l *LogResponse) SetContext(ctx context.Context) {
	l.ctx = ctx
}
