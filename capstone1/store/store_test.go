package store_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"studentmgmt/service"
	"studentmgmt/store"
	"studentmgmt/student"
)

func TestFileStore_MissingFileCreated(t *testing.T) {
	path := filepath.Join(t.TempDir(), "students.json")
	fs := store.NewFileStore(path)

	got, err := fs.Load()
	if err != nil {
		t.Fatalf("Load() on missing file error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Load() = %+v, want empty slice", got)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected missing file to be created, stat error: %v", err)
	}
}

func TestFileStore_EmptyAndNullFiles(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"zero-byte file", ""},
		{"whitespace-only file", "   \n\t  "},
		{"empty json array", "[]"},
		{"json null", "null"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "students.json")
			if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
				t.Fatalf("setup WriteFile: %v", err)
			}
			got, err := store.NewFileStore(path).Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}
			if got == nil || len(got) != 0 {
				t.Errorf("Load() = %#v, want non-nil empty slice", got)
			}
		})
	}
}

func TestFileStore_CorruptJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "students.json")
	if err := os.WriteFile(path, []byte("{not valid json"), 0o600); err != nil {
		t.Fatalf("setup WriteFile: %v", err)
	}
	if _, err := store.NewFileStore(path).Load(); !errors.Is(err, store.ErrCorrupt) {
		t.Fatalf("Load() on corrupt JSON error = %v, want %v", err, store.ErrCorrupt)
	}
}

func TestFileStore_IOErrors(t *testing.T) {
	// Saving into a non-existent directory fails cleanly.
	badPath := filepath.Join(t.TempDir(), "no-such-dir", "students.json")
	if err := store.NewFileStore(badPath).Save(nil); err == nil {
		t.Errorf("Save() into missing directory succeeded, want error")
	}
	// Loading from a directory path fails cleanly (not ErrCorrupt).
	dirPath := t.TempDir()
	if _, err := store.NewFileStore(dirPath).Load(); err == nil || errors.Is(err, store.ErrCorrupt) {
		t.Errorf("Load() on directory path = %v, want non-corrupt read error", err)
	}
}

// TestFileStore_PersistenceAcrossRestart verifies that 100% of saved records (12 students)
// are restored identically after simulating an application restart, with no leftover temp files.
func TestFileStore_PersistenceAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "students.json")

	svc1, err := service.New(store.NewFileStore(path))
	if err != nil {
		t.Fatalf("initial service.New() error: %v", err)
	}

	inputs := []student.Student{
		{ID: 101, Name: "Aman Verma", Age: 21, Course: "B.Tech CSE", Marks: []float64{78, 85, 90, 72, 88}},
		{ID: 102, Name: "Riya Sharma", Age: 20, Course: "BCA", Marks: []float64{95, 92, 98, 90, 96}},
		{ID: 103, Name: "Aman Gupta", Age: 22, Course: "B.Sc", Marks: []float64{55, 48, 62, 51, 58}},
		{ID: 104, Name: "Neha Kapoor", Age: 19, Course: "B.Tech AI", Marks: []float64{90, 90, 90, 90, 90}},
		{ID: 105, Name: "Kabir Mehta", Age: 21, Course: "B.Tech CSE", Marks: []float64{80, 80, 80, 80, 80}},
		{ID: 106, Name: "Priya Nair", Age: 20, Course: "MCA", Marks: []float64{72, 75, 68, 74, 71}},
		{ID: 107, Name: "Rohan Joshi", Age: 23, Course: "B.Tech ECE", Marks: []float64{70, 70, 70, 70, 70}},
		{ID: 108, Name: "Simran Kaur", Age: 18, Course: "BBA", Marks: []float64{60, 65, 62, 68, 64}},
		{ID: 109, Name: "Dev Patel", Age: 21, Course: "B.Sc CS", Marks: []float64{60, 60, 60, 60, 60}},
		{ID: 110, Name: "Ananya Iyer", Age: 20, Course: "BA", Marks: []float64{50, 50, 50, 50, 50}},
		{ID: 111, Name: "Vikram Singh", Age: 22, Course: "B.Com", Marks: []float64{49.5, 50, 49.5, 50, 50}},
		{ID: 112, Name: "Meera Das", Age: 19, Course: "B.Tech CSE", Marks: []float64{100, 99.5, 98.5, 100, 97.5}},
	}

	for _, st := range inputs {
		if _, err := svc1.Add(st); err != nil {
			t.Fatalf("Add(%d) error: %v", st.ID, err)
		}
	}
	want := svc1.List()

	// Simulate application restart by creating a brand-new FileStore and Service.
	svc2, err := service.New(store.NewFileStore(path))
	if err != nil {
		t.Fatalf("restarted service.New() error: %v", err)
	}
	got := svc2.List()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("restored records mismatch:\ngot:  %+v\nwant: %+v", got, want)
	}

	// Ensure atomic save left no temporary files behind.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".tmp") {
			t.Errorf("leftover temp file found: %s", e.Name())
		}
	}
}
