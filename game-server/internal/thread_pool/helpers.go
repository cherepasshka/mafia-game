package threadpool

func (pool *ThreadPool) AddTask(task func()) {
	pool.queue <- task
}

func (pool *ThreadPool) Wait() {
	close(pool.queue)
	pool.wg.Wait()
}
