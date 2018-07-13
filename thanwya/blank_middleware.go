package thanwya

type blankMiddleware struct {
	middleware
}

func (bm blankMiddleware) next(students []Student) {

}

func (bm blankMiddleware) setNext(md middleware) {

}
