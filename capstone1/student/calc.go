package student

import "math"

type gradeBand struct {
	min   float64
	grade string
}

var gradeBands = []gradeBand{
	{90, "O"},
	{80, "A"},
	{70, "B"},
	{60, "C"},
	{50, "D"},
}

// FailGrade is assigned when a student's percentage is below the lowest passing band.
const FailGrade = "F"

// Total returns the sum of marks rounded to two decimal places.
func Total(marks []float64) float64 {
	var t float64
	for _, m := range marks {
		t += m
	}
	return round2(t)
}

// Percentage returns (total / max) * 100 rounded to two decimal places, or 0 if max <= 0.
func Percentage(total, max float64) float64 {
	if max <= 0 {
		return 0
	}
	return round2(total / max * 100)
}

// Grade maps a percentage to a letter grade using gradeBands.
func Grade(percentage float64) string {
	for _, b := range gradeBands {
		if percentage >= b.min {
			return b.grade
		}
	}
	return FailGrade
}

// Compute calculates total marks, percentage, and grade for the given subject marks.
func Compute(marks []float64, maxPerSubject float64) (total, percentage float64, grade string) {
	total = Total(marks)
	percentage = Percentage(total, float64(len(marks))*maxPerSubject)
	return total, percentage, Grade(percentage)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
