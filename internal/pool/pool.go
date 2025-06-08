package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const defaultEdge = 3

var (
	ErrNoWorkers = errors.New("there are no more workers")
	ErrNotExist  = errors.New("worker not exist in map")
)

type Worker struct {
	id     uint64
	cancel context.CancelFunc
}

type WorkerPool struct {
	mu      sync.Mutex
	jobs    chan string
	workers map[uint64]*Worker
	counter uint64
	wg      sync.WaitGroup
}

func NewWorkerPool(workers uint64) *WorkerPool {
	pool := &WorkerPool{
		jobs:    make(chan string, defaultEdge),
		workers: make(map[uint64]*Worker),
	}

	for i := uint64(0); i < workers; i++ {
		pool.AddWorker()
	}

	return pool
}

func (p *WorkerPool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.counter++
	id := p.counter
	ctx, cancel := context.WithCancel(context.Background())

	worker := &Worker{
		id:     id,
		cancel: cancel,
	}

	p.workers[id] = worker
	p.wg.Add(1)
	go p.runWorker(ctx, id)
}

func (p *WorkerPool) RemoveWorker() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.counter == 0 {
		return ErrNoWorkers
	}

	worker, ok := p.workers[p.counter]
	if !ok {
		return fmt.Errorf("worker %d: %w", p.counter, ErrNotExist)
	}

	worker.cancel()
	delete(p.workers, p.counter)
	p.counter--
	return nil
}

func (p *WorkerPool) Process(data string) {
	p.jobs <- data
}

func (p *WorkerPool) Close() {
	close(p.jobs)
	p.mu.Lock()

	for _, worker := range p.workers {
		worker.cancel()
	}

	p.mu.Unlock()
	p.wg.Wait()
}

func (p *WorkerPool) runWorker(ctx context.Context, id uint64) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d is dead\n", id)
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}

			time.Sleep(59 * time.Millisecond)
			fmt.Printf("worker %d processing: %s\n", id, job)
		}
	}
}
