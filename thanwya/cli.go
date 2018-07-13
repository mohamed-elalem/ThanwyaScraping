package thanwya

import (
	"sync"
)

var limiter chan struct{}
var wg sync.WaitGroup
var mutex sync.RWMutex
var logger Log
var mw fetchingMiddleware

func init() {
	limiter = make(chan struct{}, NumberOfGoRoutines)
	logger.Init()
	mw = fetchingMiddleware{nextMiddleware: &savingMiddleware{nextMiddleware: &blankMiddleware{}}}
}

func Run() {
	// resp, err := http.PostForm("http://natega.youm7.com/Home/GetResultStage1/", url.Values{"seatNumber": []string{"123456"}})
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	// defer resp.Body.Close()

	// if status := resp.StatusCode; status != http.StatusOK {
	// 	log.Printf("The request received status code %d\n", status)
	// }
	// if err := json.NewDecoder(resp.Body).Decode(&student); err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Printf("%+v\n", student)
	defer logger.Destroy()
	mw.next(make([]Student, 0))
}
