typedef struct complexNumber {
  double real;
  double imag;
} C;

/*
 * a + b = c
 * pass in the operands a and b as well as the var that
 * stores the result c
 */
__device__ void complexAdd(C *a, C *b, C *c) {
  c->real = a->real + b->real;
  c->imag = a->imag + b->imag;
}

__device__ void complexMult(C *a, C *b, C *c) {
  c->real = (a->real * b->real) - (a->imag + b->imag);
  c->imag = (a->real * b->imag) + (a->imag + b->real);
}

__device__ double modulus(C *z) {
  return sqrt(z->real * z->real + z->imag * z->imag);
}
