// Package store provides persistence abstractions and an atomic JSON FileStore.
package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"studentmgmt/student"
)

// ErrCorrupt indicates that the JSON data file cannot be decoded.
var ErrCorrupt = errors.New("data file is corrupt")

// Store defines the persistence contract required by the service layer.
type Store interface {
	Load() ([]student.Student, error)
	Save([]student.Student) error
}

// FileStore persists student records to a JSON file on disk using atomic temp-file replacement.
type FileStore struct {
	Path string
}

// NewFileStore returns a FileStore targeting path.
func NewFileStore(path string) *FileStore {
	return &FileStore{Path: path}
}

// Load reads all students from disk. If the file does not exist, an empty file is created.
// An empty or whitespace-only file is treated as an empty collection.
func (f *FileStore) Load() ([]student.Student, error) {
	data, err := os.ReadFile(f.Path)
	if errors.Is(err, os.ErrNotExist) {
		if err := f.Save(nil); err != nil {
			return nil, err
		}
		return []student.Student{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", f.Path, err)
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []student.Student{}, nil
	}
	var students []student.Student
	if err := json.Unmarshal(data, &students); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrCorrupt, f.Path, err)
	}
	if students == nil {
		students = []student.Student{}
	}
	return students, nil
}

// Save atomically writes students as indented JSON via a temporary file and rename.
func (f *FileStore) Save(students []student.Student) error {
	if students == nil {
		students = []student.Student{}
	}
	data, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		return fmt.Errorf("encode students: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(f.Path), filepath.Base(f.Path)+".tmp*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmp.Name(), f.Path); err != nil {
		return fmt.Errorf("replace %s: %w", f.Path, err)
	}
	return nil
}
