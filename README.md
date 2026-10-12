# Real-Time Mandelbrot and Julia Set Visualizer

An interactive fractal visualizer written in Go and CUDA (C++23). The application offloads fractal computation to custom GPU kernels and uses Go with SDL2 for windowing, user input, and real-time texture streaming.

README.md genenerated with AI, the rest of the project is hand written.

---

## At a Glance

- **Real-Time 1080p Rendering:** Evaluates escape-time fractals at 1920x1080 on the GPU with fluid frame rates during continuous panning and zooming.
- **Interactive Julia Set Exploration:** An in-window minimap maps complex coordinates directly to the parameter $c$, dynamically morphing Julia set structures as you drag the cursor.
- **Continuous Coloring:** Implements a continuous potential algorithm to remove discrete iteration banding, producing smooth palette transitions.
- **Go + CUDA Architecture:** Separates low-level SIMT compute into a C++23/CUDA shared library (`libmandelbrot.so`) invoked by a lightweight Go host application via Cgo.
- **Self-Contained UI:** Includes custom event dispatch, viewport projection math, and an in-house implementation of the midpoint (Bresenham) circle algorithm without external UI dependencies.

---

## Quickstart

### Prerequisites

- **OS:** Linux or Windows with WSL2
- **Hardware:** NVIDIA GPU with CUDA support
- **Toolchain:** 
  - CUDA Toolkit (`nvcc`)
  - Go 1.22+
  - GCC / Clang with C++23 support
- **Libraries & Fonts (Debian / Ubuntu):**
  ```bash
  sudo apt-get update
  sudo apt-get install -y libsdl2-dev libsdl2-ttf-dev fonts-ubuntu
  ```
  *(The application defaults to `/usr/share/fonts/truetype/ubuntu/UbuntuMono-R.ttf`)*

### Clone and Run

```bash
# 1. Clone the repository
git clone git@github.com:SocialSalt/mandelbrot.git
cd mandelbrot

# 2. Compile CUDA kernels into lib/libmandelbrot.so and launch the Go application
make run
```

To compile a standalone binary instead:
```bash
make build
./mandelbrot
```

### Controls

| Input | Action |
| :--- | :--- |
| **Left Click + Drag** | Pan viewport across the complex plane |
| **Scroll Wheel Up / Down** | Zoom in / Zoom out smoothly |
| **"Julia" Button (Top-Left)** | Toggle between Mandelbrot and Julia set rendering modes |
| **Minimap (Top-Left)** | Click and drag the indicator to change the Julia constant $c$ in real time |
| **Close Window** | Exit application |

---

## Architecture

The project splits responsibilities between high-level application control and low-level compute:

```
[ Go Host Application ]
  ├── SDL2 Event Loop (panning, zooming, window resize)
  ├── UI Components (toggle button, coordinate-mapped minimap)
  └── Memory Management (allocates host framebuffer, streams to SDL texture)
         │
         │ Cgo C-ABI Call
         ▼
[ CUDA Compute Engine (libmandelbrot.so) ]
  ├── Device Memory Allocation (cudaMalloc / cudaFree)
  ├── 2D SIMT Kernels (parallelMandelbrot, parallelJulia)
  ├── Continuous Coloring & Palette Interpolation
  └── Transfer Result back to Host (cudaMemcpy)
```

1. **Host Buffer Allocation:** Go allocates an RGB byte slice (`width * height * 3`).
2. **Compute Dispatch:** Go passes a pointer to the host buffer into `ComputeMandelbrot` or `ComputeJulia` through Cgo.
3. **GPU Execution:** The CUDA library allocates device memory, executes the 2D kernel grid, and copies the rendered pixels back to the host slice.
4. **Presentation:** Go streams the slice into an `SDL_TEXTUREACCESS_STREAMING` texture using `unsafe.SliceData`, renders UI overlays, and presents the frame.

---

## Implementation Details

### GPU Kernel & Grid Layout

