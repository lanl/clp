package main

import (
	"github.com/longshotsyndicate/clp"
	"log"
)

func main() {

	n := 5

	//make a nonzero array
	col := make([]clp.Nonzero, 0)
	for j := 0; j < n; j++ {
		col = append(col, clp.Nonzero{Index: j, Value: float64(j)})
	}

	pm := clp.NewPackedMatrix()

	for j := 0; j < n; j++ {
		pm.BufferColumn(col)
	}

	pm.AppendBufferedColumnsBatched()


	simp := clp.NewSimplex()

	log.Println("lol?")
	simp.LoadProblem(pm, nil, nil, nil, nil)

	simp.SetOptimizationDirection(clp.Maximize)

	simp.SetScaling(clp.Scaling(2))

	simp.SetPrimalTolerance(1e-9)


	// Solve the optimization problem.
	simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)

	simp.PrimalColumnSolution()


}
