#ifndef COLOR_CUH
#define COLOR_CUH

#define MAX_ITER 700
#define RADIUS 2.0e8f
#define NUM_COLORS 6
#include <cmath>

__device__ const double palette[NUM_COLORS][3] = {
    {235.0, 52.0, 52.0},  {247.0, 225.0, 25.0}, {17.0, 135.0, 11.0},
    {15.0, 235.0, 232.0}, {15.0, 15.0, 190.0},  {134.0, 13.0, 140.0},
};

__device__ inline void linear_interpolate(const double x[3], const double y[3],
                                          double t, unsigned char *r,
                                          unsigned char *g, unsigned char *b) {
  *r = (unsigned char)((1 - t) * x[0] + t * y[0]);
  *g = (unsigned char)((1 - t) * x[1] + t * y[1]);
  *b = (unsigned char)((1 - t) * x[2] + t * y[2]);
}

__device__ inline void getColor(int iters, double modz, unsigned char *r,
                                unsigned char *g, unsigned char *b) {
  double itn = 0;
  if (iters < MAX_ITER) {
    double log_zn = std::log(modz) / 2.0;
    double nu = std::log(log_zn / std::log(2.0)) / std::log(2.0);
    itn = (double)iters + 1.0 - nu;
  } else {
    itn = (double)iters;
  }
  const double *color1 = palette[(int)std::floor(itn) % NUM_COLORS];
  const double *color2 = palette[((int)std::floor(itn) + 1) % NUM_COLORS];
  double intpart;
  linear_interpolate(color1, color2, std::modf(itn, &intpart), r, g, b);
}

#endif
