package thanwya

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
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
		limiter <- struct{}{}
		wg.Add(1)
		go func(seatNumber int) {
			fm.performFetching(seatNumber)
			fm.performNext(false)
			wg.Done()
		}(seatNumber)
	}
	wg.Wait()
	fmt.Println()
	if len(fm.students) > 0 {
		fm.performNext(true)
	}
}

func (fm *fetchingMiddleware) performNext(force bool) {
	if force || len(fm.students) >= MaxNumberOfArraySizeBeforeSave {
		mutex.Lock()
		if len(fm.students) >= MaxNumberOfArraySizeBeforeSave {
			fmt.Println()
			fm.nextMiddleware.next(fm.students)
		}
		fm.students = fm.students[:0]
		mutex.Unlock()
	}
}

func (fm *fetchingMiddleware) performFetching(seatNumber int) {
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
	fmt.Printf("\r%d students are fetched, failed attempts %d, total %d, Working goroutines %d", fm.exist, fm.notExist, fm.exist+fm.notExist, runtime.NumGoroutine())
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
