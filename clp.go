// Copyright © 2015, Los Alamos National Security, LLC
// All rights reserved.

/*
Package clp provides access to the COIN-OR Linear Program (CLP)
library.  As the name implies, CLP is a solver for linear-programming
problems:

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
value of a + b is 7.6, which is achieved when a = 2.4 and b = 5.2.

The Go package currently implements only a tiny subset of the CLP library
but a subset that suffices to solve basic optimization problems.

Relevant URLs:

• COIN-OR (COmputational INfrastructure for Operations Research): http://www.coin-or.org/

• LP (Linear Programming): https://en.wikipedia.org/wiki/Linear_programming

• CLP (COIN-OR Linear Programming): http://www.coin-or.org/projects/Clp.xml
*/
package clp

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
