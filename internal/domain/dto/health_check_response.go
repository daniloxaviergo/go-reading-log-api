package dto

import "context"

// HealthCheckResponse represents the JSON response for health check endpoints
type HealthCheckResponse struct {
	ctx     context.Context
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// NewHealthCheckResponse creates a new HealthCheckResponse with context
func NewHealthCheckResponse(status string) *HealthCheckResponse {
	return &HealthCheckResponse{
		ctx:    context.Background(),
		Status: status,
	}
}

// GetContext returns the embedded context
func (h *HealthCheckResponse) GetContext() context.Context {
	if h.ctx == nil {
		return context.Background()
	}
	return h.ctx
}

// SetContext sets the context for the HealthCheckResponse
func (h *HealthCheckResponse) SetContext(ctx context.Context) {
	h.ctx = ctx
}
