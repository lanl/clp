// Test CLP packed matrixes
// By Scott Pakin <pakin@lanl.gov>

package clp_test

import (
	"github.com/losalamos/clp"
	"testing"
)

// Test if we can create a packed matrix.
func TestCreate(t *testing.T) {
	_ = clp.NewPackedMatrix()
}

// Test if we can add columns to a packed matrix.
func TestAddColumns(t *testing.T) {
	m := clp.NewPackedMatrix()
	const edge = 100
	for i := 0; i < edge; i++ {
		major := clp.Nonzero{Index: i, Value: float64(i) * 10.0}
		minor := clp.Nonzero{Index: edge - i - 1, Value: -float64(i) * 10.0}
		col := []clp.Nonzero{major, minor}
		m.AppendColumn(col)
	}
}
