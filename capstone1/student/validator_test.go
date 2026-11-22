package student

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{"positive", 101, nil},
		{"minimum positive", 1, nil},
		{"zero", 0, ErrInvalidID},
		{"negative", -7, ErrInvalidID},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateID(tc.id)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateID(%d) = %v, want %v", tc.id, err, tc.wantErr)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid name", "Aman Verma", nil},
		{"padded valid name", "  Riya  ", nil},
		{"empty", "", ErrEmptyName},
		{"spaces only", "   \t\n ", ErrEmptyName},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateName(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateName(%q) = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestValidateCourse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid course", "B.Tech CSE", nil},
		{"empty", "", ErrEmptyCourse},
		{"whitespace only", " \t ", ErrEmptyCourse},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCourse(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateCourse(%q) = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestValidateAge(t *testing.T) {
	tests := []struct {
		name    string
		age     int
		wantErr error
	}{
		{"below min", MinAge - 1, ErrInvalidAge},
		{"zero", 0, ErrInvalidAge},
		{"negative", -10, ErrInvalidAge},
		{"exact min", MinAge, nil},
		{"typical", 21, nil},
		{"exact max", MaxAge, nil},
		{"above max", MaxAge + 1, ErrInvalidAge},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAge(tc.age)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateAge(%d) = %v, want %v", tc.age, err, tc.wantErr)
			}
		})
	}
}

func TestValidateMark(t *testing.T) {
	tests := []struct {
		name    string
		mark    float64
		wantErr error
	}{
		{"zero", 0, nil},
		{"mid decimal", 75.5, nil},
		{"exact max", MaxPerSubject, nil},
		{"negative", -0.01, ErrInvalidMarks},
		{"above max", MaxPerSubject + 0.01, ErrInvalidMarks},
		{"NaN", math.NaN(), ErrInvalidMarks},
		{"positive Inf", math.Inf(1), ErrInvalidMarks},
		{"negative Inf", math.Inf(-1), ErrInvalidMarks},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMark(tc.mark, MaxPerSubject)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateMark(%v) = %v, want %v", tc.mark, err, tc.wantErr)
			}
		})
	}
}

func TestValidateMarks(t *testing.T) {
	tests := []struct {
		name        string
		marks       []float64
		wantErr     error
		wantContain string
	}{
		{"valid 5 marks", []float64{78, 85, 90, 72, 88}, nil, ""},
		{"nil slice", nil, ErrMarksCount, "expected 5, got 0"},
		{"too few", []float64{80, 85, 90}, ErrMarksCount, "expected 5, got 3"},
		{"too many", []float64{80, 85, 90, 70, 60, 50}, ErrMarksCount, "expected 5, got 6"},
		{"invalid mark at subject 3", []float64{95, 92, 105, 90, 96}, ErrInvalidMarks, "subject 3"},
		{"NaN at subject 1", []float64{math.NaN(), 80, 80, 80, 80}, ErrInvalidMarks, "subject 1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMarks(tc.marks, NumSubjects, MaxPerSubject)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ValidateMarks(%v) = %v, want %v", tc.marks, err, tc.wantErr)
			}
			if tc.wantContain != "" && !strings.Contains(err.Error(), tc.wantContain) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantContain)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	valid := Student{
		ID:     101,
		Name:   "Aman Verma",
		Age:    21,
		Course: "B.Tech CSE",
		Marks:  []float64{78, 85, 90, 72, 88},
	}

	tests := []struct {
		name    string
		mutate  func(*Student)
		wantErr error
	}{
		{"valid record", func(s *Student) {}, nil},
		{"invalid ID", func(s *Student) { s.ID = 0 }, ErrInvalidID},
		{"empty name", func(s *Student) { s.Name = "   " }, ErrEmptyName},
		{"invalid age", func(s *Student) { s.Age = 3 }, ErrInvalidAge},
		{"empty course", func(s *Student) { s.Course = "" }, ErrEmptyCourse},
		{"wrong marks count", func(s *Student) { s.Marks = []float64{80, 90} }, ErrMarksCount},
		{"out-of-range mark", func(s *Student) { s.Marks[4] = 101 }, ErrInvalidMarks},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := valid
			st.Marks = append([]float64(nil), valid.Marks...)
			tc.mutate(&st)
			err := Validate(st)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Validate(%+v) = %v, want %v", st, err, tc.wantErr)
			}
		})
	}
}
