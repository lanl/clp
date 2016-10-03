package main

import (
	"github.com/longshotsyndicate/clp"
)

func main() {

	n := 192

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

	simp.SetOptimizationDirection(clp.Maximize)

	simp.SetScaling(clp.Scaling(2))

	simp.SetPrimalTolerance(1e-9)


	// Solve the optimization problem.
	simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)

	simp.PrimalColumnSolution()


}
