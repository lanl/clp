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
