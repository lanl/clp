// Test CLP packed matrixes
// By Scott Pakin <pakin@lanl.gov>

package clp_test

import (
	"bytes"
	"github.com/losalamos/clp"
	"testing"
)

// Test if we can create a packed matrix.
func TestCreateMatrix(t *testing.T) {
	_ = clp.NewPackedMatrix()
}

// Append a given number of columns to a (presumably empty) matrix to produce
// an NxN sparse matrix with nonzeroes only on the major and minor diagonals.
func addColumns(m *clp.PackedMatrix, n int) {
	for i := 0; i < n; i++ {
		major := clp.Nonzero{Index: i, Value: float64(i) * 10.0}
		minor := clp.Nonzero{Index: n - i - 1, Value: -float64(i) * 10.0}
		col := []clp.Nonzero{major, minor}
		m.AppendColumn(col)
	}
}

// Test if we can add columns to a packed matrix.
func TestAddColumns(t *testing.T) {
	m := clp.NewPackedMatrix()
	addColumns(m, 1000)
}

// Test if we can add columns to a packed matrix and get the expected dump in
// return.
func TestDumpMatrix(t *testing.T) {
	m := clp.NewPackedMatrix()
	addColumns(m, 5)
	expected := `Dumping matrix...

colordered: 1
major: 5   minor: 5
vec 0 has length 2 with entries:
                      0               0.0000000000000000000000000
                      4              -0.0000000000000000000000000
vec 1 has length 2 with entries:
                      1              10.0000000000000000000000000
                      3             -10.0000000000000000000000000
vec 2 has length 2 with entries:
                      2              20.0000000000000000000000000
                      2             -20.0000000000000000000000000
vec 3 has length 2 with entries:
                      3              30.0000000000000000000000000
                      1             -30.0000000000000000000000000
vec 4 has length 2 with entries:
                      4              40.0000000000000000000000000
                      0             -40.0000000000000000000000000

Finished dumping matrix
`
	var buf bytes.Buffer
	m.DumpMatrix(&buf)
	actual := buf.String()
	if actual != expected {
		t.Logf("Expected output follows:\n%s", expected)
		t.Logf("Actual output follows:\n%s", actual)
		t.Fatalf("Mismatch between expected and actual matrix contents")
	}
}
