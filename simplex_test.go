// Test the CLP simplex model
// By Scott Pakin <pakin@lanl.gov>

package clp_test

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/lanl/clp"
)

// Test if we can create a simplex model.
func TestCreateSimplex(t *testing.T) {
	_ = clp.NewSimplex()
}

// Test if we can set solve iteration limit.
func TestSimplexSetIters(t *testing.T) {
	s := clp.NewSimplex()
	maxIter := 10
	s.SetMaxIterations(maxIter)
	maxIterBack := s.MaxIterations()
	if maxIter != maxIterBack {
		t.Fatal("Cannot set max iterations")
	}
}

// Test if we can set solve time limit.
func TestSimplexSetSeconds(t *testing.T) {
	s := clp.NewSimplex()
	maxSeconds := 12.1
	s.SetMaxSeconds(maxSeconds)
	maxSecondsBack := s.MaxSeconds()
	if !closeTo(maxSeconds/maxSecondsBack, 1.0, 0.01) {
		t.Fatalf("Cannot set max seconds (wanted %v but saw %v", maxSeconds, maxSecondsBack)
	}
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
	secStatus := simp.SecondaryStatus()
	if secStatus != clp.SecondaryNone {
		t.Fatalf("Expected %d secondary status but got %d", clp.SecondaryNone, secStatus)
	}
}

func TestZeroSolve(t *testing.T) {
	mat := clp.NewPackedMatrix()
	mat.AppendColumn([]clp.Nonzero{
		{Index: 0, Value: 1.0},
	})
	// force a second all-0 row into the matrix
   mat.SetDimensions(2,1)
	rb := []clp.Bounds{
		{Lower: 0, Upper: 0},
		{Lower: 0, Upper: 0},
	}
	obj := []float64{1.0}
	simp := clp.NewSimplex()
	simp.LoadProblem(mat, nil, obj, rb, nil)
	simp.SetOptimizationDirection(clp.Minimize)

	// Solve the optimization problem.
	simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
	soln := simp.PrimalColumnSolution()
   if soln == nil {
      t.Error("got nil solution when testing Zero case")
   }
	// the real sign of success is that we got here without a panic

}

// Test if we can solve the same problem as above but with the "easy" interface.
func TestEasyPrimalSolve(t *testing.T) {
	// Set up the following problem: Minimize a + 2b subject to {4 ≤ a + b
	// ≤ 9, -5 ≤ 3a − b ≤ 3}.
	simp := clp.NewSimplex()
	simp.EasyLoadDenseProblem(
		[]float64{1.0, 2.0}, // a + 2b
		nil,                 // No explicit bounds on A or B
		[][]float64{
			// LB  A    B   UB
			{4.0, 1.0, 1.0, 9.0},   // 4 ≤ a + b ≤ 9
			{-5.0, 3.0, -1.0, 3.0}, // -5 ≤ 3a − b ≤ 3
		})
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
	secStatus := simp.SecondaryStatus()
	if secStatus != clp.SecondaryNone {
		t.Fatalf("Expected %d secondary status but got %d", clp.SecondaryNone, secStatus)
	}
}

// Ensure that we can both query and change the primal tolerance used in a
// simplex model.
func TestGetSetSimplexPrimalTolerance(t *testing.T) {
	simp := clp.NewSimplex()

	initial := simp.PrimalTolerance()
	simp.SetPrimalTolerance(initial * 2.0)
	reset := simp.PrimalTolerance()

	if reset != initial*2.0 {
		t.Fatalf("Expected %f but observed %f", initial*2.0, reset)
	}
}

// Maximize a + b subject to both 0 ≤ 2a + b ≤ 10 and 3 ≤ 2b − a ≤ 8.
func ExampleSimplex_LoadProblem() {
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

// Maximize a + b subject to both 0 ≤ 2a + b ≤ 10 and 3 ≤ 2b − a ≤ 8.
func ExampleSimplex_EasyLoadDenseProblem() {
	// Set up the problem.
	simp := clp.NewSimplex()
	simp.EasyLoadDenseProblem(
		//         A    B
		[]float64{1.0, 1.0}, // a + b
		nil,                 // No explicit bounds on A or B
		[][]float64{
			// LB  A    B    UB
			{0.0, 2.0, 1.0, 10.0}, // 0 ≤ 2a + b ≤ 10
			{3.0, -1.0, 2.0, 8.0}, // 3 ≤ -a + 2b ≤ 8
		})
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

// Test if we can solve a problem with far more inequalities than variables.
func TestEasyManyIneqs(t *testing.T) {
	// Set up the following problem: Minimize a subject to {1 ≤ a, 2 ≤ a,
	// …, N ≤ a}.
	const nIneqs = 100
	simp := clp.NewSimplex()
	inf := math.Inf(1)
	ineqs := make([][]float64, nIneqs)
	for i := range ineqs {
		ineqs[i] = []float64{float64(i + 1), 1.0, inf}
	}
	simp.EasyLoadDenseProblem([]float64{1.0}, nil, ineqs)
	simp.SetOptimizationDirection(clp.Minimize)

	// Solve the optimization problem.
	simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
	v := simp.ObjectiveValue()
	soln := simp.PrimalColumnSolution()

	// Check the results.
	if !closeTo(soln[0], 100.0, 0.5) {
		t.Fatalf("Expected [100] but observed %v", soln)
	}
	if !closeTo(v, 100, 0.5) {
		t.Fatalf("Expected 100 but observed %.10g", v)
	}
	secStatus := simp.SecondaryStatus()
	if secStatus != clp.SecondaryNone {
		t.Fatalf("Expected %d secondary status but got %d", clp.SecondaryNone, secStatus)
	}
}

// Test if we can write an optimization problem to an MPS file.
func TestWriteMPS(t *testing.T) {
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

	// Write it to a file.  We don't verify the contents because these
	// could potentially change, even in non-meaningful ways (e.g.,
	// spacing), across versions of the CLP library.
	mps, err := ioutil.TempFile("", "clp-*.mps")
	if err != nil {
		t.Fatalf("Failed to create a temporary MPS file (%v)", err)
	}
	mpsName := mps.Name()
	mps.Close()
	defer os.Remove(mpsName)
	if !simp.WriteMPS(mpsName) {
		t.Fatalf("Failed to write a simplex model to %s", mpsName)
	}
}
