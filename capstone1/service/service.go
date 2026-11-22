// Package service implements student CRUD business rules over an in-memory index and a Store.
package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"studentmgmt/store"
	"studentmgmt/student"
)

// Sentinel errors returned by Service operations.
var (
	ErrDuplicateID = errors.New("a student with this ID already exists")
	ErrNotFound    = errors.New("student not found")
)

// Service manages student records in ascending ID order with O(1) ID lookup.
type Service struct {
	store    store.Store
	students []student.Student
	index    map[int]int
}

// New initializes a Service by loading and validating existing records from st.
func New(st store.Store) (*Service, error) {
	list, err := st.Load()
	if err != nil {
		return nil, err
	}
	copied := make([]student.Student, len(list))
	for i, rec := range list {
		if err := student.Validate(rec); err != nil {
			return nil, fmt.Errorf("record %d (ID %d): %w", i+1, rec.ID, err)
		}
		copied[i] = clone(rec)
	}
	s := &Service{store: st}
	if err := s.replace(copied); err != nil {
		return nil, err
	}
	return s, nil
}

// Add validates st, computes its results, persists the updated list, and returns the saved record.
func (s *Service) Add(st student.Student) (student.Student, error) {
	if err := student.Validate(st); err != nil {
		return student.Student{}, err
	}
	if _, ok := s.index[st.ID]; ok {
		return student.Student{}, fmt.Errorf("ID %d: %w", st.ID, ErrDuplicateID)
	}
	st = withResults(st)
	if err := s.commit(append(s.snapshot(), st)); err != nil {
		return student.Student{}, err
	}
	return clone(st), nil
}

// List returns a defensive copy of all students ordered by ascending ID.
func (s *Service) List() []student.Student {
	return s.snapshot()
}

// SearchByID returns the student with the given id, or a wrapped ErrNotFound.
func (s *Service) SearchByID(id int) (student.Student, error) {
	i, ok := s.index[id]
	if !ok {
		return student.Student{}, fmt.Errorf("ID %d: %w", id, ErrNotFound)
	}
	return clone(s.students[i]), nil
}

// SearchByName returns all students whose names contain keyword (case-insensitive).
func (s *Service) SearchByName(keyword string) []student.Student {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	var out []student.Student
	if keyword == "" {
		return out
	}
	for _, st := range s.students {
		if strings.Contains(strings.ToLower(st.Name), keyword) {
			out = append(out, clone(st))
		}
	}
	return out
}

// Update replaces the student with the given id, recalculates results, and persists the change.
func (s *Service) Update(id int, upd student.Student) (student.Student, error) {
	i, ok := s.index[id]
	if !ok {
		return student.Student{}, fmt.Errorf("ID %d: %w", id, ErrNotFound)
	}
	upd.ID = id
	if err := student.Validate(upd); err != nil {
		return student.Student{}, err
	}
	upd = withResults(upd)
	list := s.snapshot()
	list[i] = upd
	if err := s.commit(list); err != nil {
		return student.Student{}, err
	}
	return clone(upd), nil
}

// Delete removes the student with the given id and persists the change.
func (s *Service) Delete(id int) error {
	i, ok := s.index[id]
	if !ok {
		return fmt.Errorf("ID %d: %w", id, ErrNotFound)
	}
	list := s.snapshot()
	list = append(list[:i], list[i+1:]...)
	return s.commit(list)
}

func (s *Service) commit(list []student.Student) error {
	sort.Slice(list, func(a, b int) bool { return list[a].ID < list[b].ID })
	if err := s.store.Save(list); err != nil {
		return fmt.Errorf("saving changes: %w", err)
	}
	return s.applySorted(list)
}

func (s *Service) replace(list []student.Student) error {
	sort.Slice(list, func(a, b int) bool { return list[a].ID < list[b].ID })
	return s.applySorted(list)
}

func (s *Service) applySorted(list []student.Student) error {
	idx := make(map[int]int, len(list))
	for i, st := range list {
		if _, dup := idx[st.ID]; dup {
			return fmt.Errorf("ID %d: %w", st.ID, ErrDuplicateID)
		}
		idx[st.ID] = i
	}
	s.students, s.index = list, idx
	return nil
}

func (s *Service) snapshot() []student.Student {
	out := make([]student.Student, len(s.students))
	for i, st := range s.students {
		out[i] = clone(st)
	}
	return out
}

func withResults(st student.Student) student.Student {
	st = clone(st)
	st.Name = strings.TrimSpace(st.Name)
	st.Course = strings.TrimSpace(st.Course)
	st.Total, st.Percentage, st.Grade = student.Compute(st.Marks, student.MaxPerSubject)
	return st
}

func clone(st student.Student) student.Student {
	st.Marks = append([]float64(nil), st.Marks...)
	return st
}
