# Phase III Test Report: Console-Based Student Management System in Go

**ENSP465: New Age Programming Languages Lab · Capstone Project (Phase III)**
Divyanshi · Roll No. `2301730198` · Semester VII (B.Tech CSE) · Section C · Faculty: Dr. Hemraj Kumawat

---

## 1. Test Objectives & Strategy

Phase I and Phase II established three measurable verification targets for Phase III:

1. **Automated Unit Tests (≥ 70% coverage)** covering the calculator (`student/calc.go`), validator (`student/validator.go`), and service layer (`service/service.go`) using table-driven tests and an in-memory mock `Store`.
2. **Calculation Accuracy (100% agreement)** between program output and hand-computed results across a test set of **at least 10 sample students** spanning all grade bands and boundary values.
3. **Persistence Integrity (100% restored)** verifying that every saved record is restored identically after restarting the application, including edge cases for missing, empty (0-byte / whitespace / `[]` / `null`), and corrupt data files.

In addition, automated scripted UI integration tests (`ui/ui_test.go`) were implemented using Go's injected `io.Reader` / `io.Writer` streams, and static analysis (`gofmt`, `go vet`, and `-race`) was run across the entire module.

---

## 2. Test Suite Overview

| Test File | Package | Focus & Techniques |
|-----------|---------|--------------------|
| `student/calc_test.go` | `student` | Table-driven tests for `Total`, `Percentage`, `Grade` (all band boundaries + `NaN` / negative), and 12-student `Compute` verification. |
| `student/validator_test.go` | `student` | Table-driven tests for `ValidateID`, `ValidateName`, `ValidateCourse`, `ValidateAge`, `ValidateMark` (including `NaN`, `+Inf`, `-Inf`), `ValidateMarks` (wrong counts, subject-index errors), and `Validate`. |
| `service/service_test.go` | `service` | Unit tests with an in-memory `mockStore` (injectable load/save failures): `New` (sorting, duplicates, invalid stored records), `Add`, `List`, `SearchByID`, `SearchByName`, `Update`, `Delete` (index rebuild), defensive-copy isolation, and save-first rollback guarantees. |
| `store/store_test.go` | `store_test` | Real filesystem tests using `t.TempDir()`: missing file auto-creation, empty 0-byte/whitespace/`[]`/`null` files, corrupt JSON (`ErrCorrupt`), I/O errors, 12-student save & simulated restart round-trip, and atomic temp-file cleanup check. |
| `ui/ui_test.go` | `ui_test` | Scripted end-to-end console session covering all 6 menu options, invalid choices, `cancel` at main menu and sub-prompts, re-prompting on bad inputs, Enter-to-keep defaults on update, `y/n` delete confirmation, EOF handling, and reader error propagation. |

---

## 3. Calculation Verification (12 Sample Students)

Every record below was computed by hand and compared against `student.Compute(marks, 100.0)` in `TestCompute_SampleDataset` (`student/calc_test.go`):

| ID | Student Name | Course | Marks (5 Subjects) | Hand Total | Program Total | Hand % | Program % | Hand Grade | Program Grade | Status |
|----|--------------|--------|--------------------|------------|---------------|--------|-----------|------------|---------------|--------|
| 101 | Aman Verma | B.Tech CSE | 78, 85, 90, 72, 88 | 413 | 413 | 82.60% | 82.60% | A | A | Pass |
| 102 | Riya Sharma | BCA | 95, 92, 98, 90, 96 | 471 | 471 | 94.20% | 94.20% | O | O | Pass |
| 103 | Aman Gupta | B.Sc | 55, 48, 62, 51, 58 | 274 | 274 | 54.80% | 54.80% | D | D | Pass |
| 104 | Neha Kapoor | B.Tech AI | 90, 90, 90, 90, 90 | 450 | 450 | 90.00% | 90.00% | O | O | Pass |
| 105 | Kabir Mehta | B.Tech CSE | 80, 80, 80, 80, 80 | 400 | 400 | 80.00% | 80.00% | A | A | Pass |
| 106 | Priya Nair | MCA | 72, 75, 68, 74, 71 | 360 | 360 | 72.00% | 72.00% | B | B | Pass |
| 107 | Rohan Joshi | B.Tech ECE | 70, 70, 70, 70, 70 | 350 | 350 | 70.00% | 70.00% | B | B | Pass |
| 108 | Simran Kaur | BBA | 60, 65, 62, 68, 64 | 319 | 319 | 63.80% | 63.80% | C | C | Pass |
| 109 | Dev Patel | B.Sc CS | 60, 60, 60, 60, 60 | 300 | 300 | 60.00% | 60.00% | C | C | Pass |
| 110 | Ananya Iyer | BA | 50, 50, 50, 50, 50 | 250 | 250 | 50.00% | 50.00% | D | D | Pass |
| 111 | Vikram Singh | B.Com | 49.5, 50, 49.5, 50, 50 | 249 | 249 | 49.80% | 49.80% | F | F | Pass |
| 112 | Meera Das | B.Tech CSE | 100, 99.5, 98.5, 100, 97.5 | 495.5 | 495.5 | 99.10% | 99.10% | O | O | Pass |

