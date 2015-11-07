// Simplex model

package clp

// #include "clp-interface.h"
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
			c_free(p)
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
		colLB = c_malloc(nc, C.double(0.0))
		colUB = c_malloc(nc, C.double(0.0))
		for i, b := range cb {
			c_SetArrayDouble(colLB, i, b.Lower)
			c_SetArrayDouble(colUB, i, b.Upper)
		}
		s.allocs = append(s.allocs, colLB, colUB)
	}

	// Next, we convert obj to a C vector, cObj.
	var cObj unsafe.Pointer
	if obj != nil {
		cObj = c_malloc(nc, C.double(0.0))
		for i, v := range obj {
			c_SetArrayDouble(cObj, i, v)
		}
		s.allocs = append(s.allocs, cObj)
	}

	// Then, we convert rb to two C vectors, rowLB and rowUB.
	var rowLB, rowUB unsafe.Pointer
	if rb != nil {
		rowLB = c_malloc(nr, C.double(0.0))
		rowUB = c_malloc(nr, C.double(0.0))
		for i, b := range rb {
			c_SetArrayDouble(rowLB, i, b.Lower)
			c_SetArrayDouble(rowUB, i, b.Upper)
		}
		s.allocs = append(s.allocs, rowLB, rowUB)
	}

	// Finally, we convert rowObj to a C vector, rObj.
	var rObj unsafe.Pointer
	if rowObj != nil {
		rObj = c_malloc(nr, C.double(0.0))
		for i, v := range rowObj {
			c_SetArrayDouble(rObj, i, v)
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
	KeepWorkAreas        StartFinishOptions = 1 // Do not delete work areas and factorization at end
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
