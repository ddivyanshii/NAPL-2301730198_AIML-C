package student

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	// MinAge is the minimum valid student age.
	MinAge = 5
	// MaxAge is the maximum valid student age.
	MaxAge = 100
)

// Sentinel validation errors returned by field-level and record-level checks.
var (
	ErrInvalidID    = errors.New("ID must be a positive whole number")
	ErrEmptyName    = errors.New("name must not be empty")
	ErrEmptyCourse  = errors.New("course must not be empty")
	ErrInvalidAge   = fmt.Errorf("age must be between %d and %d", MinAge, MaxAge)
	ErrInvalidMarks = errors.New("invalid marks")
	ErrMarksCount   = errors.New("wrong number of subject marks")
)

// ValidateID checks that id is a positive integer.
func ValidateID(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return nil
}

// ValidateName checks that name contains non-whitespace characters.
func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}
	return nil
}

// ValidateCourse checks that course contains non-whitespace characters.
func ValidateCourse(course string) error {
	if strings.TrimSpace(course) == "" {
		return ErrEmptyCourse
	}
	return nil
}

// ValidateAge checks that age lies within [MinAge, MaxAge].
func ValidateAge(age int) error {
	if age < MinAge || age > MaxAge {
		return ErrInvalidAge
	}
	return nil
}

// ValidateMark checks that a single subject mark is finite and within [0, max].
func ValidateMark(m, max float64) error {
	if math.IsNaN(m) || math.IsInf(m, 0) || m < 0 || m > max {
		return fmt.Errorf("%w: marks must be between 0 and %g", ErrInvalidMarks, max)
	}
	return nil
}

// ValidateMarks checks that marks has length n and every entry lies within [0, max].
func ValidateMarks(marks []float64, n int, max float64) error {
	if len(marks) != n {
		return fmt.Errorf("%w: expected %d, got %d", ErrMarksCount, n, len(marks))
	}
	for i, m := range marks {
		if err := ValidateMark(m, max); err != nil {
			return fmt.Errorf("subject %d: %w", i+1, err)
		}
	}
	return nil
}

// Validate checks every input field of s (derived fields are ignored).
func Validate(s Student) error {
	if err := ValidateID(s.ID); err != nil {
		return err
	}
	if err := ValidateName(s.Name); err != nil {
		return err
	}
	if err := ValidateAge(s.Age); err != nil {
		return err
	}
	if err := ValidateCourse(s.Course); err != nil {
		return err
	}
	return ValidateMarks(s.Marks, NumSubjects, MaxPerSubject)
}
