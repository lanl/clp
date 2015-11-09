// Test the CLP simplex model
// By Scott Pakin <pakin@lanl.gov>

package clp_test

import (
	"fmt"
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
	// Set up the following problem: Minimize a + 2b subject to {4 ≤ a + b
	// ≤ 9, -5 ≤ 3a − b ≤ 3}.
	mat := clp.NewPackedMatrix()
	mat.AppendColumn([]clp.Nonzero{
		{Index: 0, Value: 1.0}, // a
		{Index: 1, Value: 3.0}, // 3a
	})
	mat.AppendColumn([]clp.Nonzero{
		{Index: 0, Value: 1.0},  // b
		{Index: 1, Value: -1.0}, // -b
	})
	rb := []clp.Bounds{
		{Lower: 4, Upper: 9},  // [4, 9]
		{Lower: -5, Upper: 3}, // [-5, 3]
	}
	obj := []float64{1.0, 2.0} // a + 2b
	simp := clp.NewSimplex()
	simp.LoadProblem(mat, nil, obj, rb, nil)
	simp.SetOptimizationDirection(clp.Minimize)

	// Solve the optimization problem.
	simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
	v := simp.ObjectiveValue()
	soln := simp.PrimalColumnSolution()

	// Check the results.
	if !closeTo(soln[0], 1.75, 0.005) || !closeTo(soln[1], 2.25, 0.005) {
		t.Fatalf("Expected [1.75 2.25] but observed %v", soln)
	}
	if !closeTo(v, 6.25, 0.005) {
		t.Fatalf("Expected 6.25 but observed %.10g", v)
	}
}

// Maximize a + b subject to both 0 ≤ 2a + b ≤ 10 and 3 ≤ 2b − a ≤ 8.
func Example_complete() {
	// Set up the problem.
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
	val := simp.ObjectiveValue()
	soln := simp.PrimalColumnSolution()

	// Output the results.
	fmt.Printf("a = %.1f\nb = %.1f\na + b = %.1f\n", soln[0], soln[1], val)
	// Output:
	// a = 2.4
	// b = 5.2
	// a + b = 7.6
}
