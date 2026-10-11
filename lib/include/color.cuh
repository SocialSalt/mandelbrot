#ifndef COLOR_CUH
#define COLOR_CUH

#define MAX_ITER 1000
#define MAX_ITERF 1000.0f
#include <cmath>

__device__ inline void getColor(int itrs, double modz, double rad,
                                unsigned char *r, unsigned char *g,
                                unsigned char *b) {

  if (itrs == MAX_ITER) {
    *r = 0;
    *g = 0;
    *b = 0;
  } else {
    double v = MAX_ITERF - (std::log(modz) / std::log(rad));
    *r = (unsigned char)(itrs * (200.0f / MAX_ITERF));
    *g = (unsigned char)(itrs * (190.0f / MAX_ITERF));
    *b = (unsigned char)(itrs * (235.0f / MAX_ITERF));
  }
}

#endif
