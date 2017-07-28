package jobqueue

import (
//"log"
)

type Worker struct {
	MaxJobs    int      //任务队列长度
	JobChannel chan Job //任务队列
	Quit       chan bool
}

func NewWorker(maxJobs int) *Worker {
	//log.Printf("[worker] create a worker, queue len is %d\n", maxJobs)
	return &Worker{
		MaxJobs:    maxJobs,
		JobChannel: make(chan Job, maxJobs),
		Quit:       make(chan bool, 1),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.JobChannel:
				job.Handler(job.Input, job.Output)
			case <-w.Quit:
				//log.Println("[worker] a worker stoped.")
				return
			}
		}
	}()
}

func (w *Worker) Push(job Job) {
	w.JobChannel <- job
}

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
