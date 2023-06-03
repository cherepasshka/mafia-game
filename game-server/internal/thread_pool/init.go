package threadpool

import (
	"sync"
)

type ThreadPool struct {
	queue chan func()
	wg    sync.WaitGroup
}

var threadPool *ThreadPool = nil

func New(maxWorkers int) *ThreadPool {
	pool := &ThreadPool{
		queue: make(chan func()),
	}
	pool.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer pool.wg.Done()
			for task := range pool.queue {
				task()
			}
		}()
	}
	if threadPool == nil {
		threadPool = pool
	}
	return pool
}

func GetThreadPool() *ThreadPool {
	return threadPool
}
