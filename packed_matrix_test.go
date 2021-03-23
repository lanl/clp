// Test CLP packed matrixes
// By Scott Pakin <pakin@lanl.gov>

package clp_test

import (
	"bytes"
	"github.com/lanl/clp"
	"testing"
)

// Test if we can create a packed matrix.
func TestCreateMatrix(t *testing.T) {
	_ = clp.NewPackedMatrix()
}

// Test if we can reserve space in a packed matrix.
func TestReserve(t *testing.T) {
	m := clp.NewPackedMatrix()
	m.Reserve(100, 100, false)
}

// Append a given number of columns to a (presumably empty) matrix to produce
// an NRxNC sparse matrix with nonzeroes only on the major and minor diagonals.
func addColumns(m *clp.PackedMatrix, nr, nc int) {
	for i := 0; i < nc; i++ {
		major := clp.Nonzero{Index: i % nr, Value: float64(i) * 10.0}
		minor := clp.Nonzero{Index: nr - i%nr - 1, Value: -float64(i) * 10.0}
		col := []clp.Nonzero{major, minor}
		m.AppendColumn(col)
	}
}

// Test if we can add columns to a packed matrix.
func TestAddColumns(t *testing.T) {
	m := clp.NewPackedMatrix()
	addColumns(m, 1000, 1000)
}

// Test if we can remove columns from a packed matrix.
func TestDeleteColumns(t *testing.T) {
	m := clp.NewPackedMatrix()
	addColumns(m, 100, 5)
	_, c := m.Dims()
	if c != 5 {
		t.Fatalf("Expected 5 colunms but saw %d", c)
	}
	m.DeleteColumns([]int{1, 3, 4})
	_, c = m.Dims()
	if c != 2 {
		t.Fatalf("Expected 2 colunms but saw %d", c)
	}
}

// Test if we can query a packed matrix's dimensions.
func TestDims(t *testing.T) {
	for _, trials := range [...][2]int{
		{66, 7},  // Wide
		{55, 68}, // Tall
	} {
		nr, nc := trials[0], trials[1]
		m := clp.NewPackedMatrix()
		addColumns(m, nr, nc)
		r, c := m.Dims()
		if r != nr || c != nc {
			t.Fatalf("Expected %dx%d but saw %dx%d", nr, nc, r, c)
		}
	}
}

// Test if we can add columns to a packed matrix and get the expected dump in
// return.
func TestDumpMatrix(t *testing.T) {
	m := clp.NewPackedMatrix()
	addColumns(m, 5, 5)
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
