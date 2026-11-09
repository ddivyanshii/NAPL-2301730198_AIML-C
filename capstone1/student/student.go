package student

const (
	NumSubjects   = 5
	MaxPerSubject = 100.0
)

type Student struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Age        int       `json:"age"`
	Course     string    `json:"course"`
	Marks      []float64 `json:"marks"`
	Total      float64   `json:"total"`
	Percentage float64   `json:"percentage"`
	Grade      string    `json:"grade"`
}
