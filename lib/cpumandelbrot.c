#include "include/mandelbrot.h"
#include <math.h>

#define MAX_ITER 100

typedef struct complexNumber {
  long double real;
  long double imag;
} C;

long double complexAbs(C *c) {
  return sqrtl((c->real * c->real) + (c->imag * c->imag));
}

void complexAdd(C *z, C *cnst, C *res) {
  res->real = z->real + cnst->real;
  res->imag = z->imag + cnst->imag;
}

void complexMult(C *x, C *y, C *res) {
  res->real = (x->real * y->real) - (x->imag * y->imag);
  res->imag = (x->real * y->imag) + (x->imag * y->real);
}

int mandelbrot(C *c) {
  C z = {0.0, 0.0};
  C zsq;
  for (int i = 0; i < MAX_ITER; i++) {
    if (complexAbs(&z) > 2) {
      return i;
    }
    complexMult(&z, &z, &zsq);
    complexAdd(&zsq, c, &z);
  }
  return MAX_ITER;
}

void getColor(int itrs, unsigned char *r, unsigned char *g, unsigned char *b) {
  *r = (unsigned char)(itrs * 2.0f);
  *g = (unsigned char)(itrs * 1.9f);
  *b = (unsigned char)(itrs * 2.35f);
}

void computeMandelbrot(unsigned char *image, double real_min, double imag_min,
                       double inc_real, double inc_imag, int img_w, int img_h,
                       int channels) {

  for (int x = 0; x < img_w; x++) {
    for (int y = 0; y < img_h; y++) {
      double real = real_min + ((double)x * inc_real);
      double imag = imag_min + ((double)y * inc_imag);
      C c = {real, imag};
      int iters = mandelbrot(&c);
      float ratio = (iters * 1.0f) / MAX_ITER;
      unsigned char r, g, b;

      getColor(iters, &r, &g, &b);

      int pixel_index = (x + (y * img_w)) * channels;
      image[pixel_index + 0] = r;
      image[pixel_index + 1] = g;
      image[pixel_index + 2] = b;
    }
  }
}

void CpuMandelbrot(unsigned char *res, double real_min, double real_max,
                   double imag_min, double imag_max, int img_w, int img_h,
                   int channels) {

  long double inc_real = (real_max - real_min) / img_w;
  long double inc_imag = (imag_max - imag_min) / img_h;
  int bytes = img_w * img_h * channels * sizeof(unsigned char);

  computeMandelbrot(res, real_min, imag_min, inc_real, inc_imag, img_w, img_h,
                    channels);
}
