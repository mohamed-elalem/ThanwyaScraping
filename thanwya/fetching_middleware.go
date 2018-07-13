package thanwya

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type fetchingMiddleware struct {
	middleware
	exist          int
	notExist       int
	nextMiddleware middleware
	students       []Student
}

func (fm *fetchingMiddleware) setNext(md middleware) {
	fm.nextMiddleware = md
}

func (fm *fetchingMiddleware) next(students []Student) {
	fm.students = students
	for seatNumber := InitialSeatNumber; seatNumber <= LastSeatNumber; seatNumber++ {
		go func(seatNumber int) {
			wg.Add(1)
			fm.performFetching(seatNumber)
			wg.Done()
		}(seatNumber)
	}
	wg.Wait()
	fmt.Println()
	fm.nextMiddleware.next(fm.students)
}

func (fm *fetchingMiddleware) performFetching(seatNumber int) {
	limiter <- struct{}{}
	defer func() {
		<-limiter
	}()
	student, err := fm.fetchStudent(seatNumber)
	if err != nil {
		logger.Log(err.Error(), WARNING)
		fm.notExist++
	} else {
		fm.appendStudent(student)
	}
	fm.printProgress()
}

func (fm *fetchingMiddleware) printProgress() {
	fmt.Printf("\r%d students are fetched, failed attempts %d", fm.exist, fm.notExist)
}

func (fm *fetchingMiddleware) fetchStudent(seatNumber int) ([]Student, error) {
	var student []Student

	resp, err := http.PostForm(BaseURL, url.Values{"seatNumber": {strconv.Itoa(seatNumber)}})

	if err != nil {
		logger.Log(err.Error(), ERROR)
		return student, err
	}

	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		logger.Log("Request responded with status "+strconv.Itoa(status), WARNING)
	}

	if err := json.NewDecoder(resp.Body).Decode(&student); err != nil {
		logger.Log(err.Error(), ERROR)
	}

	if len(student) == 0 {
		return student, errors.New("Student with seat number " + strconv.Itoa(seatNumber) + " does not exist.")
	}
	return student, nil
}

func (fm *fetchingMiddleware) appendStudent(student []Student) {
	// For now it's not required to synchronize since i will
	// First fetch then insert
	mutex.Lock()
	defer mutex.Unlock()
	if len(student) > 0 {
		fm.students = append(fm.students, student[0])
		fm.exist++
	}
}
