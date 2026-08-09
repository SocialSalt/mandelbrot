#include "include/mandelbrot.h"
#include <crt/device_fp128_functions.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>

#define MAX_ITER 1000
#define NUM_THREADS 32 * 32
#define NUM_BLOCKS (int)ceil(((long double)IMG_W * IMG_H) / NUM_THREADS)

// fp 128
typedef struct complex128 {
  __float128 real;
  __float128 imag;
} C128;

__device__ __float128 modsq(C128 *c) {
  return (c->real * c->real) + (c->imag * c->imag);
}

__device__ void add(C128 *z, C128 *cnst, C128 *res) {
  res->real = z->real + cnst->real;
  res->imag = z->imag + cnst->imag;
}

__device__ void mult(C128 *x, C128 *y, C128 *res) {
  res->real = (x->real * y->real) - (x->imag * y->imag);
  res->imag = (x->real * y->imag) + (x->imag * y->real);
}

__device__ int mandelbrot(C128 *c) {
  C128 z = {0.0, 0.0};
  C128 zsq;
  for (int i = 0; i < MAX_ITER; i++) {
    if (modsq(&z) > 4.0) {
      return i;
    }
    mult(&z, &z, &zsq);
    add(&zsq, c, &z);
  }
  return MAX_ITER;
}

__device__ void getColor128(int itrs, unsigned char *r, unsigned char *g,
                            unsigned char *b) {

  if (itrs == MAX_ITER) {
    *r = 0;
    *g = 0;
    *b = 0;
  } else {
    *r = (unsigned char)(itrs * 2.0f);
    *g = (unsigned char)(itrs * 1.9f);
    *b = (unsigned char)(itrs * 2.35f);
  }
}

__global__ void parallelMandelbrot128(unsigned char *image, __float128 real_min,
                                      __float128 real_scale,
                                      __float128 imag_min,
                                      __float128 imag_scale, int img_w,
                                      int img_h, int channels) {
  int x = (blockIdx.x * blockDim.x) + threadIdx.x;
  int y = (blockIdx.y * blockDim.y) + threadIdx.y;
  if (x < img_w && y < img_h) {
    printf("%d, %d", x, y);
    __float128 real = real_min + ((__float128)x * real_scale);
    __float128 imag = imag_min + ((__float128)y * imag_scale);
    C128 c = {real, imag};

    int iters = mandelbrot(&c);

    unsigned char r, g, b;
    getColor128(iters, &r, &g, &b);

    int pixel_index = (x + (y * img_w)) * channels;
    image[pixel_index + 0] = r;
    image[pixel_index + 1] = g;
    image[pixel_index + 2] = b;
  }
}
__float128 parseFloat128(const char *s) { return strtof128(s, NULL); }

void computeMandelbrot128(unsigned char *image, int img_w, int img_h,
                          __float128 real_center, __float128 imag_center,
                          __float128 real_width, __float128 imag_height,
                          int channels) {

  __float128 real_max = real_center + real_width / 2.0;
  __float128 real_min = real_center - real_width / 2.0;
  __float128 imag_max = imag_center + imag_height / 2.0;
  __float128 imag_min = imag_center - imag_height / 2.0;

  // printf("%f + %fi, W: %f H: %f\n", (double)real_center, (double)imag_center,
  //        (double)real_width, (double)imag_height);

  __float128 real_scale = (real_max - real_min) / img_w;
  __float128 imag_scale = (imag_max - imag_min) / img_h;
  int bytes = img_w * img_h * channels * sizeof(unsigned char);

  unsigned char *dev_image;
  cudaMalloc(&dev_image, img_w * img_h * channels);

  dim3 block_dim(32, 32, 1);
  dim3 grid_dim(ceil((float)img_w / 32), ceil((float)img_h / 32), 1);

  parallelMandelbrot128<<<grid_dim, block_dim>>>(
      dev_image, real_min, real_scale, imag_min, imag_scale, img_w, img_h,
      channels);
  cudaDeviceSynchronize();
  cudaMemcpy(image, dev_image, bytes, cudaMemcpyDeviceToHost);
  cudaFree(dev_image);
}

extern "C" {
void ComputeMandelbrot128Double(unsigned char *image, int img_w, int img_h,
                                double real_center, double imag_center,
                                double real_width, double imag_height,
                                int channels) {
  computeMandelbrot128(
      image, img_w, img_h, static_cast<__float128>(real_center),
      static_cast<__float128>(imag_center), static_cast<__float128>(real_width),
      static_cast<__float128>(imag_height), channels);
}

void ComputeMandelbrot128String(unsigned char *image, int img_w, int img_h,
                                const char *real_center,
                                const char *imag_center, const char *real_width,
                                const char *imag_height, int channels

) {
  computeMandelbrot128(image, img_w, img_h, parseFloat128(real_center),
                       parseFloat128(imag_center), parseFloat128(real_width),
                       parseFloat128(imag_height), channels);
}
}
