package dto

import (
	"context"
	"strconv"
)

// JSONAPIEnvelope represents the JSON:API response envelope
// Format: { "data": { "type": "...", "attributes": {...} } }
// The data field can contain either a single JSONAPIData object or an array of them
type JSONAPIEnvelope struct {
	Data interface{} `json:"data"`
}

// JSONAPIData represents the data object within JSON:API envelope
type JSONAPIData struct {
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
	ID         interface{} `json:"id,omitempty"`
}

// NewJSONAPIEnvelope creates a new JSON:API envelope with the given data
// Supports both single data objects and collections (arrays)
func NewJSONAPIEnvelope(data interface{}) *JSONAPIEnvelope {
	return &JSONAPIEnvelope{Data: data.(JSONAPIData)}
}

// NewJSONAPIEnvelopeWithArray creates a new JSON:API envelope with an array of data objects
// For collections, the data field contains an array of JSONAPIData objects
func NewJSONAPIEnvelopeWithArray(data []JSONAPIData) *JSONAPIEnvelope {
	// Ensure we have a non-nil slice to preserve type information
	if data == nil {
		data = []JSONAPIData{}
	}
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
		ID:         strconv.FormatInt(project.ID, 10), // ID as string per JSON:API spec
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
