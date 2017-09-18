package gofr

import (
	"image"
	"math/rand"
	"regexp"
	"runtime"
	"testing"
)

var n_cpu = runtime.NumCPU()

func parameters() Parameters {
	a := complex(-2.6, -2.1)
	b := complex(1.6, 2.1)
	w := 512
	h := int(float64(w) * (imag(b) - imag(a)) / (real(b) - real(a)))

	return Parameters{
		ImageWidth:   w,
		ImageHeight:  h,
		Min:          a,
		Max:          b,
		MaxI:         1000,
		ColorFunc:    "mono",
		RenderFunc:   "mandelbrot",
		EscapeRadius: 4.0,
		MemberColor:  "#000000",
		Power:        2,
	}
}

func contexts(p *Parameters) []*Context {
	img := image.NewNRGBA64(image.Rect(0, 0, p.ImageWidth, p.ImageHeight))
	return MakeContexts(img, n_cpu, p)
}

func TestVersion(t *testing.T) {
	pattern := regexp.MustCompile(`\d+\.\d+\.\d+`)
	if !pattern.MatchString(Version) {
		t.Errorf("Expected %#v to match \"%s\".", Version, pattern.String())
	}
}

func TestRenderImage(t *testing.T) {
	c := make(chan bool)
	p := parameters()
	contexts := contexts(&p)
	Render(n_cpu, contexts, c)
}

func TestMandelbrot(t *testing.T) {
	c := make(chan bool)
	p := parameters()
	contexts := contexts(&p)

	rc := Mandelbrot(contexts[0], c)
	if rc != 0 {
		t.Errorf("Mandelbrot didn't return 0: %d", rc)
	}
}

func TestEscape(t *testing.T) {
	p := parameters()
	contexts := contexts(&p)
	z_in := complex(0.1*rand.Float64(), 0.1*rand.Float64())
	z_out := complex(2.0*rand.Float64(), 2.0*rand.Float64())

	i, _ := Escape(contexts[0], z_in, p.MaxI)
	if i != p.MaxI {
		t.Errorf("Incorrectly calculated point in set: %v iterations: %v != %v", z_in, i, p.MaxI)
	}
	i, _ = Escape(contexts[0], z_out, p.MaxI)
	if i >= p.MaxI {
		t.Errorf("Incorrectly calculated point not in set: %v iterations: %v >= %v", z_out, i, p.MaxI)
	}
}

func BenchmarkRenderImage(b *testing.B) {
	c := make(chan bool)
	p := parameters()
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		Render(n_cpu, contexts, c)
	}
}

func BenchmarkMandelbrot(b *testing.B) {
	c := make(chan bool)
	p := parameters()
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		Mandelbrot(contexts[0], c)
	}
}

func BenchmarkColorMono(b *testing.B) {
	p := parameters()
	p.ColorFunc = "mono"
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		ColorMono(contexts[0], complex(0, 0), 0, 0, 0, p.MaxI)
	}
}

func BenchmarkColorMonoStripe(b *testing.B) {
	p := parameters()
	p.ColorFunc = "stripe"
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		ColorMono(contexts[0], complex(0, 0), 0, 0, 0, p.MaxI)
	}
}

func BenchmarkColorBands(b *testing.B) {
	p := parameters()
	p.ColorFunc = "bands"
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		ColorMono(contexts[0], complex(0, 0), 0, 0, 0, p.MaxI)
	}
}

func BenchmarkColorSmooth(b *testing.B) {
	p := parameters()
	p.ColorFunc = "smooth"
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		ColorMono(contexts[0], complex(0, 0), 0, 0, 0, p.MaxI)
	}
}

func BenchmarkEscapeIn(b *testing.B) {
	p := parameters()
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		z := complex(0.1*rand.Float64(), 0.1*rand.Float64())
		Escape(contexts[0], z, p.MaxI)
	}
}

func BenchmarkEsceapeOut(b *testing.B) {
	p := parameters()
	contexts := contexts(&p)

	for i := 0; i < b.N; i++ {
		z := complex(2.0*rand.Float64(), 2.0*rand.Float64())
		Escape(contexts[0], z, p.MaxI)
	}
}