**Result:** 12 / 12 matched (**100% agreement**).

---

## 4. Automated Test Execution & Coverage Results

Command executed inside `capstone1/`:

```bash
go test -v -race -coverprofile=/tmp/cover.out ./... && go tool cover -func=/tmp/cover.out
```

### Package Coverage Summary

| Package | Target | Achieved | Status |
|---------|--------|----------|--------|
| `studentmgmt/student` (`calc.go`, `validator.go`) | ≥ 70.0% | **100.0%** | Exceeded |
| `studentmgmt/service` (`service.go`) | ≥ 70.0% | **100.0%** | Exceeded |
| `studentmgmt/ui` (`ui.go`) | N/A | **93.8%** | Pass |
| `studentmgmt/store` (`store.go`) | N/A | **78.9%** | Pass |
| **Total (all packages incl. `main.go`)** | ≥ 70.0% | **91.2%** | Exceeded |

### Captured `go tool cover -func` Output

```text
studentmgmt/main.go:13:                 main                    0.0%
studentmgmt/service/service.go:28:      New                     100.0%
studentmgmt/service/service.go:48:      Add                     100.0%
studentmgmt/service/service.go:63:      List                    100.0%
studentmgmt/service/service.go:68:      SearchByID              100.0%
studentmgmt/service/service.go:77:      SearchByName            100.0%
studentmgmt/service/service.go:92:      Update                  100.0%
studentmgmt/service/service.go:111:     Delete                  100.0%
studentmgmt/service/service.go:121:     commit                  100.0%
studentmgmt/service/service.go:129:     replace                 100.0%
studentmgmt/service/service.go:134:     applySorted             100.0%
studentmgmt/service/service.go:146:     snapshot                100.0%
studentmgmt/service/service.go:154:     withResults             100.0%
studentmgmt/service/service.go:162:     clone                   100.0%
studentmgmt/store/store.go:30:          NewFileStore            100.0%
studentmgmt/store/store.go:36:          Load                    93.3%
studentmgmt/store/store.go:61:          Save                    68.2%
studentmgmt/student/calc.go:22:         Total                   100.0%
studentmgmt/student/calc.go:31:         Percentage              100.0%
studentmgmt/student/calc.go:39:         Grade                   100.0%
studentmgmt/student/calc.go:49:         Compute                 100.0%
studentmgmt/student/calc.go:55:         round2                  100.0%
studentmgmt/student/validator.go:28:    ValidateID              100.0%
studentmgmt/student/validator.go:36:    ValidateName            100.0%
studentmgmt/student/validator.go:44:    ValidateCourse          100.0%
studentmgmt/student/validator.go:52:    ValidateAge             100.0%
studentmgmt/student/validator.go:60:    ValidateMark            100.0%
studentmgmt/student/validator.go:68:    ValidateMarks           100.0%
studentmgmt/student/validator.go:81:    Validate                100.0%
studentmgmt/ui/ui.go:27:                New                     100.0%
studentmgmt/ui/ui.go:32:                Run                     100.0%
studentmgmt/ui/ui.go:71:                handleStreamError       100.0%
studentmgmt/ui/ui.go:79:                reportOpError           100.0%
studentmgmt/ui/ui.go:87:                printMenu               100.0%
studentmgmt/ui/ui.go:99:                addStudent              95.7%
studentmgmt/ui/ui.go:129:               viewAll                 100.0%
studentmgmt/ui/ui.go:140:               search                  89.3%
studentmgmt/ui/ui.go:182:               updateStudent           86.4%
studentmgmt/ui/ui.go:208:               deleteStudent           85.7%
studentmgmt/ui/ui.go:240:               askID                   85.7%
studentmgmt/ui/ui.go:252:               askDetails              92.1%
studentmgmt/ui/ui.go:303:               ask                     100.0%
studentmgmt/ui/ui.go:324:               read                    100.0%
studentmgmt/ui/ui.go:339:               printTable              100.0%
studentmgmt/ui/ui.go:355:               println                 100.0%
studentmgmt/ui/ui.go:356:               printf                  100.0%
studentmgmt/ui/ui.go:357:               formatMark              100.0%
studentmgmt/ui/ui.go:359:               parseID                 100.0%
total:                                  (statements)            91.2%
```

