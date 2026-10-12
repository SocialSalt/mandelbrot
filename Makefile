build-cuda:
	 nvcc --compiler-options '-fPIC' -std=c++23 -shared -o lib/libmandelbrot.so lib/mandelbrot.cu lib/julia.cu

run:
	 nvcc --compiler-options '-fPIC' -std=c++23 -shared -o lib/libmandelbrot.so lib/mandelbrot.cu lib/julia.cu
	go run ./cmd/mandelbrot.go

build:
	 nvcc --compiler-options '-fPIC' -std=c++23 -shared -o lib/libmandelbrot.so lib/mandelbrot.cu lib/julia.cu
	go build ./cmd/mandelbrot.go

