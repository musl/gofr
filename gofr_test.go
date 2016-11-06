package gofr

import (
	"flag"
	"image"
	"math/rand"
	"os"
	"runtime"
	"testing"
)

var a = complex(-2.6, -2.1)
var b = complex(1.6, 2.1)
var w = 1024
var h = int(float64(w) * (imag(b) - imag(a)) / (real(b) - real(a)))

var p = Parameters{
	ImageWidth:  w,
	ImageHeight: h,
	Min:         a,
	Max:         b,
	MaxI:        1000,
	ColorFunc:   ColorMono,
}

var img = image.NewNRGBA64(image.Rect(0, 0, p.ImageWidth, p.ImageHeight))
var n_cpu = runtime.NumCPU()
var contexts = MakeContexts(img, n_cpu, &p)
var context = contexts[0]
var z_in = complex(0.01, 0.01)
var z_out = complex(-2.0, -2.0)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestRenderImage(t *testing.T) {
	Render(n_cpu, contexts, Mandelbrot)
}

func TestMandelbrot(t *testing.T) {
	rc := Mandelbrot(context)
	if rc != 0 {
		t.Errorf("Mandelbrot didn't return 0: %d", rc)
	}
}

func TestEscape(t *testing.T) {
	i, _ := Escape(context, z_in, p.MaxI)
	if i != p.MaxI {
		t.Errorf("Incorrectly calculated point in set: %v iterations: %v", z_in, p.MaxI)
	}
	i, _ = Escape(context, z_out, p.MaxI)
	if i >= p.MaxI {
		t.Errorf("Incorrectly calculated point not in set: %v iterations: %v", z_out, p.MaxI)
	}
}

func BenchmarkRenderImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render(1, contexts, Mandelbrot)
	}
}

func BenchmarkMandelbrot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Mandelbrot(context)
	}
}

func BenchmarkColorMono(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ColorMono(context, complex(0, 0), 0, 0, 0, p.MaxI)
	}
}

func BenchmarkEscapeIn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		z := complex(0.1*rand.Float64(), 0.1*rand.Float64())
		Escape(context, z, p.MaxI)
	}
}

func BenchmarkEsceapeOut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		z := complex(2.0*rand.Float64(), 2.0*rand.Float64())
		Escape(context, z, p.MaxI)
	}
}
