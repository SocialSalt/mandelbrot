#ifndef COMPLEX_CUH
#define COMPLEX_CUH

#include <cmath>
#include <math.h>

typedef struct complexNumber {
  double real;
  double imag;
} C;

__device__ inline double modulus(C *c) {
  return sqrtf((c->real * c->real) + (c->imag * c->imag));
}

__device__ __forceinline__ double modsq(C *c) {
  return (c->real * c->real) + (c->imag * c->imag);
}

__device__ __forceinline__ void add(C *z, C *cnst, C *res) {
  res->real = z->real + cnst->real;
  res->imag = z->imag + cnst->imag;
}

__device__ __forceinline__ void mult(C *x, C *y, C *res) {
  res->real = (x->real * y->real) - (x->imag * y->imag);
  res->imag = (x->real * y->imag) + (x->imag * y->real);
}

#endif
