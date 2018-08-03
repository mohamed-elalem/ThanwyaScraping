package thanwya

import (
	"database/sql"
	"log"
)

var (
	db *sql.DB
)

func init() {
	var err error
	connStr := "user=" + DatabaseUser + " dbname=" + DatabaseName + " sslmode=" + DatabaseSSLMode
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func closeDB() {
	db.Close()
}
