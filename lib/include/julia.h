// julia.h
#ifndef JULIA_H
#define JULIA_H

#ifdef __cplusplus
extern "C" {
#endif

void ComputeJulia(unsigned char *image, int img_w, int img_h,
                  double real_center, double imag_center, double real_width,
                  double imag_height, int channels, double c_real,
                  double c_imag);

#ifdef __cplusplus
}
#endif

#endif
