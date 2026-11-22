package student

import (
	"math"
	"testing"
)

func TestTotal(t *testing.T) {
	tests := []struct {
		name  string
		marks []float64
		want  float64
	}{
		{"empty slice", nil, 0},
		{"all zeros", []float64{0, 0, 0, 0, 0}, 0},
		{"all perfect", []float64{100, 100, 100, 100, 100}, 500},
		{"typical integers", []float64{78, 85, 90, 72, 88}, 413},
		{"decimal marks", []float64{80.25, 75.5, 91.25, 68.0, 84.5}, 399.5},
		{"floating point rounding", []float64{0.1, 0.2, 0.3, 0.1, 0.2}, 0.9},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Total(tc.marks); got != tc.want {
				t.Errorf("Total(%v) = %v, want %v", tc.marks, got, tc.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		name  string
		total float64
		max   float64
		want  float64
	}{
		{"zero max returns 0", 400, 0, 0},
		{"negative max returns 0", 400, -500, 0},
		{"zero total", 0, 500, 0},
		{"full score", 500, 500, 100},
		{"exact decimal", 413, 500, 82.6},
		{"repeating decimal rounded to 2 places", 200, 300, 66.67},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Percentage(tc.total, tc.max); got != tc.want {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tc.total, tc.max, got, tc.want)
			}
		})
	}
}

func TestGrade(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		want       string
	}{
		{"perfect 100", 100, "O"},
		{"boundary O lower", 90, "O"},
		{"just below O", 89.99, "A"},
		{"boundary A lower", 80, "A"},
		{"just below A", 79.99, "B"},
		{"boundary B lower", 70, "B"},
		{"just below B", 69.99, "C"},
		{"boundary C lower", 60, "C"},
		{"just below C", 59.99, "D"},
		{"boundary D lower", 50, "D"},
		{"just below D", 49.99, "F"},
		{"zero percent", 0, "F"},
		{"negative percent", -5, "F"},
		{"NaN defaults to F", math.NaN(), "F"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Grade(tc.percentage); got != tc.want {
				t.Errorf("Grade(%v) = %q, want %q", tc.percentage, got, tc.want)
			}
		})
	}
}

// TestCompute_SampleDataset verifies 100% agreement between Compute and hand-computed
// results across a test set of 12 sample students spanning every grade band and boundary.
func TestCompute_SampleDataset(t *testing.T) {
	tests := []struct {
		name     string
		marks    []float64
		wantTot  float64
		wantPct  float64
		wantGrad string
	}{
		{"101 Aman Verma", []float64{78, 85, 90, 72, 88}, 413, 82.6, "A"},
		{"102 Riya Sharma", []float64{95, 92, 98, 90, 96}, 471, 94.2, "O"},
		{"103 Aman Gupta", []float64{55, 48, 62, 51, 58}, 274, 54.8, "D"},
		{"104 Neha Kapoor (O boundary)", []float64{90, 90, 90, 90, 90}, 450, 90, "O"},
		{"105 Kabir Mehta (A boundary)", []float64{80, 80, 80, 80, 80}, 400, 80, "A"},
		{"106 Priya Nair", []float64{72, 75, 68, 74, 71}, 360, 72, "B"},
		{"107 Rohan Joshi (B boundary)", []float64{70, 70, 70, 70, 70}, 350, 70, "B"},
		{"108 Simran Kaur", []float64{60, 65, 62, 68, 64}, 319, 63.8, "C"},
		{"109 Dev Patel (C boundary)", []float64{60, 60, 60, 60, 60}, 300, 60, "C"},
		{"110 Ananya Iyer (D boundary)", []float64{50, 50, 50, 50, 50}, 250, 50, "D"},
		{"111 Vikram Singh (just below D)", []float64{49.5, 50, 49.5, 50, 50}, 249, 49.8, "F"},
		{"112 Meera Das (decimal marks)", []float64{100, 99.5, 98.5, 100, 97.5}, 495.5, 99.1, "O"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTot, gotPct, gotGrad := Compute(tc.marks, MaxPerSubject)
			if gotTot != tc.wantTot || gotPct != tc.wantPct || gotGrad != tc.wantGrad {
				t.Errorf("Compute(%v) = (%v, %v, %q), want (%v, %v, %q)",
					tc.marks, gotTot, gotPct, gotGrad, tc.wantTot, tc.wantPct, tc.wantGrad)
			}
		})
	}
}
