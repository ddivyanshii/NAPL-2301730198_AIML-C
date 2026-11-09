package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"studentmgmt/service"
	"studentmgmt/student"
)

var errCancelled = errors.New("cancelled")

type UI struct {
	svc *service.Service
	in  *bufio.Scanner
	out io.Writer
}

func New(svc *service.Service, in io.Reader, out io.Writer) *UI {
	return &UI{svc: svc, in: bufio.NewScanner(in), out: out}
}

func (u *UI) Run() error {
	u.println("\nWelcome to the Student Management System")
	for {
		u.printMenu()
		choice, err := u.read("Enter choice (1-6): ")
		if err != nil {
			return u.endOf(err)
		}
		switch choice {
		case "1":
			err = u.addStudent()
		case "2":
			err = u.viewAll()
		case "3":
			err = u.search()
		case "4":
			err = u.updateStudent()
		case "5":
			err = u.deleteStudent()
		case "6":
			u.println("Goodbye! All changes are saved.")
			return nil
		default:
			u.println("  ✗ Invalid choice. Please enter a number from 1 to 6.")
			continue
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return u.endOf(err)
			}
			u.endOf(err)
		}
	}
}

func (u *UI) endOf(err error) error {
	switch {
	case errors.Is(err, io.EOF):
		u.println("\nInput closed. Goodbye!")
	case errors.Is(err, errCancelled):
		u.println("  Cancelled.")
	default:
		u.println("  ✗ Error:", err)
	}
	return nil
}

func (u *UI) printMenu() {
	u.println("\n========== STUDENT MANAGEMENT ==========")
	u.println(" 1. Add student")
	u.println(" 2. View all students")
	u.println(" 3. Search student")
	u.println(" 4. Update student")
	u.println(" 5. Delete student")
	u.println(" 6. Exit")
	u.println("========================================")
	u.println(" (type 'cancel' at any prompt to go back)")
}

func (u *UI) addStudent() error {
	u.println("\n--- Add Student ---")
	idStr, err := u.ask("Student ID", "", func(v string) error {
		id, err := parseID(v)
		if err != nil {
			return err
		}
		if _, err := u.svc.SearchByID(id); err == nil {
			return service.ErrDuplicateID
		}
		return nil
	})
	if err != nil {
		return err
	}
	id, _ := strconv.Atoi(idStr)

	st, err := u.askDetails(student.Student{ID: id}, false)
	if err != nil {
		return err
	}
	saved, err := u.svc.Add(st)
	if err != nil {
		return err
	}
	u.println("  ✓ Student added successfully.")
	u.printTable([]student.Student{saved})
	return nil
}

func (u *UI) viewAll() error {
	list := u.svc.List()
	if len(list) == 0 {
		u.println("\nNo student records yet.")
		return nil
	}
	u.printf("\n--- All Students (%d) ---\n", len(list))
	u.printTable(list)
	return nil
}

func (u *UI) search() error {
	u.println("\n--- Search ---")
	u.println(" 1. By ID")
	u.println(" 2. By name (partial match)")
	choice, err := u.read("Choice: ")
	if err != nil {
		return err
	}
	switch choice {
	case "1":
		v, err := u.ask("Student ID", "", func(v string) error { _, err := parseID(v); return err })
		if err != nil {
			return err
		}
		id, _ := strconv.Atoi(v)
		st, err := u.svc.SearchByID(id)
		if err != nil {
			return err
		}
		u.printTable([]student.Student{st})
	case "2":
		kw, err := u.ask("Name contains", "", func(v string) error {
			if strings.TrimSpace(v) == "" {
				return errors.New("enter some text to search for")
			}
			return nil
		})
		if err != nil {
			return err
		}
		res := u.svc.SearchByName(kw)
		if len(res) == 0 {
			u.printf("  No students match %q.\n", kw)
			return nil
		}
		u.printf("  %d match(es):\n", len(res))
		u.printTable(res)
	default:
		u.println("  ✗ Invalid choice.")
	}
	return nil
}

func (u *UI) updateStudent() error {
	u.println("\n--- Update Student ---")
	v, err := u.ask("Student ID to update", "", func(v string) error { _, err := parseID(v); return err })
	if err != nil {
		return err
	}
	id, _ := strconv.Atoi(v)
	cur, err := u.svc.SearchByID(id)
	if err != nil {
		return err
	}
	u.println("  Current record (press Enter to keep a value):")
	u.printTable([]student.Student{cur})

	st, err := u.askDetails(cur, true)
	if err != nil {
		return err
	}
	saved, err := u.svc.Update(id, st)
	if err != nil {
		return err
	}
	u.println("  ✓ Student updated successfully.")
	u.printTable([]student.Student{saved})
	return nil
}