The CUDA kernels (`lib/mandelbrot.cu` and `lib/julia.cu`) partition the screen using 32x32 thread blocks (1,024 threads per block):

- Grid dimensions: `ceil(width / 32)` by `ceil(height / 32)`.
- Each thread $(x, y)$ calculates its corresponding coordinate on the complex plane based on the viewport center and boundaries:
  - **Mandelbrot:** $z_0 = 0$, $c = \text{pixel coordinate}$
  - **Julia:** $z_0 = \text{pixel coordinate}$, $c = \text{user-selected constant}$
- Complex arithmetic primitives are marked `__device__ __forceinline__` to minimize overhead and maximize register reuse.

### Continuous Potential Coloring

Standard escape-time algorithms index colors using discrete integer iteration counts, creating visible step banding. The engine implements a continuous potential formula in `lib/include/color.cuh`:

$$\nu = \frac{\ln\left(\frac{\ln |z_n|}{\ln 2}\right)}{\ln 2}$$
$$i_{\text{smooth}} = i + 1 - \nu$$

Where $|z_n|$ is the escape radius modulus and $i$ is the iteration count. The fractional part of $i_{\text{smooth}}$ drives linear interpolation between adjacent colors in a 6-stop RGB palette, producing smooth gradients.

### Minimap Parameter Mapping & Circle Rasterizer

The minimap (`internal/userinterface/minimap.go`) defines a coordinate projection between screen pixels within a 200x200 UI rectangle and the complex domain $[-1, 1] \times [-1, 1]$. Mouse drag events dynamically update $c = c_{\text{real}} + i \cdot c_{\text{imag}}$ and trigger immediate recomputation of the Julia set.

To draw the selector handle without third-party graphics packages, the project implements the midpoint (Bresenham) circle algorithm in Go (`internal/userinterface/circle.go`), using 8-way symmetry to plot points directly on the SDL renderer.

---

## Engineering Trade-offs

### Host-to-Device Memory Roundtrips
Currently, each frame allocates device memory, executes the kernel, copies the framebuffer back to system RAM via `cudaMemcpy`, and uploads it to an SDL2 streaming texture.
- **Why this choice:** It cleanly decouples CUDA compilation from the GUI library, avoiding driver-specific OpenGL/Vulkan interop setup across different development environments.
- **Alternative:** Direct graphics interop (e.g., CUDA-OpenGL Pixel Buffer Objects) would allow the kernel to write directly into texture memory on the GPU, avoiding PCIe bus transfers entirely.

### Warp Divergence
On NVIDIA GPUs, threads in a 32-thread warp execute instructions in lockstep. Near fractal boundaries, threads inside the set execute the maximum iteration count (700 iterations), while adjacent outside threads escape early. Active threads stall while the slowest thread finishes.
- Standard mitigations include early checking for the main cardioid and period-2 bulb ($q = (x - 1/4)^2 + y^2; \ q(q + (x - 1/4)) < \frac{1}{4} y^2$) to immediately classify interior pixels without iterating.

### Cgo Dispatch Frequency
Calling C from Go introduces a small stack transition cost. Because this dispatch occurs once per frame for the entire 2-million-pixel buffer, the call overhead is negligible compared to the kernel execution time.

---

## Work in Progress: Extended Precision (Float128)

Standard 64-bit IEEE 754 floating-point numbers (`float64`) reach numerical precision limits at zoom levels around $10^{-14}$, causing pixelation and loss of detail.

Initial support for 128-bit precision is under development:
- `lib/mandelbrot128.cu`: Initial kernel using `__float128` via `<crt/device_fp128_functions.h>`.
- `internal/utils/utils.go`: Go utilities using `math/big` set to 113 bits of precision (`FP128_PRECISION = 113`), matching binary128 semantics.

Remaining work includes connecting the 128-bit kernel into the active Go event loop and comparing software-emulated quad-precision against perturbation rendering algorithms.
