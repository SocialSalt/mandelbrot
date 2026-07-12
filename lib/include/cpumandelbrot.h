#ifndef MANDELBROT_H
#define MANDELBROT_H

void CpuMandelbrot(unsigned char *res, double real_min, double real_max,
                   double imag_min, double imag_max, int img_w, int img_h,
                   int channels);

#endif
