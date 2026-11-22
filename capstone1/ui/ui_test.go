package ui_test

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"studentmgmt/service"
	"studentmgmt/store"
	"studentmgmt/ui"
)

func runSession(t *testing.T, svc *service.Service, input string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	err := ui.New(svc, strings.NewReader(input), &out).Run()
	return out.String(), err
}

func newTestService(t *testing.T) *service.Service {
	t.Helper()
	path := filepath.Join(t.TempDir(), "students.json")
	svc, err := service.New(store.NewFileStore(path))
	if err != nil {
		t.Fatalf("service.New() error: %v", err)
	}
	return svc
}

func TestUI_EndToEndCRUD(t *testing.T) {
	svc := newTestService(t)

	script := strings.Join([]string{
		// Empty list check.
		"2",
		// Invalid menu choice & cancel at main menu stays in loop.
		"9",
		"cancel",
		// Add student 101 with initial invalid ID, age, and mark re-prompts.
		"1",
		"abc",
		"101",
		"Aman Verma",
		"notanumber",
		"3",
		"21",
		"B.Tech CSE",
		"bad",
		"78",
		"85",
		"105",
		"90",
		"72",
		"88",
		// Duplicate ID then cancel.
		"1",
		"101",
		"cancel",
		// Add second student 102.
		"1",
		"102",
		"Riya Sharma",
		"20",
		"BCA",
		"95",
		"92",
		"98",
		"90",
		"96",
		// View all students.
		"2",
		// Search by ID (found and not found).
		"3", "1", "101",
		"3", "1", "999",
		// Search by name (empty re-prompt, match, no match, invalid search option).
		"3", "2", "", "aman",
		"3", "2", "zzz",
		"3", "9",
		// Update student 101: keep all fields with Enter except Subject 5 -> 95.
		"4", "101", "", "", "", "", "", "", "", "95",
		// Update unknown student.
		"4", "555",
		// Delete student 102: invalid answer, then 'n' abort, then 'y' confirm.
		"5", "102", "maybe", "n",
		"5", "102", "y",
		// Delete unknown student.
		"5", "999",
		// Exit cleanly.
		"6",
	}, "\n") + "\n"

	out, err := runSession(t, svc, script)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	checks := []string{
		"No student records yet.",
		"Invalid choice. Please enter a number from 1 to 6.",
		"Already at the main menu.",
		"ID must be a whole number",
		"age must be a whole number",
		"age must be between 5 and 100",
		"marks must be a number",
		"marks must be between 0 and 100",
		"413/500",
		"82.60%",
		"a student with this ID already exists",
		"Cancelled.",
		"471/500",
		"94.20%",
		"All Students (2)",
		"ID 999: student not found",
		"enter some text to search for",
		"1 match(es):",
		`No students match "zzz".`,
		"420/500",
		"84.00%",
		"ID 555: student not found",
		"please answer y or n",
		"Deletion aborted.",
		"Student deleted.",
		"Goodbye! All changes are saved.",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("UI output missing expected text %q", want)
		}
	}

	// Verify final persisted state has only student 101 with updated total 420.
	remaining := svc.List()
	if len(remaining) != 1 || remaining[0].ID != 101 || remaining[0].Total != 420 {
		t.Fatalf("final state = %+v, want only ID 101 with Total=420", remaining)
	}
}

func TestUI_EOFHandling(t *testing.T) {
	// EOF at main menu.
	out, err := runSession(t, newTestService(t), "")
	if err != nil {
		t.Fatalf("Run() on empty input error: %v", err)
	}
	if !strings.Contains(out, "Input closed. Goodbye!") {
		t.Errorf("expected EOF goodbye message, got %q", out)
	}

	// EOF mid-operation (during Add Student prompt).
	out, err = runSession(t, newTestService(t), "1\n101\n")
	if err != nil {
		t.Fatalf("Run() on mid-prompt EOF error: %v", err)
	}
	if !strings.Contains(out, "Input closed. Goodbye!") {
		t.Errorf("expected EOF goodbye message mid-prompt, got %q", out)
	}
}

type failingReader struct{ err error }

func (f failingReader) Read([]byte) (int, error) { return 0, f.err }

func TestUI_ReaderErrorPropagated(t *testing.T) {
	wantErr := errors.New("stream broken")
	var out bytes.Buffer
	err := ui.New(newTestService(t), failingReader{err: wantErr}, &out).Run()
	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}
