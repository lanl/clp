// Simplex model

package clp

// #include "clp-interface.h"
// #include "stdlib.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// A Simplex represents solves linear-programming problems using the simplex
// method.
type Simplex struct {
	model  *C.clp_object    // Pointer to a ClpSimplex
	allocs []unsafe.Pointer // Row/column data to which the ClpSimplex points
}

// NewSimplex creates a new simplex model.
func NewSimplex() *Simplex {
	s := &Simplex{
		model:  C.new_simplex_model(),
		allocs: make([]unsafe.Pointer, 0, 64),
	}
	runtime.SetFinalizer(s, func(s *Simplex) {
		// When we're finished with it, free the model and all the
		// memory it referred to.
		C.free_simplex_model(s.model)
		for _, p := range s.allocs {
			cFree(p)
		}
	})
	return s
}

// Bounds represents the lower and upper bound on a value.
type Bounds struct {
	Lower float64
	Upper float64
}

// LoadProblem loads a problem into a simplex model.  It takes as an argument a
// matrix (with inequalities in rows and coefficients in columns), the upper
// and lower column bounds, the coefficients of the column objective function,
// the upper and lower row bounds, and the coefficients of the row objective
// function.  Any of these arguments except for the matrix can be nil.  When
// nil, the column bounds default to {0, ∞} for each row; the column and row
// objective functions default to 0 for all coefficients; and the row bounds
// default to {−∞, +∞} for each column.
func (s *Simplex) LoadProblem(m Matrix, cb []Bounds, obj []float64, rb []Bounds, rowObj []float64) {
	// Because of the the way the C++ API works, m can't be an arbitrary
	// implementation of the Matrix interface.  We therefore check that it
	// wraps one of the interfaces CLP knows about and abort if not.
	matrix, ok := m.(*PackedMatrix)
	if !ok {
		panic(fmt.Sprintf("clp: Simplex.LoadProblem cannot currently accept a Matrix of type %T", m))
	}

	// Get the matrix dimensions.
	nr, nc := m.Dims()

	// It's not safe to pass Go-allocated memory to C.  Hence, we use C's
	// malloc to allocate the memory, which we free in the Simplex
	// finalizer.  First, we convert cb to two C vectors, colLB and colUB.
	var colLB, colUB unsafe.Pointer
	if cb != nil {
		colLB = cMalloc(nc, C.double(0.0))
		colUB = cMalloc(nc, C.double(0.0))
		for i, b := range cb {
			cSetArrayDouble(colLB, i, b.Lower)
			cSetArrayDouble(colUB, i, b.Upper)
		}
		s.allocs = append(s.allocs, colLB, colUB)
	}

	// Next, we convert obj to a C vector, cObj.
	var cObj unsafe.Pointer
	if obj != nil {
		cObj = cMalloc(nc, C.double(0.0))
		for i, v := range obj {
			cSetArrayDouble(cObj, i, v)
		}
		s.allocs = append(s.allocs, cObj)
	}

	// Then, we convert rb to two C vectors, rowLB and rowUB.
	var rowLB, rowUB unsafe.Pointer
	if rb != nil {
		rowLB = cMalloc(nr, C.double(0.0))
		rowUB = cMalloc(nr, C.double(0.0))
		for i, b := range rb {
			cSetArrayDouble(rowLB, i, b.Lower)
			cSetArrayDouble(rowUB, i, b.Upper)
		}
		s.allocs = append(s.allocs, rowLB, rowUB)
	}

	// Finally, we convert rowObj to a C vector, rObj.
	var rObj unsafe.Pointer
	if rowObj != nil {
		rObj = cMalloc(nr, C.double(0.0))
		for i, v := range rowObj {
			cSetArrayDouble(rObj, i, v)
		}
		s.allocs = append(s.allocs, rObj)
	}

	// With all of our parameters ready, we can call our C wrapper function.
	C.simplex_load_problem(s.model, matrix.matrix,
		(*C.double)(colLB), (*C.double)(colUB), (*C.double)(cObj),
		(*C.double)(rowLB), (*C.double)(rowUB), (*C.double)(rObj))
}

// An OptDirection specifies the direction of optimization (maximize, minimize,
// or ignore).
type OptDirection float64

// These constants are used to specify the objective sense in
// Simplex.SetOptimizationDirection.
const (
	Ignore   OptDirection = 0.0  // Ignore the objective function
	Maximize              = -1.0 // Maximize the objective function
	Minimize              = 1.0  // Minimize the objective function
)

// SetOptimizationDirection specifies whether the objective function should be
// minimized, maximized, or ignored.
func (s *Simplex) SetOptimizationDirection(d OptDirection) {
	C.simplex_set_opt_dir(s.model, C.double(d))
}

