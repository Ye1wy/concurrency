package pool

import (
	"context"
	"fmt"
	"time"
)

const defaultEdge = 3

type WorkerPool struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	workers   uint64
	jobs      chan string
}

func NewWorkerPool(workers uint64) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		ctx:       ctx,
		ctxCancel: cancel,
		workers:   workers,
		jobs:      make(chan string, defaultEdge),
	}

	pool.initWorkers()

	return pool
}

func (p *WorkerPool) Add() {
	p.workers++
	go p.worker(p.workers, p.jobs)
}

func (p *WorkerPool) Delete() {
	p.ctxCancel()
	p.workers--
}

func (p *WorkerPool) initWorkers() {
	for id := 1; uint64(id) <= p.workers; id++ {
		go p.worker(uint64(id), p.jobs)
	}
}

func (p *WorkerPool) ProcessData(data string) {
	p.jobs <- data
}

func (p *WorkerPool) CloseJobChan() {
	close(p.jobs)
}

func (p *WorkerPool) worker(id uint64, jobs chan string) {
	for j := range jobs {
		select {
		case <-p.ctx.Done():
			if id == (p.workers + 1) {
				jobs <- j
				p.ctx, p.ctxCancel = context.WithCancel(context.Background())
				return
			}
		default:
			time.Sleep(time.Second)
		}

		fmt.Println("output: worker", id, j)
	}
}
