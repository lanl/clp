#ifndef _CLP_INTERFACE_H_
#define _CLP_INTERFACE_H_

#ifdef __cplusplus
extern "C" {
#endif

  // A clp_object* is an opaque pointer to an arbitrary C++ object.
  typedef char clp_object;

  // Declare all of our wrapper functions.
  extern clp_object* new_packed_matrix (void);
  extern void free_packed_matrix (clp_object* matrix);
  extern void pm_append_col (clp_object* matrix, const int vecsize,
                             const int* vecind, const double* vecelem);
  extern void pm_get_dims (clp_object* matrix, int* nrows, int* ncols);
  extern void pm_get_sparse_data (clp_object* matrix, const int** starts,
                                  const int** lengths, const int** indices,
                                  const double** elements);
  extern clp_object* new_simplex_model (void);
  extern void free_simplex_model (clp_object* model);
  extern void simplex_load_problem (clp_object* model, clp_object* matrix,
                                    const double* collb, const double* colub,
                                    const double* obj,
                                    const double* rowlb, const double* rowub,
                                    const double* rowObj);
  extern void simplex_set_opt_dir (clp_object* model, double dir);
  extern int simplex_primal (clp_object* model, int vp, int sfo);
  extern int simplex_dual (clp_object* model, int vp, int sfo);
  extern int simplex_barrier (clp_object* model, int xover);
  extern int simplex_red_grad (clp_object* model, int phase);
  extern void simplex_get_dims (clp_object* model, int* nrows, int* ncols);
  extern void simplex_scaling (clp_object* model, int mode);
  extern double* simplex_get_prim_col_soln (clp_object* model);
  extern double* simplex_get_dual_col_soln (clp_object* model);
  extern double* simplex_get_prim_row_soln (clp_object* model);
  extern double* simplex_get_dual_row_soln (clp_object* model);
  extern double simplex_obj_val (clp_object* model);

  extern void simplex_primal_set_tolerance(clp_object* model, double tolerance);
  extern double simplex_primal_get_tolerance(clp_object* model);

  extern void set_max_iterations(clp_object* model, int max_iter);
  extern int max_iterations(clp_object* model);
  extern void set_max_seconds(clp_object* model, double max_seconds);
  extern double max_seconds(clp_object* model);

#ifdef __cplusplus
}
#endif

#endif
