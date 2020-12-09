/*
Package clp provides an interface to the COIN-OR Linear Programming (CLP)
library.  As the name implies, CLP is a solver for linear-programming problems:

    CLP is a high quality open-source LP solver. Its main strengths are its
    Dual and Primal Simplex algorithms. It also has a barrier algorithm for
    Linear and Quadratic objectives. There are limited facilities for Nonlinear
    and Quadratic objectives using the Simplex algorithm. It is available as a
    library and as a standalone solver. It was written by John Forrest, jjforre
    at us.ibm.com

Linear programming is an optimization technique.  Given an objective function,
such as a + b, and a set of constraints in the form of linear inequalities,
such as 0 ≤ 2a + b ≤ 10 and 3 ≤ 2b − a ≤ 8, find values for the variables that
maximizes or minimizes the objective function.  In this example, the maximum
value of a + b is 7.6, which is achieved when a = 2.4 and b = 5.2.  The example
code associated with Simplex.LoadProblem shows how to set up and solve this
precise problem using an API based directly on CLP's C++ API.  The example code
associated with Simplex.EasyLoadDenseProblem shows how to specify the same
problem using a more equation-oriented API specific to the clp package.

The clp package currently exposes only a tiny subset of the CLP library.

Relevant URLs:

• COIN-OR (COmputational INfrastructure for Operations Research): http://www.coin-or.org/

• LP (Linear Programming): https://en.wikipedia.org/wiki/Linear_programming

• CLP (COIN-OR Linear Programming): http://www.coin-or.org/projects/Clp.xml
*/
package clp

// #cgo pkg-config: clp
// #include <stdlib.h>
import "C"
import (
	"reflect"
	"unsafe"
)

// A Nonzero represents an element in a sparse row or column.
type Nonzero struct {
	Index int     // Zero-based element offset
	Value float64 // Value at that offset
}

// A Matrix sparsely represents a set of linear expressions.  Each column
// represents a variable, each row represents an expression, and each cell
// containing a coefficient.  Bounds on rows and columns are applied during
// model initialization.
type Matrix interface {
	AppendColumn(col []Nonzero) // Append a column given values for all of its nonzero elements
	Dims() (rows, cols int)     // Return the matrix's dimensions
}

// cMalloc asks C to allocate memory.  For convenience to Go, the arguments
// are like calloc's except that the size argument is a value, which cMalloc
// will take the size of.  cMalloc panics on error (typically, out of memory).
func cMalloc(nmemb int, sizeVal interface{}) unsafe.Pointer {
	size := reflect.TypeOf(sizeVal).Size()
	mem := C.malloc(C.size_t(uintptr(nmemb) * size))
	if mem == nil {
		panic("clp: malloc failed")
	}
	return mem
}

// cFree asks C to free memory.
func cFree(mem unsafe.Pointer) {
	C.free(mem)
}

// cSetArrayInt assigns a[i] = v where a is a C.int array allocated by
// cMalloc and i and v are Go ints.
func cSetArrayInt(a unsafe.Pointer, i, v int) {
	eSize := unsafe.Sizeof(C.int(0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	*(*C.int)(ptr) = C.int(v)
}

// c__GetArrayInt returns a[i] as a Go int where a is a C.int array allocated
// by cMalloc and i is a Go ints.
func cGetArrayInt(a unsafe.Pointer, i int) int {
	eSize := unsafe.Sizeof(C.int(0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	return int(*(*C.int)(ptr))
}

// cSetArrayDouble assigns a[i] = v where a is a C.double array allocated by
// cMalloc, i is an int, and v is a Go float64.
func cSetArrayDouble(a unsafe.Pointer, i int, v float64) {
	eSize := unsafe.Sizeof(C.double(0.0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	*(*C.double)(ptr) = C.double(v)
}

// cGetArrayDouble returns a[i] as a Go float64 where a is a C.double array
// allocated by cMalloc and i is an int.
func cGetArrayDouble(a unsafe.Pointer, i int) float64 {
	eSize := unsafe.Sizeof(C.double(0.0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	return float64(*(*C.double)(ptr))
}

// cNewArrayDoubleFrom creates a new c double array
// and initializes it from the given slice.
func cNewArrayDoubleFrom(s []float64) unsafe.Pointer {
	n := len(s)
	cArray := cMalloc(n, C.double(0.0))
	for i, v := range s {
		cSetArrayDouble(cArray, i, v)
	}
	return cArray
}

// cNewArrayIntFrom creates a new c int array
// and initializes it from the given slice.
func cNewArrayIntFrom(s []int) unsafe.Pointer {
	n := len(s)
	cArray := cMalloc(n, C.int(0))
	for i, v := range s {
		cSetArrayInt(cArray, i, v)
	}
	return cArray
}

// cCopyArrayInt copies the given c int array into the slice.
func cCopyArrayInt(s []int, u unsafe.Pointer) {
	for i := range s {
		s[i] = cGetArrayInt(u, i)
	}
}

// cCopyArrayDouble copies the given c double array into the slice.
func cCopyArrayDouble(s []float64, u unsafe.Pointer) {
	for i := range s {
		s[i] = cGetArrayDouble(u, i)
	}
}
