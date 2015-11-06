package clp

// #cgo pkg-config: clp
// #include <stdlib.h>
// #include "clp-interface.h"
import "C"
import "runtime"
import "unsafe"

// A PackedMatrix is a sparse representation of linear-programming matrix, with
// each column representing a variable, each row representing an expression,
// and each cell containing a coefficient.
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
			C.free(p)
		}
	})
	return m
}

// A Nonzero represents an element in a sparse row or column.
type Nonzero struct {
	Index int     // Zero-based element offset
	Value float64 // Value at that offset
}

// AppendCol appends a sparse column to a packed matrix.  The column is
// specified as a slice of {row number, value} pairs.
func (pm *PackedMatrix) AppendColumn(col []Nonzero) {
	// It's not safe to pass Go-allocated memory to C.  Hence, we use C's
	// malloc to allocate the memory, which we free in the PackedMatrix
	// finalizer.
	nElts := len(col)
	var cInt C.int
	intB := unsafe.Sizeof(cInt)
	rows := C.malloc(C.size_t(intB * uintptr(nElts)))
	if rows == nil {
		panic("clp: malloc failed")
	}
	pm.allocs = append(pm.allocs, rows)
	var cDbl C.double
	dblB := unsafe.Sizeof(cDbl)
	vals := C.malloc(C.size_t(dblB * uintptr(nElts)))
	if vals == nil {
		panic("clp: malloc failed")
	}
	pm.allocs = append(pm.allocs, vals)

	// Convert from the given array of two-element structs to two flat
	// vectors, and replace Go datatypes with C datatypes.
	for i, c := range col {
		ptr := unsafe.Pointer(uintptr(rows) + uintptr(i)*intB)
		*(*C.int)(ptr) = C.int(c.Index)
		ptr = unsafe.Pointer(uintptr(vals) + uintptr(i)*dblB)
		*(*C.double)(ptr) = C.double(c.Value)
	}

	// Tell our C wrapper function to append the column.
	C.pm_append_col(pm.matrix, C.int(nElts), (*C.int)(rows), (*C.double)(vals))
}
