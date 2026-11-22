# Console-Based Student Management System in Go

**ENSP465: New Age Programming Languages Lab · Capstone Project (Phases I–III Complete)**
Divyanshi · Roll No. `2301730198` · Semester VII (B.Tech CSE) · Section C · Faculty: Dr. Hemraj Kumawat

A lightweight, offline, menu-driven console application written in Go for managing student records. It supports full CRUD operations, automatic calculation of total marks, percentage, and letter grade, strict input validation, and atomic JSON file persistence using **only the Go standard library**.

---

## Features

- **Add student:** Record ID, name, age, course, and marks for 5 subjects (`0–100` each) with immediate field-level validation and duplicate-ID detection.
- **View all students:** Display all records in an aligned table (`text/tabwriter`), automatically sorted in ascending ID order.
- **Search student:** Look up an exact student by ID (`O(1)` map index) or find students by partial, case-insensitive name match.
- **Update student:** Edit any field while pressing `Enter` to keep current values; total, percentage, and grade are recalculated automatically.
- **Delete student:** Preview the record and confirm (`y/n`) before removal.
- **Cancel anywhere:** Type `cancel` at any prompt to abort the current operation and return cleanly to the main menu.
- **Atomic JSON persistence:** Automatically loads `students.json` on start-up (creating it if absent) and saves after every change via a temporary file + atomic rename; failed saves never corrupt in-memory state.

## Grading Scheme

| Percentage | Grade | Meaning |
|------------|-------|---------|
| 90 and above | **O** | Outstanding |
| 80 – 89.99 | **A** | Excellent |
| 70 – 79.99 | **B** | Good |
| 60 – 69.99 | **C** | Satisfactory |
| 50 – 59.99 | **D** | Pass |
| Below 50 | **F** | Fail |

## Prerequisites

- **Go 1.21 or newer** (`go version` to verify). No third-party packages, database server, or internet connection required.

## Build & Run

```bash
# Run directly
go run .

# Or build a standalone binary
go build -o studentmgmt .
./studentmgmt

# Use a custom JSON data file path
./studentmgmt -data mydata.json
```

## Testing & Verification

Run the full automated test suite (unit tests, persistence tests, and scripted UI tests):

```bash
# Run all tests with race detection and coverage summary
go test -v -race -cover ./...

# Generate a function-level coverage breakdown
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

# Static analysis and formatting checks
go vet ./...
gofmt -l .
```

### Coverage Summary

| Package | Responsibility | Statement Coverage |
|---------|----------------|-------------------|
| `studentmgmt/student` | Data model, `validator.go`, `calc.go` | **100.0%** |
| `studentmgmt/service` | CRUD rules, `map[int]int` index, save-first semantics | **100.0%** |
| `studentmgmt/ui` | Menu loop, prompts, validation feedback, table rendering | **93.8%** |
| `studentmgmt/store` | `Store` interface & atomic JSON `FileStore` | **78.9%** |
| **Overall** | Entire module (including `main.go`) | **91.2%** |

All core targets from Phase I (≥ 70% coverage on calculator, validator, and service; 100% agreement on ≥ 10 hand-computed sample students; 100% persistence across restart) are exceeded.

## Project Layout

```text
capstone1/
├── main.go                       # Entry point: parses -data flag and wires store → service → ui
├── go.mod                        # Go module definition (stdlib only)
├── student/
│   ├── student.go                # Student struct with JSON tags and constants
│   ├── calc.go                   # Total, Percentage, Grade, Compute
│   ├── calc_test.go              # Table-driven boundary tests + 12-student dataset verification
│   ├── validator.go              # Field-level and whole-record validation + sentinel errors
│   └── validator_test.go         # Table-driven validation tests (including NaN/Inf)
├── service/
│   ├── service.go                # In-memory sorted slice + map[int]int index + CRUD operations
│   └── service_test.go           # Unit tests using an in-memory mock Store
├── store/
│   ├── store.go                  # Store interface + atomic JSON FileStore
│   └── store_test.go             # Missing/empty/corrupt file & 12-student restart persistence tests
├── ui/
│   ├── ui.go                     # Menu-driven console interface with injected io.Reader/io.Writer
│   └── ui_test.go                # Scripted end-to-end UI workflow, EOF, and error tests
└── docs/
    ├── TEST_REPORT.md            # Phase III test report, 12-student verification, coverage results
    └── DOCUMENTATION.md          # Short documentation and captured console screens
```

## Documentation & Deliverables

- [Final Report (Phase III PDF)](2301730198_Divyanshi_ENSP465_Final_Report.pdf)
- [Project Presentation (PPTX)](2301730198_Divyanshi_ENSP465_Presentation.pptx)
- [Phase II Progress Report (PDF)](Phase_2_Progress_Report_Capstone_Project.pdf)
- [Phase I Proposal & System Design (PDF)](2301730198_Divyanshi_ENSP465_Phase1.pdf)
- [Test Report (Phase III)](docs/TEST_REPORT.md): Automated test results, 12-student hand-computed verification, coverage breakdown, and Phase III refinements.
- [Short Documentation](docs/DOCUMENTATION.md): Captured interactive console screens and worked examples.
