#include <coin/ClpSimplex.hpp>
#include "clp-interface.h"

extern "C" {

  // Create a new CoinPackedMatrix.
  clp_object* new_packed_matrix (void)
  {
    CoinPackedMatrix* matrix = new CoinPackedMatrix();
    return (clp_object*)matrix;
  }

  // Free an existing CoinPackedMatrix.
  void free_packed_matrix (clp_object* matrix)
  {
    delete (CoinPackedMatrix*)matrix;
  }

  // Append a (sparse) column to a CoinPackedMatrix.
  void pm_append_col (clp_object* matrix, const int vecsize,
		      const int* vecind, const double* vecelem)
  {
    ((CoinPackedMatrix*)matrix)->appendCol(vecsize, vecind, vecelem);
  }

  // Dump a CoinPackedMatrix to a file in a human-readable format.
  void pm_dump_matrix (clp_object* matrix, const char* fname)
  {
    ((CoinPackedMatrix*)matrix)->dumpMatrix(fname);
  }  
}
