package service

import (
	"errors"
	"reflect"
	"testing"

	"studentmgmt/student"
)

// mockStore is an in-memory store.Store test double with injectable errors.
type mockStore struct {
	data      []student.Student
	loadErr   error
	saveErr   error
	saveCalls int
}

func (m *mockStore) Load() ([]student.Student, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	out := make([]student.Student, len(m.data))
	for i, s := range m.data {
		s.Marks = append([]float64(nil), s.Marks...)
		out[i] = s
	}
	return out, nil
}

func (m *mockStore) Save(list []student.Student) error {
	m.saveCalls++
	if m.saveErr != nil {
		return m.saveErr
	}
	out := make([]student.Student, len(list))
	for i, s := range list {
		s.Marks = append([]float64(nil), s.Marks...)
		out[i] = s
	}
	m.data = out
	return nil
}

func sampleStudent(id int, name string, marks []float64) student.Student {
	return student.Student{
		ID:     id,
		Name:   name,
		Age:    20,
		Course: "B.Tech CSE",
		Marks:  marks,
	}
}

func TestNew(t *testing.T) {
	t.Run("loads and sorts records by ID", func(t *testing.T) {
		ms := &mockStore{
			data: []student.Student{
				sampleStudent(103, "Aman Gupta", []float64{55, 48, 62, 51, 58}),
				sampleStudent(101, "Aman Verma", []float64{78, 85, 90, 72, 88}),
			},
		}
		svc, err := New(ms)
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}
		got := svc.List()
		if len(got) != 2 || got[0].ID != 101 || got[1].ID != 103 {
			t.Fatalf("List() order = %+v, want IDs [101, 103]", got)
		}
	})

	t.Run("propagates store load error", func(t *testing.T) {
		wantErr := errors.New("disk read failed")
		if _, err := New(&mockStore{loadErr: wantErr}); !errors.Is(err, wantErr) {
			t.Fatalf("New() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("rejects duplicate IDs in stored data", func(t *testing.T) {
		ms := &mockStore{
			data: []student.Student{
				sampleStudent(101, "First", []float64{80, 80, 80, 80, 80}),
				sampleStudent(101, "Second", []float64{90, 90, 90, 90, 90}),
			},
		}
		if _, err := New(ms); !errors.Is(err, ErrDuplicateID) {
			t.Fatalf("New() error = %v, want %v", err, ErrDuplicateID)
		}
	})

	t.Run("rejects invalid record in stored data", func(t *testing.T) {
		bad := sampleStudent(101, "Bad Age", []float64{80, 80, 80, 80, 80})
		bad.Age = 2
		if _, err := New(&mockStore{data: []student.Student{bad}}); !errors.Is(err, student.ErrInvalidAge) {
			t.Fatalf("New() error = %v, want %v", err, student.ErrInvalidAge)
		}
	})
}

func TestAdd(t *testing.T) {
	ms := &mockStore{}
	svc, err := New(ms)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	// Add valid record with surrounding spaces to verify trimming & computation.
	in := student.Student{
		ID:     101,
		Name:   "  Aman Verma  ",
		Age:    21,
		Course: "  B.Tech CSE ",
		Marks:  []float64{78, 85, 90, 72, 88},
	}
	saved, err := svc.Add(in)
	if err != nil {
		t.Fatalf("Add() error: %v", err)
	}
	if saved.Name != "Aman Verma" || saved.Course != "B.Tech CSE" {
		t.Errorf("expected trimmed fields, got Name=%q Course=%q", saved.Name, saved.Course)
	}
	if saved.Total != 413 || saved.Percentage != 82.6 || saved.Grade != "A" {
		t.Errorf("computed = (%v, %v, %q), want (413, 82.6, A)", saved.Total, saved.Percentage, saved.Grade)
	}
	if ms.saveCalls != 1 || len(ms.data) != 1 {
		t.Errorf("expected 1 save call with 1 record, got calls=%d len=%d", ms.saveCalls, len(ms.data))
	}

	// Mutating caller's input or returned Marks must not mutate service state.
	in.Marks[0] = 0
	saved.Marks[1] = 0
	got, _ := svc.SearchByID(101)
	if !reflect.DeepEqual(got.Marks, []float64{78, 85, 90, 72, 88}) {
		t.Errorf("internal Marks mutated: %v", got.Marks)
	}

	// Duplicate ID is rejected.
	if _, err := svc.Add(sampleStudent(101, "Duplicate", []float64{50, 50, 50, 50, 50})); !errors.Is(err, ErrDuplicateID) {
		t.Errorf("duplicate Add() error = %v, want %v", err, ErrDuplicateID)
	}

	// Validation failure is rejected without calling Save.
	prevCalls := ms.saveCalls
	bad := sampleStudent(102, "", []float64{50, 50, 50, 50, 50})
	if _, err := svc.Add(bad); !errors.Is(err, student.ErrEmptyName) {
		t.Errorf("invalid Add() error = %v, want %v", err, student.ErrEmptyName)
	}
	if ms.saveCalls != prevCalls {
		t.Errorf("Save was called on validation failure")
	}

	// Store save failure leaves in-memory state unchanged.
	saveErr := errors.New("disk full")
	ms.saveErr = saveErr
	if _, err := svc.Add(sampleStudent(102, "Riya Sharma", []float64{95, 92, 98, 90, 96})); !errors.Is(err, saveErr) {
		t.Errorf("Add() save error = %v, want %v", err, saveErr)
	}
	if len(svc.List()) != 1 {
		t.Errorf("in-memory list changed after failed save: len=%d", len(svc.List()))
	}
	if _, err := svc.SearchByID(102); !errors.Is(err, ErrNotFound) {
		t.Errorf("index changed after failed save: SearchByID(102) = %v", err)
	}
}

func TestSearchAndList(t *testing.T) {
	ms := &mockStore{}
	svc, _ := New(ms)

	_, _ = svc.Add(sampleStudent(103, "Aman Gupta", []float64{55, 48, 62, 51, 58}))
	_, _ = svc.Add(sampleStudent(101, "Aman Verma", []float64{78, 85, 90, 72, 88}))
	_, _ = svc.Add(sampleStudent(102, "Riya Sharma", []float64{95, 92, 98, 90, 96}))

	// List returns records sorted by ascending ID.
	all := svc.List()
	if len(all) != 3 || all[0].ID != 101 || all[1].ID != 102 || all[2].ID != 103 {
		t.Fatalf("List() IDs not sorted: %+v", all)
	}

	// Mutating returned slice does not affect internal state.
	all[0].Name = "Tampered"
	all[0].Marks[0] = -99
	fresh, _ := svc.SearchByID(101)
	if fresh.Name != "Aman Verma" || fresh.Marks[0] != 78 {
		t.Errorf("List() did not return a defensive copy: %+v", fresh)
	}

	// SearchByID not found.
	if _, err := svc.SearchByID(999); !errors.Is(err, ErrNotFound) {
		t.Errorf("SearchByID(999) error = %v, want %v", err, ErrNotFound)
	}

	// SearchByName: partial, case-insensitive, whitespace-trimmed.
	matches := svc.SearchByName("  AMAN ")
	if len(matches) != 2 || matches[0].ID != 101 || matches[1].ID != 103 {
		t.Errorf("SearchByName(\"  AMAN \") = %+v, want IDs [101, 103]", matches)
	}
	if res := svc.SearchByName("   "); len(res) != 0 {
		t.Errorf("SearchByName blank = %+v, want empty", res)
	}
	if res := svc.SearchByName("zzz"); len(res) != 0 {
		t.Errorf("SearchByName(\"zzz\") = %+v, want empty", res)
	}
}

func TestUpdate(t *testing.T) {
	ms := &mockStore{}
	svc, _ := New(ms)
	_, _ = svc.Add(sampleStudent(101, "Aman Verma", []float64{78, 85, 90, 72, 88}))

	// Update unknown ID -> ErrNotFound.
	if _, err := svc.Update(555, sampleStudent(555, "Nobody", []float64{50, 50, 50, 50, 50})); !errors.Is(err, ErrNotFound) {
		t.Errorf("Update(555) error = %v, want %v", err, ErrNotFound)
	}

	// Update with invalid field -> validation error, original unchanged.
	bad := sampleStudent(101, "Aman Verma", []float64{78, 85, 90, 72, 105})
	if _, err := svc.Update(101, bad); !errors.Is(err, student.ErrInvalidMarks) {
		t.Errorf("invalid Update() error = %v, want %v", err, student.ErrInvalidMarks)
	}

	// Update with mismatched upd.ID still preserves target ID 101 and recalculates results.
	upd := student.Student{
		ID:     999,
		Name:   "Aman Verma",
		Age:    22,
		Course: "B.Tech CSE",
		Marks:  []float64{78, 85, 90, 72, 95},
	}
	got, err := svc.Update(101, upd)
	if err != nil {
		t.Fatalf("Update(101) error: %v", err)
	}
	if got.ID != 101 || got.Age != 22 || got.Total != 420 || got.Percentage != 84 || got.Grade != "A" {
		t.Errorf("Update(101) = %+v, want ID=101 Age=22 Total=420 Pct=84 Grade=A", got)
	}

	// Failed save leaves previous values intact.
	saveErr := errors.New("write error")
	ms.saveErr = saveErr
	_, err = svc.Update(101, sampleStudent(101, "Should Not Stick", []float64{10, 10, 10, 10, 10}))
	if !errors.Is(err, saveErr) {
		t.Errorf("Update() save error = %v, want %v", err, saveErr)
	}
	after, _ := svc.SearchByID(101)
	if after.Name != "Aman Verma" || after.Total != 420 {
		t.Errorf("in-memory record changed after failed save: %+v", after)
	}
}

func TestDelete(t *testing.T) {
	ms := &mockStore{}
	svc, _ := New(ms)
	_, _ = svc.Add(sampleStudent(101, "Aman Verma", []float64{78, 85, 90, 72, 88}))
	_, _ = svc.Add(sampleStudent(102, "Riya Sharma", []float64{95, 92, 98, 90, 96}))
	_, _ = svc.Add(sampleStudent(103, "Aman Gupta", []float64{55, 48, 62, 51, 58}))

	// Unknown ID -> ErrNotFound.
	if err := svc.Delete(999); !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete(999) error = %v, want %v", err, ErrNotFound)
	}

	// Failed save leaves all records intact.
	saveErr := errors.New("io error")
	ms.saveErr = saveErr
	if err := svc.Delete(102); !errors.Is(err, saveErr) {
		t.Errorf("Delete(102) save error = %v, want %v", err, saveErr)
	}
	if len(svc.List()) != 3 {
		t.Fatalf("expected 3 records after failed delete, got %d", len(svc.List()))
	}

	// Successful delete removes 102 and rebuilds index for 103.
	ms.saveErr = nil
	if err := svc.Delete(102); err != nil {
		t.Fatalf("Delete(102) unexpected error: %v", err)
	}
	if _, err := svc.SearchByID(102); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected 102 to be deleted, got err=%v", err)
	}
	st103, err := svc.SearchByID(103)
	if err != nil || st103.Name != "Aman Gupta" {
		t.Errorf("index broken for remaining record 103: st=%+v err=%v", st103, err)
	}
	// Verify updating 103 works after 102 was removed ahead of it.
	if _, err := svc.Update(103, sampleStudent(103, "Aman Gupta Jr", []float64{60, 60, 60, 60, 60})); err != nil {
		t.Errorf("Update(103) after Delete(102) failed: %v", err)
	}
}
