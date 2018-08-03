package thanwya

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	_ "github.com/lib/pq"
)

var (
	queryTemplate *template.Template
	wheel         string
	done          chan struct{}
	wheelIndex    int
)

const (
	insertQueryTemplate = `
INSERT INTO students 
(seat_number, name, total_grade, student_type, school, department_name, branch, number_of_failures)
VALUES
{{ range . }}
({{.SeatNumber}}, '{{.Name}}', '{{.TotalDegree}}', '{{.StudentType}}', '{{.SchoolName}}', '{{.DepartmentName}}', '{{.Section | branch}}', {{.NumberOfFailures}}) {{if not (isLastElement $ .)}},{{end}}
{{ end }}
`
)

func isLastElement(students []Student, student Student) bool {
	return student.SeatNumber == students[len(students)-1].SeatNumber
}

func init() {
	queryTemplate = template.Must(template.New("insertQuery").Funcs(template.FuncMap{"isLastElement": isLastElement, "branch": branch}).Parse(insertQueryTemplate))
	wheel = "/-\\|"
	done = make(chan struct{})
}

type savingMiddleware struct {
	middleware
	nextMiddleware middleware
	totalStudents  int
}

func (sm *savingMiddleware) next(students []Student) {
	defer func() {
		done <- struct{}{}
	}()
	go sm.wheelChanger()

	for i := 0; i < len(students); i += NumberOfCuncurrentInserts {
		query := bytes.Buffer{}
		end := min(i+NumberOfCuncurrentInserts, len(students))
		if err := queryTemplate.ExecuteTemplate(&query, "insertQuery", students[i:end]); err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec(query.String()); err != nil {
			log.Fatalf("%+v", err)
		}
		// fmt.Printf("\rStudents %d..%d inserted successfully...", i+1, end)
		sm.printProgress(end, len(students))
	}
	fmt.Println()
	sm.nextMiddleware.next(students)
}

func (sm *savingMiddleware) printProgress(current, total int) {
	fmt.Printf("\r")
	fmt.Printf("Saving ")
	sm.printWheel()
	fmt.Printf(" ")
	printProgress(current, total)
}

func (sm *savingMiddleware) setNext(md middleware) {
	sm.nextMiddleware = md
}

func (sm *savingMiddleware) wheelChanger() {
	exitLoop := false
	for {
		select {
		case <-done:
			exitLoop = true
		default:
		}
		if exitLoop {
			break
		}
		time.Sleep(RefreshRate)
		wheelIndex = (wheelIndex + 1) % len(wheel)
	}
}

func (sm *savingMiddleware) printWheel() {
	fmt.Printf(string(wheel[wheelIndex]))
}