func (u *UI) deleteStudent() error {
	u.println("\n--- Delete Student ---")
	v, err := u.ask("Student ID to delete", "", func(v string) error { _, err := parseID(v); return err })
	if err != nil {
		return err
	}
	id, _ := strconv.Atoi(v)
	st, err := u.svc.SearchByID(id)
	if err != nil {
		return err
	}
	u.printTable([]student.Student{st})
	ans, err := u.ask("Delete this student? (y/n)", "", func(v string) error {
		switch strings.ToLower(v) {
		case "y", "yes", "n", "no":
			return nil
		}
		return errors.New("please answer y or n")
	})
	if err != nil {
		return err
	}
	if strings.HasPrefix(strings.ToLower(ans), "n") {
		u.println("  Deletion aborted.")
		return nil
	}
	if err := u.svc.Delete(id); err != nil {
		return err
	}
	u.println("  ✓ Student deleted.")
	return nil
}

func (u *UI) askDetails(base student.Student, withDefaults bool) (student.Student, error) {
	def := func(s string) string {
		if withDefaults {
			return s
		}
		return ""
	}
	st := student.Student{ID: base.ID}
	var err error

	if st.Name, err = u.ask("Name", def(base.Name), student.ValidateName); err != nil {
		return st, err
	}
	ageStr, err := u.ask("Age", def(strconv.Itoa(base.Age)), func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("age must be a whole number")
		}
		return student.ValidateAge(n)
	})
	if err != nil {
		return st, err
	}
	st.Age, _ = strconv.Atoi(ageStr)

	if st.Course, err = u.ask("Course", def(base.Course), student.ValidateCourse); err != nil {
		return st, err
	}

	for i := 0; i < student.NumSubjects; i++ {
		d := ""
		if withDefaults && i < len(base.Marks) {
			d = formatMark(base.Marks[i])
		}
		label := fmt.Sprintf("Marks for Subject %d (0-%g)", i+1, student.MaxPerSubject)
		v, err := u.ask(label, d, func(v string) error {
			m, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return errors.New("marks must be a number")
			}
			return student.ValidateMark(m, student.MaxPerSubject)
		})
		if err != nil {
			return st, err
		}
		m, _ := strconv.ParseFloat(v, 64)
		st.Marks = append(st.Marks, m)
	}
	return st, nil
}

func (u *UI) ask(label, def string, check func(string) error) (string, error) {
	for {
		prompt := label + ": "
		if def != "" {
			prompt = fmt.Sprintf("%s [%s]: ", label, def)
		}
		v, err := u.read(prompt)
		if err != nil {
			return "", err
		}
		if v == "" && def != "" {
			v = def
		}
		if err := check(v); err != nil {
			u.println("  ✗", err)
			continue
		}
		return v, nil
	}
}

func (u *UI) read(prompt string) (string, error) {
	fmt.Fprint(u.out, prompt)
	if !u.in.Scan() {
		if err := u.in.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	v := strings.TrimSpace(u.in.Text())
	if strings.EqualFold(v, "cancel") {
		return "", errCancelled
	}
	return v, nil
}

func (u *UI) printTable(list []student.Student) {
	tw := tabwriter.NewWriter(u.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNAME\tAGE\tCOURSE\tMARKS\tTOTAL\tPERCENT\tGRADE")
	fmt.Fprintln(tw, "--\t----\t---\t------\t-----\t-----\t-------\t-----")
	for _, s := range list {
		parts := make([]string, len(s.Marks))
		for i, m := range s.Marks {
			parts[i] = formatMark(m)
		}
		fmt.Fprintf(tw, "%d\t%s\t%d\t%s\t%s\t%g/%g\t%.2f%%\t%s\n",
			s.ID, s.Name, s.Age, s.Course, strings.Join(parts, " "),
			s.Total, float64(len(s.Marks))*student.MaxPerSubject, s.Percentage, s.Grade)
	}
	tw.Flush()
}

func (u *UI) println(a ...any)          { fmt.Fprintln(u.out, a...) }
func (u *UI) printf(f string, a ...any) { fmt.Fprintf(u.out, f, a...) }
func formatMark(m float64) string       { return strconv.FormatFloat(m, 'f', -1, 64) }

func parseID(v string) (int, error) {
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, errors.New("ID must be a whole number")
	}
	return id, student.ValidateID(id)
}
