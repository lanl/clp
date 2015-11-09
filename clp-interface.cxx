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

  // Retrieve a CoinPackedMatrix's rows and columns.
  void pm_get_dims (clp_object* matrix, int* nrows, int* ncols)
  {
    *nrows = ((CoinPackedMatrix*)matrix)->getNumRows();
    *ncols = ((CoinPackedMatrix*)matrix)->getNumCols();
  }

  // Dump a CoinPackedMatrix to a file in a human-readable format.
  void pm_dump_matrix (clp_object* matrix, const char* fname)
  {
    ((CoinPackedMatrix*)matrix)->dumpMatrix(fname);
  }

  // Create a new ClpSimplex.
  clp_object* new_simplex_model (void)
  {
    ClpSimplex* model = new ClpSimplex();
    model->messageHandler()->setLogLevel(0);  // Not Go-like to log to a hard-wired location.
    return (clp_object*)model;
  }

  // Free an existing ClpSimplex.
  void free_simplex_model (clp_object* model)
  {
    delete (ClpSimplex*)model;
  }

  // Load a problem into a ClpSimplex.
  void simplex_load_problem (clp_object* model, clp_object* matrix,
                             const double* collb,
                             const double* colub,
                             const double* obj,
                             const double* rowlb,
                             const double* rowub,
                             const double* rowObj)
  {
    ((ClpSimplex*)model)->loadProblem(*(CoinPackedMatrix*)matrix,
                                      collb, colub, obj,
                                      rowlb, rowub, rowObj);
  }

  // Set the optimization direction.
  void simplex_set_opt_dir (clp_object* model, double dir)
  {
    ((ClpSimplex*)model)->setOptimizationDirection(dir);
  }

  // Solve a model using the primal method.
  int simplex_primal (clp_object* model, int vp, int sfo)
  {
    return ((ClpSimplex*)model)->primal(vp, sfo);
  }

  // Solve a model using the dual method.
  int simplex_dual (clp_object* model, int vp, int sfo)
  {
    return ((ClpSimplex*)model)->dual(vp, sfo);
  }

  // Solve a model using the barrier method.
  int simplex_barrier (clp_object* model, int xover)
  {
    return ((ClpSimplex*)model)->barrier(bool(xover));
  }

  // Solve a model using the reduced-gradient method.
  int simplex_red_grad (clp_object* model, int phase)
  {
    return ((ClpSimplex*)model)->reducedGradient(phase);
  }

  // Retrieve a simplex's rows and columns.
  void simplex_get_dims (clp_object* model, int* nrows, int* ncols)
  {
    *nrows = ((ClpSimplex*)model)->getNumRows();
    *ncols = ((ClpSimplex*)model)->getNumCols();
  }

  // Set or unset problem scaling.
  void simplex_scaling (clp_object* model, int mode)
  {
    ((ClpSimplex*)model)->scaling(mode);
  }

  // Return a model's primal column solution.
  double* simplex_get_prim_col_soln (clp_object* model)
  {
    return ((ClpSimplex*)model)->primalColumnSolution();
  }

  // Return a model's dual column solution.
  double* simplex_get_dual_col_soln (clp_object* model)
  {
    return ((ClpSimplex*)model)->dualColumnSolution();
  }

  // Return a model's primal row solution.
  double* simplex_get_prim_row_soln (clp_object* model)
  {
    return ((ClpSimplex*)model)->primalRowSolution();
  }

  // Return a model's dual row solution.
  double* simplex_get_dual_row_soln (clp_object* model)
  {
    return ((ClpSimplex*)model)->dualRowSolution();
  }

  // Return the value of the objective function.
  double simplex_obj_val (clp_object* model)
  {
    return ((ClpSimplex*)model)->objectiveValue();
  }
}
