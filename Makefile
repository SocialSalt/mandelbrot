build-cuda:
	 nvcc --compiler-options '-fPIC' -arch=sm_100 -std=c++23 -shared -o lib/libmandelbrot.so lib/mandelbrot.cu -Xlinker -lquadmath
