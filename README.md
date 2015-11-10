clp
===

Description
-----------

The `clp` package provides a [Go](https://golang.org/) interface to the [COIN-OR Linear Programming](http://www.coin-or.org/projects/Clp.xml) (CLP) library, part of the [COIN-OR](http://www.coin-or.org/) (COmputational INfrastructure for Operations Research) suite.

[Linear programming](https://en.wikipedia.org/wiki/Linear_programming) is a method for maximizing or minimizing a linear expression subject to a set of constraints expressed as inequalities.

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
