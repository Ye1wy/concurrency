package main

import (
	"log"
	"strconv"
	"worker-pool/internal/pool"
)

func main() {
	pool := pool.NewWorkerPool(3)

	for i := 0; i < 1000; i++ {
		pool.Process(strconv.Itoa(i))

		if i == 2 || i == 500 {
			pool.AddWorker()
		}

		if i == 10 || i == 600 {
			if err := pool.RemoveWorker(); err != nil {
				log.Println(err)
			}
		}
	}

	pool.Close()
}
