package util

import (
	. "github.com/orange-jacky/albums/common/jobqueue"
	"sync"
)

type jobQueue struct {
	*Worker
}

var (
	jobq      *jobQueue
	jobq_once sync.Once
)

//创建单实例
func JobQueue() *jobQueue {
	jobq_once.Do(func() {
		jobq = &jobQueue{NewWorker(100)}
	})
	return jobq
}

func GetJobQueue() *jobQueue {
	return jobq
}
