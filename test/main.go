package main

import (
	"log"
	"github.com/longshotsyndicate/clp"
	"sync"
	"time"

	_ "net/http/pprof"
	"net/http"
)

func main() {
	//we want to start up a shitload of goroutines and then have them all spamming addcolumn and causing
	//malloc churn to measure contention

	go func() {
		log.Printf("CPU profiler running on port 6060.")
		http.ListenAndServe(":6060", nil)
	}()


	goroutines := 30
	n := 1000

	wg := sync.WaitGroup{}

	wg.Add(1)

	for i := 0; i < goroutines; i++ {
		go func() {
			for {
				start := time.Now()
				pm := clp.NewPackedMatrix()

				//make a nonzero array
				col := make([]clp.Nonzero, 0)
				for j := 0; j < n; j++ {
					col = append(col, clp.Nonzero{Index: j, Value: 0.0 + float64(n) * 0.01})
				}

				for j := 0; j < n; j++ {

					//make a bunch of non

					pm.AppendColumn(col)

				}

				dur := time.Now().Sub(start)

				log.Printf("%d %.2f", i, dur.Seconds())
			}

		}()


		time.Sleep(5*time.Second)
	}


	wg.Wait()

}
