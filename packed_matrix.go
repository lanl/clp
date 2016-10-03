// Packed matrices

package clp

// #include "clp-interface.h"
import "C"
import (
	"fmt"
	"io"
	"unsafe"
	"log"
	"runtime"
)

// A PackedMatrix is a basic implementation of the Matrix interface.
type PackedMatrix struct {
	matrix *C.clp_object    // Pointer to a CoinPackedMatrix
	allocs []unsafe.Pointer // Row/column data to which the CoinPackedMatrix points

	pendingColumns [][]Nonzero
	totalDataLen int
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




// AppendColumn appends a sparse column to a packed matrix.  The column is
// specified as a slice of {row number, value} pairs.
// Additionally, we do not call into cgo at this point or incur any malloc hits
// As this trashes multicore performance. Once all columns are buffered, we can call a final
// AppendBufferedColumnsBatched() to flush all waiting columns to the matrix for solving
func (pm *PackedMatrix) BufferColumn(col []Nonzero) {
	pm.pendingColumns = append(pm.pendingColumns, col)
	pm.totalDataLen += len(col)
}


//Flushes all buffered columns to the matrix in one go with minimal malloc calls
func (pm *PackedMatrix) AppendBufferedColumnsBatched() {
	//so we need to allocate a few chunks of memory here to fit the
	//CoinPackedMatrix::appendCols signature

	numCols := len(pm.pendingColumns)


	if numCols == 0 {
		return
	}


	columnStarts := c_malloc(numCols+1, C.int(0))
	rowIndices := c_malloc(pm.totalDataLen, C.int(0))
	rowElements := c_malloc(pm.totalDataLen, C.double(0.0))
	pm.allocs = append(pm.allocs, columnStarts)
	pm.allocs = append(pm.allocs, rowIndices)
	pm.allocs = append(pm.allocs, rowElements)

	dataPosition := 0

	for col, colData := range pm.pendingColumns {
		//columnStarts[col] = C.int(dataPosition)
		c_SetArrayInt(columnStarts, col, dataPosition)
		for _, nz := range colData {

			c_SetArrayInt(rowIndices, dataPosition, nz.Index)
			c_SetArrayDouble(rowElements, dataPosition, nz.Value)

			dataPosition++
		}
	}

	c_SetArrayInt(columnStarts, len(pm.pendingColumns), dataPosition)



	C.pm_append_cols(
		pm.matrix,
		C.int(numCols),
		(*C.int)(columnStarts),
		(*C.int)(rowIndices),
		(*C.double)(rowElements),
		C.int(0),
	)

	pm.pendingColumns = nil
	pm.totalDataLen = 0
}

func (pm *PackedMatrix) Reserve(cols int, maxSize int, create int) {
	//  extern void reserve_packed_matrix (clp_object* model, int newMaxMajorDim, int newMaxSize, int create);
	C.reserve_packed_matrix(pm.matrix, C.int(cols), C.int(maxSize), C.int(create))
}


// AppendColumn appends a sparse column to a packed matrix.  The column is
// specified as a slice of {row number, value} pairs.
func (pm *PackedMatrix) AppendColumns(cols [][]Nonzero) {
	for _, col := range cols {
		pm.AppendColumn(col)
	}
}

// Dims returns a packed matrix's dimensions (rows and columns).
func (pm *PackedMatrix) Dims() (rows, cols int) {
	var r, c C.int
	C.pm_get_dims(pm.matrix, &r, &c)
	rows = int(r)
	cols = int(c)
	return
}

// SparseData returns a packed matrix's data in a sparse representation.  It
// corresponds to the getVectorStarts(), getVectorLengths(), getIndices(), and
// getElements() methods in the CLP library's CoinPackedMatrix class.
func (pm *PackedMatrix) SparseData() (starts, lengths, indices []int, elements []float64) {
	// Retrieve pointers into the matrix's internal state.
	var cstarts *C.int
	var clens *C.int
	var cidxs *C.int
	var celts *C.double
	C.pm_get_sparse_data(pm.matrix, &cstarts, &clens, &cidxs, &celts)

	// Convert from C arrays to Go slices.  We assume column ordering
	// because we don't yet give the user the ability to change the
	// ordering from the default column-ordered.
	_, nc := pm.Dims()
	starts = make([]int, 0)
	lengths = make([]int, 0)
	log.Println("nc", nc)
	for i := 0; i < nc ; i++ {
		starts = append(starts, c_GetArrayInt(unsafe.Pointer(cstarts), i))
		lengths = append(lengths, c_GetArrayInt(unsafe.Pointer(clens), i))
	}
	indices = make([]int, 0, nc)
	elements = make([]float64, 0, nc)
	log.Println(starts)
	log.Println(lengths)
	for i := 0; i < nc; i++ {
		for j := starts[i]; j < starts[i]+lengths[i]; j++ {
			indices = append(indices, c_GetArrayInt(unsafe.Pointer(cidxs), j))
			elements = append(elements, c_GetArrayDouble(unsafe.Pointer(celts), j))
		}
	}
	return
}

// DenseData returns a packed matrix's data in a dense representation.  This
// method has no exact equivalent in the CLP library.  It is merely a
// convenient wrapper for SparseMatrix that makes it easy to work with smaller
// matrices.
func (pm *PackedMatrix) DenseData() [][]float64 {
	// Create a dense matrix to populate and return.
	nr, nc := pm.Dims()
	mat := make([][]float64, nr)
	for r := range mat {
		mat[r] = make([]float64, nc)
	}

	// Populate the dense matrix from the sparse representation.
	starts, lengths, indices, elements := pm.SparseData()
	for c, st := range starts {
		iend := st + lengths[c]
		for i := st; i < iend; i++ {
			r := indices[i]
			mat[r][c] = elements[i]
		}
	}
	return mat
}

// DumpMatrix outputs a packed matrix in a human-readable format.  This method
// is intended primarily to help with testing and debugging.
func (pm *PackedMatrix) DumpMatrix(w io.Writer) error {
	// Reproduce CoinPackedMatrix::dumpMatrix() from CoinPackedMatrix.cpp.
	// We don't call the original C++ method because it writes to a file,
	// while we'd prefer to use an io.Writer.
	starts, lengths, indices, elements := pm.SparseData()
	var err error
	printf := func(format string, a ...interface{}) {
		// Borrow the error-checking trick from "Errors are values"
		// (https://blog.golang.org/errors-are-values).
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, format, a...)
	}
	printf("Dumping matrix...\n\n")
	printf("colordered: %d\n", 1) // Only column ordered is currently supported.
	minor, major := pm.Dims()
	printf("major: %d   minor: %d\n", major, minor)
	for i := 0; i < major; i++ {
		printf("vec %d has length %d with entries:\n", i, lengths[i])
		for j := starts[i]; j < starts[i]+lengths[i]; j++ {
			printf("        %15d  %40.25f\n", indices[j], elements[j])
		}
	}
	printf("\nFinished dumping matrix\n")
	return err
}
