package gc

import (
	"sync"
	"time"

	"github.com/concourse/atc/worker"
)

type WorkerPool struct {
	workerPool       worker.Client
	maxJobsPerWorker int

	workers  map[string]worker.Worker
	workersL *sync.Mutex

	workerJobs  map[string]int
	workerJobsL *sync.Mutex
}

type Job interface {
	Run(worker.Worker)
}

type JobFunc func(worker.Worker)

func (f JobFunc) Run(workerClient worker.Worker) {
	f(workerClient)
}

func NewWorkerPool(workerPool worker.Client, maxJobsPerWorker int) *WorkerPool {
	pool := &WorkerPool{
		workerPool:       workerPool,
		maxJobsPerWorker: maxJobsPerWorker,
	}

	go pool.syncWorkersLoop()

	return pool
}

func (pool *WorkerPool) Queue(workerName string, job Job) {
	if !pool.startJob(workerName) {
		// drop the job on the floor; it'll be queued up again later
		return
	}

	pool.workersL.Lock()
	workerClient, found := pool.workers[workerName]
	pool.workersL.Unlock()

	if !found {
		// drop the job on the floor; it'll be queued up again later
		return
	}

	go func() {
		defer pool.finishJob(workerName)
		job.Run(workerClient)
	}()
}

func (pool *WorkerPool) startJob(workerName string) bool {
	pool.workerJobsL.Lock()
	defer pool.workerJobsL.Unlock()

	if pool.workerJobs[workerName] == pool.maxJobsPerWorker {
		return false
	}

	pool.workerJobs[workerName]++

	return true
}

func (pool *WorkerPool) finishJob(workerName string) {
	pool.workerJobsL.Lock()
	pool.workerJobs[workerName]--
	pool.workerJobsL.Unlock()
}

func (pool *WorkerPool) syncWorkersLoop() {
	pool.syncWorkers()

	ticker := time.NewTicker(30 * time.Second) // XXX: parameterize same as default worker TTL (...which might actually live on the worker side...)

	for {
		select {
		case <-ticker.C:
			pool.syncWorkers()
		}
	}
}

func (pool *WorkerPool) syncWorkers() {
	workers, err := pool.workerPool.RunningWorkers(nil) // XXX: logger
	if err != nil {
		// XXX: log
		return
	}

	workerMap := map[string]worker.Worker{}
	for _, worker := range workers {
		workerMap[worker.Name()] = worker
	}

	pool.workersL.Lock()
	pool.workers = workerMap
	pool.workersL.Unlock()
}
