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

The application has a menu-driven console interface. Every screen below is real output captured from the program; the text after each prompt is what the user typed. Typing `cancel` at any prompt returns to the main menu.

### 6.1 Main menu

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
Enter choice (1-6): 
```

### 6.2 Add student

Invalid answers are rejected with a message and the same prompt is shown again.

```
Enter choice (1-6): 1

--- Add Student ---
Student ID: abc
  ✗ ID must be a whole number
Student ID: 101
Name: Aman Verma
Age: 21
Course: B.Tech CSE
Marks for Subject 1 (0-100): 78
Marks for Subject 2 (0-100): 85
Marks for Subject 3 (0-100): 90
Marks for Subject 4 (0-100): 72
Marks for Subject 5 (0-100): 88
  ✓ Student added successfully.
ID   NAME        AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----        ---  ------      -----           -----    -------  -----
101  Aman Verma  21   B.Tech CSE  78 85 90 72 88  413/500  82.60%   A
```

### 6.3 Input validation (age and marks out of range)

```
Enter choice (1-6): 1

--- Add Student ---
Student ID: 102
Name: Riya Sharma
Age: 3
  ✗ age must be between 5 and 100
Age: 20
Course: BCA
Marks for Subject 1 (0-100): 95
Marks for Subject 2 (0-100): 92
Marks for Subject 3 (0-100): 105
  ✗ invalid marks: marks must be between 0 and 100
Marks for Subject 3 (0-100): 98
Marks for Subject 4 (0-100): 90
Marks for Subject 5 (0-100): 96
  ✓ Student added successfully.
ID   NAME         AGE  COURSE  MARKS           TOTAL    PERCENT  GRADE
--   ----         ---  ------  -----           -----    -------  -----
102  Riya Sharma  20   BCA     95 92 98 90 96  471/500  94.20%   O
```

### 6.4 Duplicate ID and cancel

```
Enter choice (1-6): 1

--- Add Student ---
Student ID: 101
  ✗ a student with this ID already exists
Student ID: cancel
  Cancelled.
```

### 6.5 View all students

Records are shown in a table sorted by ID.

```
Enter choice (1-6): 2

--- All Students (3) ---
ID   NAME         AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----         ---  ------      -----           -----    -------  -----
101  Aman Verma   21   B.Tech CSE  78 85 90 72 88  413/500  82.60%   A
102  Riya Sharma  20   BCA         95 92 98 90 96  471/500  94.20%   O
103  Aman Gupta   22   B.Sc        55 48 62 51 58  274/500  54.80%   D
```

### 6.6 Search by ID and by name

```
Enter choice (1-6): 3

--- Search ---
 1. By ID
 2. By name (partial match)
Choice: 1
Student ID: 101
ID   NAME        AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----        ---  ------      -----           -----    -------  -----
101  Aman Verma  21   B.Tech CSE  78 85 90 72 88  413/500  82.60%   A
Enter choice (1-6): 3

--- Search ---
 1. By ID
 2. By name (partial match)
Choice: 2
Name contains: aman
  2 match(es):
ID   NAME        AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----        ---  ------      -----           -----    -------  -----
101  Aman Verma  21   B.Tech CSE  78 85 90 72 88  413/500  82.60%   A
103  Aman Gupta  22   B.Sc        55 48 62 51 58  274/500  54.80%   D
```

### 6.7 Search with no result or unknown ID

```
Enter choice (1-6): 3

--- Search ---
 1. By ID
 2. By name (partial match)
Choice: 2
Name contains: zzz
  No students match "zzz".
Enter choice (1-6): 3

--- Search ---
 1. By ID
 2. By name (partial match)
Choice: 1
Student ID: 999
  ✗ Error: ID 999: student not found
```

### 6.8 Update student

The current value is shown in brackets; pressing Enter keeps it. Total, percentage and grade are recalculated automatically (here Subject 5 changes from 88 to 95).

```
Enter choice (1-6): 4

--- Update Student ---
Student ID to update: 101
  Current record (press Enter to keep a value):
ID   NAME        AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----        ---  ------      -----           -----    -------  -----
101  Aman Verma  21   B.Tech CSE  78 85 90 72 88  413/500  82.60%   A
Name [Aman Verma]: 
Age [21]: 
Course [B.Tech CSE]: 
Marks for Subject 1 (0-100) [78]: 
Marks for Subject 2 (0-100) [85]: 
Marks for Subject 3 (0-100) [90]: 
Marks for Subject 4 (0-100) [72]: 
Marks for Subject 5 (0-100) [88]: 95
  ✓ Student updated successfully.
ID   NAME        AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----        ---  ------      -----           -----    -------  -----
101  Aman Verma  21   B.Tech CSE  78 85 90 72 95  420/500  84.00%   A
```

### 6.9 Delete student (with confirmation)

Answering `n` aborts the deletion; answering `y` removes the record.

```
Enter choice (1-6): 5

--- Delete Student ---
Student ID to delete: 102
ID   NAME         AGE  COURSE  MARKS           TOTAL    PERCENT  GRADE
--   ----         ---  ------  -----           -----    -------  -----
102  Riya Sharma  20   BCA     95 92 98 90 96  471/500  94.20%   O
Delete this student? (y/n): n
  Deletion aborted.
Enter choice (1-6): 5

--- Delete Student ---
Student ID to delete: 102
ID   NAME         AGE  COURSE  MARKS           TOTAL    PERCENT  GRADE
--   ----         ---  ------  -----           -----    -------  -----
102  Riya Sharma  20   BCA     95 92 98 90 96  471/500  94.20%   O
Delete this student? (y/n): y
  ✓ Student deleted.
Enter choice (1-6): 2

--- All Students (2) ---
ID   NAME        AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----        ---  ------      -----           -----    -------  -----
101  Aman Verma  21   B.Tech CSE  78 85 90 72 95  420/500  84.00%   A
103  Aman Gupta  22   B.Sc        55 48 62 51 58  274/500  54.80%   D
```

### 6.10 Invalid menu choice and unknown ID on update

```
Enter choice (1-6): 9
  ✗ Invalid choice. Please enter a number from 1 to 6.
Enter choice (1-6): 4

--- Update Student ---
Student ID to update: 555
  ✗ Error: ID 555: student not found
```

### 6.11 Exit and persistence

Choosing 6 exits the program. After restarting, the saved records are loaded again from `students.json`.

```
Enter choice (1-6): 6
Goodbye! All changes are saved.
```

After restart, option 2 shows the same records:

```
Enter choice (1-6): 2

--- All Students (2) ---
ID   NAME        AGE  COURSE      MARKS           TOTAL    PERCENT  GRADE
--   ----        ---  ------      -----           -----    -------  -----
101  Aman Verma  21   B.Tech CSE  78 85 90 72 95  420/500  84.00%   A
103  Aman Gupta  22   B.Sc        55 48 62 51 58  274/500  54.80%   D
```

Contents of the saved `students.json`:

```json
[
  {
    "id": 101,
    "name": "Aman Verma",
    "age": 21,
    "course": "B.Tech CSE",
    "marks": [78, 85, 90, 72, 95],
    "total": 420,
    "percentage": 84,
    "grade": "A"
  },
  {
    "id": 103,
    "name": "Aman Gupta",
    "age": 22,
    "course": "B.Sc",
    "marks": [55, 48, 62, 51, 58],
    "total": 274,
    "percentage": 54.8,
    "grade": "D"
  }
]
```

(The file on disk lists each mark on its own line because of Go's indented JSON output; it is shown compactly here.)

## 7. How to Run

```bash
go run .
```

