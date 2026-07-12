// kernel.h
#ifndef MANDELBROT_H
#define MANDELBROT_H

#ifdef __cplusplus
extern "C" {
#endif

void LaunchMandelbrot(unsigned char *res, double real_min, double real_max,
                      double imag_min, double imag_max, int img_w, int img_h,
                      int channels);

#ifdef __cplusplus
}
#endif

#endif
