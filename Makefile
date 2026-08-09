build-cuda:
	 nvcc --compiler-options '-fPIC' -std=c++23 -shared -o lib/libmandelbrot.so lib/mandelbrot.cu
