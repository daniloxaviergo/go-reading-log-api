package dto

import (
	"context"
	"time"
)

// LogResponse represents the JSON response for Log entities
// Complies with JSON:API specification
// Note: Relationships are handled separately in the JSON:API envelope,
// not embedded in attributes. This struct contains only attribute data.
type LogResponse struct {
	ctx       context.Context
	ID        int64      `json:"id"`
	Data      *time.Time `json:"data"`
	StartPage int        `json:"start_page"`
	EndPage   int        `json:"end_page"`
	Note      *string    `json:"note"`
}

// NewLogResponseWithProject creates a new LogResponse with project relationship
// Note: Relationships are now handled at the JSON:API envelope level,
// not within the LogResponse struct itself.
func NewLogResponseWithProject(id int64, data *time.Time, startPage int, endPage int, projectID int64) *LogResponse {
	return &LogResponse{
		ctx:       context.Background(),
		ID:        id,
		Data:      data,
		StartPage: startPage,
		EndPage:   endPage,
		Note:      nil,
	}
}

// NewLogResponse creates a new LogResponse with context
func NewLogResponse(id int64, data *time.Time, startPage int, endPage int) *LogResponse {
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
