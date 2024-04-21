package internal

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

type NewTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

const (
	Low    = "Low"
	Medium = "Medium"
	High   = "High"
)

const (
	Open  = "Open"
	Doing = "Doing"
	Done  = "Done"
)
