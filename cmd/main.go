package main

import (
	"strconv"
	"time"
	"worker-pool/internal/pool"
)

func main() {
	// data := []string{"aboba", "amogus", "biba", "boba", "phoba", "phoga", "sigame", "cross", "phos", "nok", "kok", "hoh", "HOA", "pipa", "pipe"}

	workers := pool.NewWorkerPool(3)

	for i := 0; i < 5000; i++ {
		workers.ProcessData(strconv.Itoa(i))

		if i == 2 || i == 4 || i == 400 {
			workers.Add()
		}

		if i == 5 || i == 1000 {
			workers.Delete()
		}
	}

	workers.CloseJobChan()

	time.Sleep(30 * time.Second)
}
