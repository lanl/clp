// Packed matrices

package clp

// #include "clp-interface.h"
import "C"
import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"unsafe"
)

// A PackedMatrix is a basic implementation of the Matrix interface.
type PackedMatrix struct {
	matrix *C.clp_object    // Pointer to a CoinPackedMatrix
	allocs []unsafe.Pointer // Row/column data to which the CoinPackedMatrix points
}

// NewPackedMatrix allocates a new, empty, packed matrix.
func NewPackedMatrix() *PackedMatrix {
	m := &PackedMatrix{
		matrix: C.new_packed_matrix(),
		allocs: make([]unsafe.Pointer, 0, 64),
	}
	runtime.SetFinalizer(m, func(m *PackedMatrix) {
		// When we're finished with it, free the matrix and all the
		// memory it referred to.
		C.free_packed_matrix(m.matrix)
		for _, p := range m.allocs {
			c_free(p)
		}
	})
	return m
}

// AppendColumn appends a sparse column to a packed matrix.  The column is
// specified as a slice of {row number, value} pairs.
func (pm *PackedMatrix) AppendColumn(col []Nonzero) {
	// It's not safe to pass Go-allocated memory to C.  Hence, we use C's
	// malloc to allocate the memory, which we free in the PackedMatrix
	// finalizer.
	nElts := len(col)
	rows := c_malloc(nElts, C.int(0))
	pm.allocs = append(pm.allocs, rows)
	vals := c_malloc(nElts, C.double(0.0))
	pm.allocs = append(pm.allocs, vals)

	// Convert from the given array of two-element structs to two flat
	// vectors, and replace Go datatypes with C datatypes.
	for i, c := range col {
		c_SetArrayInt(rows, i, c.Index)
		c_SetArrayDouble(vals, i, c.Value)
	}

	// Tell our C wrapper function to append the column.
	C.pm_append_col(pm.matrix, C.int(nElts), (*C.int)(rows), (*C.double)(vals))
}

// Dims returns a packed matrix's dimensions (rows and columns).
func (pm *PackedMatrix) Dims() (rows, cols int) {
	var r, c C.int
	C.pm_get_dims(pm.matrix, &r, &c)
	rows = int(r)
	cols = int(c)
	return
}

// DumpMatrix outputs a packed matrix in a human-readable format.  This method
// is intended primarily to help with testing and debugging.
func (pm *PackedMatrix) DumpMatrix(w io.Writer) error {
	// CLP's dumpMatrix function accepts a filename as an argument (or NULL
	// for standard output).  To make DumpMatrix more Go-like, we write to
	// a temporary file, then read the result back into an io.Writer.  Yes,
	// that's quite kludgy, but this method is intended to be primarily a
	// test/debug function, not a critical component of application
	// execution.
	out, err := ioutil.TempFile("", "clp-")
	if err != nil {
		return err
	}
	outName := out.Name()
	defer os.Remove(outName)
	fn := C.CString(outName)
	defer c_free(unsafe.Pointer(fn))
	C.pm_dump_matrix(pm.matrix, fn)
	in, err := os.Open(outName)
	if err != nil {
		return err
	}
	defer in.Close()
	_, err = io.Copy(w, in)
	return err
}
