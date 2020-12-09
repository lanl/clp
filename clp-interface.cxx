#include <ClpSimplex.hpp>
#include "clp-interface.h"

extern "C" {

  // Create a new CoinPackedMatrix.
  clp_object* new_packed_matrix (void)
  {
    CoinPackedMatrix* matrix = new CoinPackedMatrix();
    return (clp_object*)matrix;
  }

  void reserve (clp_object* matrix,
		int newMaxMajorDim,
		int newMaxSize,
		int create)
  {
    ((CoinPackedMatrix*)matrix)->reserve(newMaxMajorDim, (CoinBigIndex)newMaxSize, bool(create));
  }

  void set_dimensions (clp_object* matrix,
		int numrows,
		int numcols)
  {
    ((CoinPackedMatrix*)matrix)->setDimensions(numrows, numcols);
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

  // Retrieve a CoinPackedMatrix's data in a sparse representation.
  void pm_get_sparse_data (clp_object* matrix,
                           const int** starts,
                           const int** lengths,
                           const int** indices,
                           const double** elements)
  {
    CoinPackedMatrix* pm = (CoinPackedMatrix*)matrix;
    *starts = pm->getVectorStarts();
    *lengths = pm->getVectorLengths();
    *indices = pm->getIndices();
    *elements = pm->getElements();
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

  void simplex_primal_set_tolerance(clp_object* model, double tolerance)
  {
    ((ClpSimplex*)model)->setPrimalTolerance(tolerance);
  }

  double simplex_primal_get_tolerance(clp_object* model)
  {
    return ((ClpSimplex*)model)->primalTolerance();
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

  void set_max_iterations(clp_object* model, int max_iter)
  {
    ((ClpModel*)model)->setMaximumIterations(max_iter);
  }

  int max_iterations(clp_object* model)
  {
    return ((ClpModel*)model)->maximumIterations();
  }

  void set_max_seconds(clp_object* model, double max_seconds)
  {
    ((ClpModel*)model)->setMaximumSeconds(max_seconds);
  }

  double max_seconds(clp_object* model)
  {
    return ((ClpModel*)model)->maximumSeconds();
  }

  int secondary_status(clp_object* model)
  {
    return ((ClpModel*)model)->secondaryStatus();
  }

  int write_mps(clp_object* model, const char * filename)
  {
    return ((ClpModel*)model)->writeMps(filename);
  }

  int primal_ranging(clp_object* model, const int number_check, const int* which,
                     double* value_increase, int* sequence_increase,
                     double* value_decrease, int* sequence_decrease)
  {
    return ((ClpSimplex*)model)->primalRanging(number_check, which,
                                               value_increase, sequence_increase,
                                               value_decrease, sequence_decrease);
  }

  int dual_ranging(clp_object* model, const int number_check, const int* which,
                   double* cost_increase, int* sequence_increase,
                   double* cost_decrease, int* sequence_decrease,
                   double* value_increase, double* value_decrease)
  {
    return ((ClpSimplex*)model)->dualRanging(number_check, which,
                                             cost_increase, sequence_increase,
                                             cost_decrease, sequence_decrease,
                                             value_increase, value_decrease);
  }
}
