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
  extern void pm_dump_matrix (clp_object* matrix, const char* fname);

#ifdef __cplusplus
}
#endif

#endif