// SetMaxIterations sets the maximum number of iterations for a solve.
func (s *Simplex) SetMaxIterations(maxIter int) {
	C.set_max_iterations(s.model, C.int(maxIter))
}

// MaxIterations returns the maximum number of iterations for a solve.
func (s *Simplex) MaxIterations() int {
	return int(C.max_iterations(s.model))
}

// SetMaxSeconds sets the maximum number of seconds for a solve.
func (s *Simplex) SetMaxSeconds(maxSeconds float64) {
	C.set_max_seconds(s.model, C.double(maxSeconds))
}

// MaxSeconds returns the maximum number of seconds for a solve.
func (s *Simplex) MaxSeconds() float64 {
	return float64(C.max_seconds(s.model))
}

// SecondaryStatus returns the secondary status of a model.
func (s *Simplex) SecondaryStatus() SimplexStatus {
	return SimplexStatus(C.secondary_status(s.model))
}

// WriteMPS writes the model to the named MPS file.  It returns true on success
// and false on failure.
func (s *Simplex) WriteMPS(filename string) bool {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return C.write_mps(s.model, cFilename) == 0
}

// A ValuesPass specifies whether to perform a values pass.
type ValuesPass int

// These constants specify the sort of value pass to perform.
const (
	NoValuesPass   ValuesPass = 0 // Use status variables to determine basis and solution
	DoValuesPass              = 1 // Do a values pass so variables not in the basis are given their current values and one pass of variables is done to clean up the basis with an equal or better objective value
	OnlyValuesPass            = 2 // Do only the values pass and then stop
)

// A StartFinishOptions is a bit vector for options related to the algorithm's
// initialization and finalization.
type StartFinishOptions uint

// These constants can be or'd together to specify start and finish options.
const (
	NoStartFinishOptions StartFinishOptions = 0 // Convenient name for no special options
	KeepWorkAreas                           = 1 // Do not delete work areas and factorization at end
	OldFactorization                        = 2 // Use old factorization if same number of rows
	ReduceInitialization                    = 4 // Skip as much initialization of work areas as possible
)

// SimplexStatus represents the status of a simplex optimization.
type SimplexStatus int

// These constants are the possible values for a SimplexStatus.
const (
	Optimal    SimplexStatus = 0
	Infeasible               = 1
	Unbounded                = 2
)

// These constants are the corresponding values for the secondary status
const (
	SecondaryNone                                  SimplexStatus = 0
	SecondaryPrimalInfeasible                                    = 1
	SecondaryScaledOptimalUnscaledPrimalInfeasible               = 2
	SecondaryScaledOptimalUnscaledDualInfeasible                 = 3
	SecondaryScaledOptimalUnscaledBothInfeasible                 = 4
	SecondaryGaveUp                                              = 5
	SecondaryFailedEmptyCheck                                    = 6
	SecondaryPostSolveNotOptimal                                 = 7
	SecondaryFailedBadElement                                    = 8
	SecondaryStoppedOnTime                                       = 9
	SecondaryStoppedPrimalInfeasible                             = 10
)

// Primal solves a simplex model with the primal method.
func (s *Simplex) Primal(vp ValuesPass, sfo StartFinishOptions) SimplexStatus {
	return SimplexStatus(C.simplex_primal(s.model, C.int(vp), C.int(sfo)))
}

// Dual solves a simplex model with the dual method.
func (s *Simplex) Dual(vp ValuesPass, sfo StartFinishOptions) SimplexStatus {
	return SimplexStatus(C.simplex_dual(s.model, C.int(vp), C.int(sfo)))
}

// Barrier solves a simplex model with the barrier method.  The argument says
// whether to cross over to simplex.
func (s *Simplex) Barrier(xover bool) SimplexStatus {
	var b C.int
	if xover {
		b = 1
	}
	return SimplexStatus(C.simplex_barrier(s.model, b))
}

// PrimalTolerance returns the tolerance currently associated with the
// variables in a simplex model.
func (s *Simplex) PrimalTolerance() float64 {
	var tolerance C.double
	tolerance = C.simplex_primal_get_tolerance(s.model)
	return float64(tolerance)
}

// SetPrimalTolerance assigns a new variable tolerance to a simplex model.
// According to the Clp documentation, "a variable is deemed primal feasible if
// it is less than the tolerance…below its lower bound and less than it above
// its upper bound".
func (s *Simplex) SetPrimalTolerance(tolerance float64) {
	C.simplex_primal_set_tolerance(s.model, C.double(tolerance))
}

// ReducedGradient solves a simplex model with the reduced-gradient method.
// The argument says whether to get a feasible solution (false) or to use a
// solution.
func (s *Simplex) ReducedGradient(phase bool) SimplexStatus {
	var b C.int
	if phase {
		b = 1
	}
	return SimplexStatus(C.simplex_red_grad(s.model, b))
}

