package main

import (
	"github.com/longshotsyndicate/clp"
	"bytes"
	"log"
	"time"
	"fmt"
)

func main() {
	rinse(10, appendColsBatch)

}


func rinse(p int, toRinse func(n int)) {
	for i:= 0; i < p ; i++ {
		go func() {
			for {
				start := time.Now()

				toRinse(1500)


				dur := time.Now().Sub(start)

				fmt.Printf("%d %.2f\n", i, dur.Seconds())
			}
		}()

		time.Sleep(5 * time.Second)
	}
}


func noMatrixDirect(n int) {

	simp := clp.NewSimplex()


	//make a nonzero array
	col := make([]clp.Nonzero, 0)
	for j := 0; j < n; j++ {
		col = append(col, clp.Nonzero{Index: j, Value: float64(j)})
	}


	for j := 0; j < n; j++ {
		simp.BufferColumn(col)
	}



	simp.LoadProblemEfficient(nil, nil, nil, nil)

}

func appendColsBatch(n int) {

	mat := clp.NewPackedMatrix()
	simp := clp.NewSimplex()

	//make a nonzero array
	col := make([]clp.Nonzero, 0)
	for j := 0; j < n; j++ {
		col = append(col, clp.Nonzero{Index: j, Value: float64(j)})
	}

	for j := 0; j < n; j++ {
		mat.BufferColumn(col)
	}

	mat.AppendBufferedColumnsBatched()

	simp.LoadProblem(mat, nil, nil, nil, nil)

}

func oneByOneSlow(n int) {

	simp := clp.NewSimplex()
	mat := clp.NewPackedMatrix()


	//make a nonzero array
	col := make([]clp.Nonzero, 0)
	for j := 0; j < n; j++ {
		col = append(col, clp.Nonzero{Index: j, Value: float64(j)})
	}


	for j := 0; j < n; j++ {
		mat.AppendColumn(col)
	}

	var buf bytes.Buffer
	mat.DumpMatrix(&buf)
	log.Println(buf.String())

	simp.LoadProblem(mat, nil, nil, nil, nil)

}
