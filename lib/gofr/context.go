package gofr

import (
	"image"
	"image/color"
)

/*
 * Easily serializeable parameters for rendering images.
 */
type Parameters struct {
	ColorFunc    string
	EscapeRadius float64
	ImageHeight  int
	ImageWidth   int
	Max          complex128
	MaxI         int
	MemberColor  string
	Min          complex128
	Scaling      int
	Power        int
}

/*
 * Context for rendering images in parallel - a basic description of
 * work for a render job to be run in one thread.
 */
type Context struct {
	ColorFunc    ColorFunc
	EscapeRadius float64
	Id           int
	Image        *image.NRGBA64
	ImageHeight  int
	ImageWidth   int
	Max          complex128
	MaxI         int
	MemberColor  color.NRGBA64
	Min          complex128
	Scaling      int
	Power        int
}

func MakeContexts(im *image.NRGBA64, n int, p *Parameters) (c []*Context) {
	r := im.Bounds()
	x := 0
	y := 0
	dx := r.Max.X / n
	dy := r.Max.Y / n
	rx := dx + (r.Max.X % n)
	ry := dy + (r.Max.Y % n)
	w := dx
	h := dy

	mc, err := MemberColorFromString(p.MemberColor)
	if err != nil {
		panic(err)
	}

	cf, err := ColorFuncFromString(p.ColorFunc)
	if err != nil {
		panic(err)
	}

	if n <= 0 {
		panic("I refuse to make zero or fewer contexts of an image.")
	}

	if n > r.Max.X || n > r.Max.Y {
		panic("I refuse to make more contexts than I have pixels.")
	}

	for i := 0; i < n; i++ {
		x = i * dx
		if i == n-1 {
			w = rx
		}

		for j := 0; j < n; j++ {
			y = j * dy
			if j == n-1 {
				h = ry
			}

			sub := im.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA64)
			nc := Context{
				ColorFunc:    cf,
				EscapeRadius: p.EscapeRadius,
				Id:           (i * n) + j,
				Image:        sub,
				ImageHeight:  p.ImageHeight,
				ImageWidth:   p.ImageWidth,
				Max:          p.Max,
				MaxI:         p.MaxI,
				MemberColor:  mc,
				Min:          p.Min,
				Scaling:      p.Scaling,
				Power:        p.Power,
			}

			c = append(c, &nc)
		}
	}

	return
}

func (self *Context) delta() (dx, dy float64) {
	rw := self.ImageWidth
	rh := self.ImageHeight

	cw := real(self.Max) - real(self.Min)
	ch := imag(self.Max) - imag(self.Min)

	dx = cw / float64(rw)
	dy = ch / float64(rh)
	return
}

/*
* Use this with EachPoint to iterate over the map of pixel coordinates
* and mapped complex points.
 */
type ContextFunc func(int, int, complex128)

/*
* Iterate over the map of pixel coordinates and complex points.
 */
func (self *Context) EachPoint(fn ContextFunc) {
	rmin := self.Image.Bounds().Min
	rmax := self.Image.Bounds().Max
	cmin := self.Min
	dx, dy := self.delta()
	var z complex128

	for x := rmin.X; x < rmax.X; x++ {
		for y := rmin.Y; y < rmax.Y; y++ {
			z = complex(real(cmin)+float64(x)*dx, imag(cmin)+float64(y)*dy)
			fn(x, y, z)
		}
	}
}