// Dims returns a model's dimensions (rows and columns).
func (s *Simplex) Dims() (rows, cols int) {
	var r, c C.int
	C.simplex_get_dims(s.model, &r, &c)
	rows = int(r)
	cols = int(c)
	return
}

// Scaling indicates how problem data are to be scaled.
type Scaling int

// These constants can be passed as an argument to Simplex.SetScaling.
const (
	NoScaling          Scaling = 0 // No scaling
	EquilibriumScaling         = 1 // Equilibrium scaling
	GeometricScaling           = 2 // Geometric scaling
	AutoScaling                = 3 // Automatic scaling
	AutoInitScaling            = 4 // Automatic scaling but like the initial solve in branch-and-bound
)

// SetScaling determines how the problem data are to be scaled.
func (s *Simplex) SetScaling(sc Scaling) {
	C.simplex_scaling(s.model, C.int(sc))
}

// PrimalColumnSolution returns the primal column solution computed by a solver.
func (s *Simplex) PrimalColumnSolution() []float64 {
	_, nc := s.Dims()
	soln := make([]float64, nc)
	cSoln := C.simplex_get_prim_col_soln(s.model)
	for i := range soln {
		soln[i] = cGetArrayDouble(unsafe.Pointer(cSoln), i)
	}
	return soln
}

// DualColumnSolution returns the dual column solution computed by a solver.
func (s *Simplex) DualColumnSolution() []float64 {
	_, nc := s.Dims()
	soln := make([]float64, nc)
	cSoln := C.simplex_get_dual_col_soln(s.model)
	for i := range soln {
		soln[i] = cGetArrayDouble(unsafe.Pointer(cSoln), i)
	}
	return soln
}

// PrimalRowSolution returns the primal row solution computed by a solver.
func (s *Simplex) PrimalRowSolution() []float64 {
	_, nc := s.Dims()
	soln := make([]float64, nc)
	cSoln := C.simplex_get_prim_row_soln(s.model)
	for i := range soln {
		soln[i] = cGetArrayDouble(unsafe.Pointer(cSoln), i)
	}
	return soln
}

// DualRowSolution returns the dual row solution computed by a solver.
func (s *Simplex) DualRowSolution() []float64 {
	_, nc := s.Dims()
	soln := make([]float64, nc)
	cSoln := C.simplex_get_dual_row_soln(s.model)
	for i := range soln {
		soln[i] = cGetArrayDouble(unsafe.Pointer(cSoln), i)
	}
	return soln
}

// ObjectiveValue returns the value of the objective function after
// optimization.
func (s *Simplex) ObjectiveValue() float64 {
	return float64(C.simplex_obj_val(s.model))
}

// EasyLoadDenseProblem has no exact equivalent in the CLP library.  It is
// merely a convenient wrapper for LoadProblem that lets callers specify
// problems in a more natural, equation-like form (as opposed to CLP's normal
// matrix form).  The main limitation is that it does not provide a
// space-efficient way to represent a sparse coefficient matrix; all
// coefficients must be specified, even when zero.  A secondary limitation is
// that it does not support a row objective function.
//
// The arguments to EasyLoadDenseProblem are the coefficients of the objective
// function, lower and upper bounds on each variable, and a matrix in which
// each row is of the form {lower bound, var_1, var_2, …, var_N, upper bound}.
func (s *Simplex) EasyLoadDenseProblem(obj []float64, varBounds [][2]float64, ineqs [][]float64) {
	// Extract the lower and upper bounds for each inequality.
	nRows := len(ineqs)
	nCols := len(ineqs[0])
	rb := make([]Bounds, nRows)
	for i, row := range ineqs {
		rb[i] = Bounds{
			Lower: row[0],
			Upper: row[nCols-1],
		}
	}

	// Add each internal column one-by-one to the model.
	mat := NewPackedMatrix()
	for c := 1; c < nCols-1; c++ {
		col := make([]Nonzero, 0)
		for r, row := range ineqs {
			if row[c] != 0.0 {
				col = append(col, Nonzero{
					Index: r,
					Value: row[c],
				})
			}
		}
		mat.AppendColumn(col)
	}

	// Convert varBounds elements from [2]float64 to Bounds.
	var cb []Bounds
	if varBounds != nil {
		cb = make([]Bounds, len(varBounds))
		for b, bnd := range varBounds {
			cb[b].Lower = bnd[0]
			cb[b].Upper = bnd[1]
		}
	}

	// Load the problem into the model.
	s.LoadProblem(mat, cb, obj, rb, nil)
}
