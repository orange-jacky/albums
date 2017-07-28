package jobqueue

type HandlerFunc func(req interface{}, resp interface{})

type Job struct {
	Input   interface{}
	Output  interface{}
	Handler HandlerFunc
}
