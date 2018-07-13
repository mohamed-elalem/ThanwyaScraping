package thanwya

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"text/template"

	_ "github.com/lib/pq"
)

var (
	db            *sql.DB
	queryTemplate *template.Template
)

const (
	insertQueryTemplate = `
INSERT INTO students 
(seat_number, name, total_grade)
VALUES
{{ range . }}
({{.SeatNumber}}, '{{.Name}}', '{{.TotalDegree}}') {{if not (isLastElement $ .)}},{{end}}
{{ end }}
`
)

func isLastElement(students []Student, student Student) bool {
	return student.SeatNumber == students[len(students)-1].SeatNumber
}

func init() {
	var err error
	connStr := "user=postgres dbname=thanwya sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	queryTemplate = template.Must(template.New("insertQuery").Funcs(template.FuncMap{"isLastElement": isLastElement}).Parse(insertQueryTemplate))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type savingMiddleware struct {
	middleware
	nextMiddleware middleware
}

func (sm *savingMiddleware) next(students []Student) {
	for i := 0; i < len(students); i += NumberOfCuncurrentInserts {
		query := bytes.Buffer{}
		end := min(i+NumberOfCuncurrentInserts, len(students))
		if err := queryTemplate.ExecuteTemplate(&query, "insertQuery", students[i:end]); err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec(query.String()); err != nil {
			log.Fatalf("%+v", err)
		}
		fmt.Printf("\rStudents %d..%d inserted successfully...", i+1, end)
	}
	fmt.Println()
	db.Close()
	sm.nextMiddleware.next(students)
}

func (sm *savingMiddleware) setNext(md middleware) {
	sm.nextMiddleware = md
}
