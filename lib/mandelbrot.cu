#include "include/color.cuh"
#include "include/complex.cuh"
#include "include/mandelbrot.h"
#include <stdlib.h>

#define NUM_THREADS 32 * 32
#define NUM_BLOCKS (int)ceil(((long double)IMG_W * IMG_H) / NUM_THREADS)

__device__ int mandelbrot(C *z, C *c) {
  C zsq;
  for (int i = 0; i < MAX_ITER; i++) {
    if (modulus(z) > 2) {
      return i;
    }
    mult(z, z, &zsq);
    add(&zsq, c, z);
  }
  return MAX_ITER;
}

__global__ void parallelMandelbrot(unsigned char *image, double real_min,
                                   double real_scale, double imag_min,
                                   double imag_scale, int img_w, int img_h,
                                   int channels) {
  int x = (blockIdx.x * blockDim.x) + threadIdx.x;
  int y = (blockIdx.y * blockDim.y) + threadIdx.y;
  if (x < img_w && y < img_h) {
    double real = real_min + ((double)x * real_scale);
    double imag = imag_min + ((double)y * imag_scale);
    C c = {real, imag};
    C z = {0.0, 0.0};

    int iters = mandelbrot(&z, &c);

    unsigned char r, g, b;
    getColor(iters, modulus(&z), 2.0, &r, &g, &b);

    int pixel_index = (x + (y * img_w)) * channels;
    image[pixel_index + 0] = r;
    image[pixel_index + 1] = g;
    image[pixel_index + 2] = b;
  }
}

extern "C" {
void ComputeMandelbrot(unsigned char *image, int img_w, int img_h,
                       double real_center, double imag_center,
                       double real_width, double imag_height, int channels) {

  double real_max = real_center + real_width / 2.0;
  double real_min = real_center - real_width / 2.0;
  double imag_max = imag_center + imag_height / 2.0;
  double imag_min = imag_center - imag_height / 2.0;

  double real_scale = (real_max - real_min) / img_w;
  double imag_scale = (imag_max - imag_min) / img_h;
  int bytes = img_w * img_h * channels * sizeof(unsigned char);

  // printf("%f + %fi, W: %f H: %f\n", (double)real_center, (double)imag_center,
  //        (double)real_width, (double)imag_height);

  unsigned char *dev_image;
  cudaMalloc(&dev_image, img_w * img_h * channels);

  dim3 block_dim(32, 32, 1);
  dim3 grid_dim(ceil((float)img_w / 32), ceil((float)img_h / 32), 1);

  parallelMandelbrot<<<grid_dim, block_dim>>>(dev_image, real_min, real_scale,
                                              imag_min, imag_scale, img_w,
                                              img_h, channels);
  cudaDeviceSynchronize();
  cudaMemcpy(image, dev_image, bytes, cudaMemcpyDeviceToHost);
  cudaFree(dev_image);
}
}
