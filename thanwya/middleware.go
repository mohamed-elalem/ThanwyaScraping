package thanwya

type middleware interface {
	next([]Student)
	setNext(middleware)
}
