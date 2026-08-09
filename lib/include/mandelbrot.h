// kernel.h
#ifndef MANDELBROT_H
#define MANDELBROT_H

#ifdef __cplusplus
extern "C" {
#endif

void ComputeMandelbrot(unsigned char *image, int img_w, int img_h,
                       double real_center, double imag_center,
                       double real_width, double imag_height, int channels);

void ComputeMandelbrot128Double(unsigned char *image, int img_w, int img_h,
                                double real_center, double imag_center,
                                double real_width, double imag_height,
                                int channels);

void ComputeMandelbrot128String(unsigned char *image, int img_w, int img_h,
                                const char *real_center,
                                const char *imag_center, const char *real_width,
                                const char *imag_height, int channels);

#ifdef __cplusplus
}
#endif

#endif
