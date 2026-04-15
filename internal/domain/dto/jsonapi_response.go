package dto

import "context"

// JSONAPIEnvelope represents the JSON:API response envelope
// Format: { "data": { "type": "...", "attributes": {...} } }
type JSONAPIEnvelope struct {
	Data JSONAPIData `json:"data"`
}

// JSONAPIData represents the data object within JSON:API envelope
type JSONAPIData struct {
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
	ID         interface{} `json:"id,omitempty"`
}

// NewJSONAPIEnvelope creates a new JSON:API envelope with the given data
func NewJSONAPIEnvelope(data JSONAPIData) *JSONAPIEnvelope {
	return &JSONAPIEnvelope{Data: data}
}

// ProjectJSONAPIResponse represents a JSON:API formatted project response
// Matches Rails JSON:API serialization format
type ProjectJSONAPIResponse struct {
	ctx        context.Context
	ID         interface{}      `json:"id"`
	Type       string           `json:"type"`
	Attributes *ProjectResponse `json:"attributes"`
}

// NewProjectJSONAPIResponse creates a new JSON:API formatted project response
func NewProjectJSONAPIResponse(project *ProjectResponse) *ProjectJSONAPIResponse {
	return &ProjectJSONAPIResponse{
		ctx:        context.Background(),
		ID:         project.ID,
		Type:       "projects",
		Attributes: project,
	}
}

// GetContext returns the embedded context
func (p *ProjectJSONAPIResponse) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the ProjectJSONAPIResponse
func (p *ProjectJSONAPIResponse) SetContext(ctx context.Context) {
	p.ctx = ctx
}
