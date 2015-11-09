// Test the CLP simplex model
// By Scott Pakin <pakin@lanl.gov>

package clp_test

import (
	"github.com/losalamos/clp"
	"math"
	"testing"
)

// Test if we can create a simplex model.
func TestCreateSimplex(t *testing.T) {
	_ = clp.NewSimplex()
}

// Test if we can load a problem into a simplex model.
func TestLoadProblem(t *testing.T) {
	s := clp.NewSimplex()
	m := clp.NewPackedMatrix()
	s.LoadProblem(m, nil, nil, nil, nil)
}

// closeTo says if two floating-point numbers are equal within some tolerance.
func closeTo(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// Test if we can solve a complete optimization problem with the simplex model.
func TestPrimalSolve(t *testing.T) {
	// Set up the following problem: Maximize a + b subject to {0 ≤ 2a + b
	// ≤ 10, 3 ≤ 2b − a ≤ 8}.
	mat := clp.NewPackedMatrix()
	mat.AppendColumn([]clp.Nonzero{
		{Index: 0, Value: 2.0},  // 2a
		{Index: 1, Value: -1.0}, // -a
	})
	mat.AppendColumn([]clp.Nonzero{
		{Index: 0, Value: 1.0}, // b
		{Index: 1, Value: 2.0}, // 2b
	})
	rb := []clp.Bounds{
		{Lower: 0, Upper: 10}, // [0, 10]
		{Lower: 3, Upper: 8},  // [3, 8]
	}
	obj := []float64{1.0, 1.0} // a + b
	simp := clp.NewSimplex()
	simp.LoadProblem(mat, nil, obj, rb, nil)
	simp.SetOptimizationDirection(clp.Maximize)

	// Solve the optimization problem.
	simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
	v := simp.ObjectiveValue()
	soln := simp.PrimalColumnSolution()

	// Check the results.
	if !closeTo(v, 7.6, 0.05) {
		t.Fatalf("Expected 7.6 but observed %.10g", v)
	}
	if !closeTo(soln[0], 2.4, 0.05) || !closeTo(soln[1], 5.2, 0.05) {
		t.Fatalf("Expected [2.4 5.2] but observed %v", soln)
	}
}
