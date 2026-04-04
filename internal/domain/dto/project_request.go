package dto

// ProjectRequest represents the incoming JSON payload for project creation
type ProjectRequest struct {
	Name      string  `json:"name"`
	TotalPage int     `json:"total_page"`
	Page      int     `json:"page"`
	StartedAt *string `json:"started_at,omitempty"`
	Reinicia  bool    `json:"reinicia,omitempty"`
}
