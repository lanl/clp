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
  extern void pm_dump_matrix (clp_object* matrix, const char* fname);
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

#ifdef __cplusplus
}
#endif

#endif
