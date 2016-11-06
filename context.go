package gofr

import (
	"image"
	"image/color"
	"log"
)

type Parameters struct {
	ImageWidth   int
	ImageHeight  int
	Min          complex128
	Max          complex128
	MaxI         int
	ColorFunc    ColorFunc
	Scaling      int
	Logger       *log.Logger
	MemberColor  color.NRGBA64
	EscapeRadius float64
}

type Context struct {
	Id         int
	Image      *image.NRGBA64
	Parameters *Parameters
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
				Id:         (i * n) + j,
				Image:      sub,
				Parameters: p,
			}

			c = append(c, &nc)
		}
	}

	return
}

func (self *Context) Delta() (dx, dy float64) {
	rw := self.Parameters.ImageWidth
	rh := self.Parameters.ImageHeight

	cw := real(self.Parameters.Max) - real(self.Parameters.Min)
	ch := imag(self.Parameters.Max) - imag(self.Parameters.Min)

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
 * Iterate over the map of pixel coordinates and complex points. Pass in
 */
func (self *Context) EachPoint(fn ContextFunc) {
	rmin := self.Image.Bounds().Min
	rmax := self.Image.Bounds().Max
	cmin := self.Parameters.Min
	dx, dy := self.Delta()
	var z complex128

	for x := rmin.X; x < rmax.X; x++ {
		for y := rmin.Y; y < rmax.Y; y++ {
			z = complex(real(cmin)+float64(x)*dx, imag(cmin)+float64(y)*dy)
			fn(x, y, z)
		}
	}
}
