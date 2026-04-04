package dto

// LogRequest represents the incoming JSON payload for log creation
type LogRequest struct {
	StartPage int     `json:"start_page"`
	EndPage   int     `json:"end_page"`
	Data      *string `json:"data,omitempty"`
	Note      *string `json:"note,omitempty"`
	Text      *string `json:"text,omitempty"`
}
