clp
===

Description
-----------

The `clp` package provides a [Go](https://golang.org/) interface to the [COIN-OR Linear Programming](http://www.coin-or.org/projects/Clp.xml) (CLP) library, part of the [COIN-OR](http://www.coin-or.org/) (COmputational INfrastructure for Operations Research) suite.

[Linear programming](https://en.wikipedia.org/wiki/Linear_programming) (LP) is a method for maximizing or minimizing a linear expression subject to a set of constraints expressed as inequalities.  As an example that's simple enough to solve by hand, what roll of three six-sided dice has the largest total value if no two dice are allowed have the same value and the difference in value between the first and second largest dice must be smaller than the difference in value between the second and third largest dice?  From an LP standpoint, the objective function we need to maximize to answer that question is *a* + *b* + *c*, where each variable represents the value on one die.  The first constraint is that each die be six sided:

* 1 ≤ *a* ≤ 6
* 1 ≤ *b* ≤ 6
* 1 ≤ *c* ≤ 6

The second constraint is that the three dice all have different values.  We specify this by imposing a total order, *a* > *b* > *c*, which we express as

* 1 ≤ a - b ≤ ∞
* 1 ≤ b - c ≤ ∞

The third constraint is that the difference in value between the first and second largest dice (*a* − *b*) is smaller than the difference in value between the second and third largest dice (*b* − *c*).  To put this in a suitable format for LP, we rewrite *a* − *b* < *b* − *c* as

* -∞ ≤ a - 2b + c ≤ -1

These constraints translate directly to Go using the `clp` package's API:
```Go
package main

import (
        "fmt"
        "github.com/losalamos/clp"
        "math"
)

func main() {
        // Set up the optimization problem.
        pinf := math.Inf(1)
        ninf := math.Inf(-1)
        simp := clp.NewSimplex()
        simp.EasyLoadDenseProblem(
                //         A    B    C
                []float64{1.0, 1.0, 1.0},
                [][2]float64{
                        // LB UB
                        {1, 6}, // 1 ≤ a ≤ 6
                        {1, 6}, // 1 ≤ b ≤ 6
                        {1, 6}, // 1 ≤ c ≤ 6
                },
                [][]float64{
                        // LB  A    B    C    UB
                        {1.0, 1.0, -1.0, 0.0, pinf},  // 1 ≤ a - b ≤ ∞
                        {1.0, 0.0, 1.0, -1.0, pinf},  // 1 ≤ b - c ≤ ∞
                        {ninf, 1.0, -2.0, 1.0, -1.0}, // -∞ ≤ a - 2b + c ≤ -1
                })
        simp.SetOptimizationDirection(clp.Maximize)

        // Solve the optimization problem.
        simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
        soln := simp.PrimalColumnSolution()

        // Output the solution.
        fmt.Printf("Die 1 = %.0f\n", soln[0])
        fmt.Printf("Die 2 = %.0f\n", soln[1])
        fmt.Printf("Die 3 = %.0f\n", soln[2])
}
```

The output is the expected
```
Die 1 = 6
Die 2 = 5
Die 3 = 3
```

Installation
------------

`clp` has been tested only on Linux.  The package requires a CLP installation to build.  To check if CLP is installed, ensure that the following command produces a list of libraries, typically along the lines of `-lClp -lCoinUtils …`, and, more importantly, issues no error messages:
```bash
pkg-config --libs clp
```

Once CLP installation is confirmed, install the `clp` package with [`go get`](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies):
```bash
go get github.com/losalamos/clp
```

Documentation
-------------

Pre-built documentation for the `clp` API is available online at <http://godoc.org/github.com/losalamos/clp>, courtesy of [GoDoc](http://godoc.org/).

License
-------

`clp` is provided under a BSD-ish license with a "modifications must be indicated" clause.  See [the LICENSE file](http://github.com/losalamos/clp/blob/master/LICENSE.md) for the full text.

This package is part of the [LANL Go Suite](http://www.lanl.gov/projects/feynman-center/technologies/software/lanl%20go%20suite.php), LA-CC-11-056.

Author
------

Scott Pakin, <pakin@lanl.gov>
