// kernel.h
#ifndef MANDELBROT_H
#define MANDELBROT_H

#ifdef __cplusplus
extern "C" {
#endif

void ComputeMandelbrot(unsigned char *image, int img_w, int img_h,
                       double real_center, double imag_center,
                       double real_width, double imag_height, int channels);

#ifdef __cplusplus
}
#endif

#endif
