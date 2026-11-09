# Console-Based Student Management System in Go

Capstone Project 1 — ENSP465 New Age Programming Languages Lab
Divyanshi · 2301730198 · Semester VII B.Tech CSE · Section C

A menu-driven console application to add, view, search, update and delete student records. It calculates total, percentage and grade automatically and saves everything to a local `students.json`. Only the Go standard library is used.

## Run

```bash
go run .
go build -o studentmgmt . && ./studentmgmt -data mydata.json
```

## Project layout

| Path | Responsibility |
|------|----------------|
| `main.go` | Wires store → service → UI |
| `ui/` | Menus, input, formatted tables (no business logic) |
| `service/` | CRUD rules, duplicate detection, `ErrNotFound` / `ErrDuplicateID` |
| `store/` | `Store` interface + atomic JSON `FileStore` |
| `student/` | `student.go` (struct, constants), `validator.go` (input rules), `calc.go` (total, percentage, grade) |

See [docs/DOCUMENTATION.md](docs/DOCUMENTATION.md) for the submission documentation.
