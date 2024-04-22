package model

type (
	Todo struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Severity    string `json:"severity"`
		Status      string `json:"status"`
	}

	NewTodoRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Severity    string `json:"severity"`
		Status      string `json:"status"`
	}
)

const (
	SeverityLow    = "LOW"
	SeverityMedium = "MED"
	SeverityHigh   = "HI"
)

const (
	StatusOpen  = "OPEN"
	StatusDoing = "DOING"
	StatusDone  = "DONE"
)
