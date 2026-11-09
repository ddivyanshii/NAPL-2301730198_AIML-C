# Short Documentation — Console-Based Student Management System in Go

**ENSP465 — New Age Programming Languages Lab · Capstone Project**
Divyanshi · 2301730198 · Semester VII (B.Tech CSE) · Section C · Faculty: Dr. Hemraj Kumawat

## 1. Project Overview

Small institutes still keep student data in paper registers or spreadsheets, where totals and grades are computed by hand, duplicates creep in, and searching is slow. This project is a lightweight, offline **console application written in Go** that manages student records through a simple menu and keeps them saved between runs. It needs no server, database or network.

## 2. Objectives

1. Provide all five CRUD operations (add, view, search, update, delete) through a menu.
2. Calculate total marks, percentage and grade automatically for every student.
3. Persist records in a JSON file and restore them on restart.
4. Validate all input so invalid data never corrupts the stored records.

## 3. Key Features

- **Add** a student (ID, name, age, course, marks in 5 subjects).
- **View** all students in an aligned table, sorted by ID.
- **Search** by ID or by partial, case-insensitive name.
- **Update** any field (press Enter to keep the current value); **delete** with y/n confirmation.
- **Automatic results** — total, percentage and grade (O ≥ 90, A ≥ 80, B ≥ 70, C ≥ 60, D ≥ 50, else F).
- **JSON persistence** — loaded at start-up, saved after every change using a temp file + rename, so a crash cannot leave a half-written file. If a save fails, the in-memory data is left unchanged.
- **Validation with friendly messages** — unique positive ID, non-empty name/course, age 5–100, marks 0–100 (NaN/Inf rejected). Invalid answers are re-asked; typing `cancel` returns to the menu.

## 4. Technologies Used

- **Go 1.21+**, standard library only: `fmt`, `bufio`, `os`, `encoding/json`, `strconv`, `strings`, `sort`, `errors`, `text/tabwriter`.
- Go features demonstrated: structs with JSON tags, slices and maps, packages, interfaces, custom/wrapped errors (`%w`, `errors.Is`), `defer`.

## 5. Implementation

Layered architecture — each package has a single responsibility:

```
main.go ──► ui ──► service ──► store (interface) ──► students.json
                      │
                      └──► student  (struct + validator.go + calc.go)
```

| Package | Role |
|---------|------|
| `ui` | Menu loop, prompts, table output. Uses injected `io.Reader`/`io.Writer` streams. |
| `service` | Holds the in-memory `[]Student` (ordered by ID) plus a `map[int]int` index for O(1) lookup and duplicate detection. Each change is saved first and only then applied in memory. |
| `student` | The `Student` struct plus two files: `validator.go` (field-level and whole-record rules returning sentinel errors) and `calc.go` (`Total`, `Percentage`, `Grade`, `Compute`; grading bands live in one table). |
| `store` | `Store` interface (`Load`, `Save`) and the atomic `FileStore`. |

**Worked example** — ID 101, Aman Verma, marks 78, 85, 90, 72, 88 → total 413/500, 82.60 %, grade **A**.

## 6. User Interface (console)

The UI is a text menu in the terminal. A real session, captured from the program:

```
========== STUDENT MANAGEMENT ==========
 1. Add student
 2. View all students
 3. Search student
 4. Update student
 5. Delete student
 6. Exit
========================================
 (type 'cancel' at any prompt to go back)
Enter choice (1-6): 2

--- All Students (2) ---
ID   NAME         AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----         ---  ------      -----           -----    -------  -----
101  Aman Verma   21   B.Tech CSE  78 85 90 72 88  413/500  82.60%   A
102  Riya Sharma  20   BCA         95 92 98 90 96  471/500  94.20%   O
```

> **Screenshots:** add terminal screenshots of each operation (add, view, search, update, delete, an invalid-input message) here before submitting.

## 7. How to Run

```bash
go run .
```