All tests passed with `-race` enabled (`0` failures, `0` data races).

---

## 5. Persistence Verification

| Scenario | Test Function | Expected Behaviour | Observed Result |
|----------|---------------|--------------------|-----------------|
| Missing data file | `TestFileStore_MissingFileCreated` | Creates empty `students.json`, returns empty slice | Pass |
| 0-byte empty file | `TestFileStore_EmptyAndNullFiles/zero-byte_file` | Returns non-nil empty slice without error | Pass |
| Whitespace-only file | `TestFileStore_EmptyAndNullFiles/whitespace-only_file` | Returns non-nil empty slice without error | Pass |
| Empty JSON `[]` and `null` | `TestFileStore_EmptyAndNullFiles` | Returns non-nil empty slice | Pass |
| Corrupt JSON syntax | `TestFileStore_CorruptJSON` | Returns wrapped `store.ErrCorrupt` | Pass |
| Unwritable / missing directory | `TestFileStore_IOErrors` | Returns descriptive error without panicking | Pass |
| 12-student save & restart | `TestFileStore_PersistenceAcrossRestart` | 12/12 records restored identically; zero `.tmp*` files left | Pass |

---

## 6. Milestone 5 Bug Fixes, Edge Cases & Clean-up

During Phase III testing and static analysis, the following issues and refinements were addressed:

1. **Empty file handling (`store/store.go`):** A 0-byte or whitespace-only `students.json` previously triggered `ErrCorrupt` from `json.Unmarshal`. `Load()` now checks `len(bytes.TrimSpace(data)) == 0` and returns an empty slice cleanly.
2. **`cancel` at main menu (`ui/ui.go`):** Typing `cancel` at `Enter choice (1-6):` previously returned `errCancelled` and exited the application. It now prints `Already at the main menu.` and stays in the menu loop.
3. **Reader I/O error handling (`ui/ui.go`):** Real scanner errors (`u.in.Err()`) during menu or sub-prompt reads are now propagated back to `main.go` instead of being swallowed or looping.
4. **Startup validation & defensive copies (`service/service.go`):** `New()` now validates every record loaded from storage (`student.Validate`) and clones `Marks` slices; `Add` and `Update` return defensive clones; redundant `sort.Slice` inside `commit` → `replace` was eliminated.
5. **Removed duplicated code (`ui/ui.go`):** Extracted `askID(label string)` to eliminate duplicated ID prompt + parse closures across `search`, `updateStudent`, and `deleteStudent`.
6. **Formatting, static analysis & GoDoc:** Ran `gofmt -w .` (fixing formatting in `student/validator.go`), verified `go vet ./...` reports zero issues, and added GoDoc comments on all exported types, constants, errors, and functions.

---

## 7. Final Milestone & Objective Status

| Milestone | Deliverables | Status |
|-----------|--------------|--------|
| 1. Phase I Submission | Proposal and system design document | Completed |
| 2. Setup & Core Model | Go module, `Student` struct, calculator, validator | Completed |
| 3. CRUD & Persistence | Service layer, atomic JSON `FileStore` | Completed |
| 4. Phase II Console UI & Demo | Menu-driven UI, intermediate report & demo | Completed |
| 5. Testing & Refinement | Automated tests (91.2% total / 100% core), edge-case fixes, `gofmt`/`go vet` | **Completed** |
| 6. Phase III Final Deliverable | Final source code, README, test report, presentation & final report | **Completed** |
