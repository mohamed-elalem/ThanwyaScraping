package thanwya

import (
	logging "log"
	"os"
	"time"
)

const (
	INFO    = 0
	WARNING = 1
	ERROR   = 2
)

var logLevels []string

var timestamp string

type Log struct {
	file *os.File
}

func (log *Log) Init() {
	logLevels = make([]string, 3)
	logLevels = append(logLevels, "INFO")
	logLevels = append(logLevels, "WARNING")
	logLevels = append(logLevels, "ERROR")
	timestamp = time.Now().Format(time.RFC850)

	fileName := "run-" + timestamp + ".log"

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logging.Fatal(err)
	}
	log.file = file
}

func (log *Log) Log(message string, level int) {
	if _, err := log.file.WriteString(logLevels[level] + " - " + time.Now().Format(time.RFC1123Z) + message + "\n"); err != nil {
		logging.Fatal(err)
	}
}

func (log *Log) Destroy() {
	log.file.Close()
}
