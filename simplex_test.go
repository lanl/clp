// Test the CLP simplex model
// By Scott Pakin <pakin@lanl.gov>

package clp_test

import (
	"github.com/losalamos/clp"
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
