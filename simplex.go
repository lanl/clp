// Simplex model

package clp

// #cgo pkg-config: clp
// #include <stdlib.h>
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
			C.free(p)
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
	if _, ok := m.(*PackedMatrix); !ok {
		panic(fmt.Sprintf("clp: Simplex.LoadProblem cannot currently accept a Matrix of type %T", m))
	}
}
