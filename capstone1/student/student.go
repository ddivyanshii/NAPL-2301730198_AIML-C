// Package student defines the Student model, input validation rules, and result calculation.
package student

const (
	// NumSubjects is the fixed number of subjects recorded per student.
	NumSubjects = 5
	// MaxPerSubject is the maximum score allowed for a single subject.
	MaxPerSubject = 100.0
)

// Student holds a student's profile, subject marks, and computed academic results.
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
