#include <math.h>
#include <stdlib.h>

#define MAX_ITER 100
#define NUM_THREADS 32 * 32
#define NUM_BLOCKS (int)ceil(((long double)IMG_W * IMG_H) / NUM_THREADS)

typedef struct complexNumber {
  long double real;
  long double imag;
} C;

__device__ long double complexAbs(C *c) {
  return sqrtf((c->real * c->real) + (c->imag * c->imag));
}

__device__ void complexAdd(C *z, C *cnst, C *res) {
  res->real = z->real + cnst->real;
  res->imag = z->imag + cnst->imag;
}

__device__ void complexMult(C *x, C *y, C *res) {
  res->real = (x->real * y->real) - (x->imag * y->imag);
  res->imag = (x->real * y->imag) + (x->imag * y->real);
}

__device__ int mandelbrot(C *c) {
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

__device__ void getColor(int itrs, unsigned char *r, unsigned char *g,
                         unsigned char *b) {
  *r = (unsigned char)(itrs * 2.0f);
  *g = (unsigned char)(itrs * 1.9f);
  *b = (unsigned char)(itrs * 2.35f);
}

__global__ void parallelMandelbrot(unsigned char *dev_image,
                                   long double real_min, long double inc_real,
                                   long double imag_min, long double inc_imag,
                                   int img_w, int img_h, int channels) {
  int x = (blockIdx.x * blockDim.x) + threadIdx.x;
  int y = (blockIdx.y * blockDim.y) + threadIdx.y;
  if (x < img_w && y < img_h) {
    long double real = real_min + ((long double)x * inc_real);
    long double imag = imag_min + ((long double)y * inc_imag);
    C c = {real, imag};
    int iters = mandelbrot(&c);
    unsigned char r, g, b;

    getColor(iters, &r, &g, &b);

    int pixel_index = (x + (y * img_w)) * channels;
    dev_image[pixel_index + 0] = r;
    dev_image[pixel_index + 1] = g;
    dev_image[pixel_index + 2] = b;
  }
}

extern "C" void LaunchMandelbrot(unsigned char *res, long double real_min,
                                 long double real_max, long double imag_min,
                                 long double imag_max, int img_w, int img_h,
                                 int channels) {

  long double inc_real = (real_max - real_min) / img_w;
  long double inc_imag = (imag_max - imag_min) / img_h;
  int bytes = img_w * img_h * channels * sizeof(unsigned char);

  unsigned char *dev_image;
  cudaMalloc(&dev_image, img_w * img_h * channels);

  dim3 block_dim(32, 32, 1);
  dim3 grid_dim(ceil((float)img_w / 32), ceil((float)img_h / 32), 1);

  parallelMandelbrot<<<grid_dim, block_dim>>>(dev_image, real_min, inc_real,
                                              imag_min, inc_imag, img_w, img_h,
                                              channels);
  cudaDeviceSynchronize();
  cudaMemcpy(res, dev_image, bytes, cudaMemcpyDeviceToHost);
  cudaFree(dev_image);
}
