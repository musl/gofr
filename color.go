package gofr

import (
	"image/color"
	"math"
)

type ColorFunc func(*Context, complex128, int, int, int, int)

func ColorSmooth(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.Set(x, y, c.Parameters.MemberColor)
		return
	}

	log_zn := math.Log(real(z)*real(z)+imag(z)*imag(z)) / 2.0
	nu := math.Log(log_zn/math.Log(2.0)) / math.Log(2.0)
	j := float64(i) + 1.0 - nu

	o := math.Pi
	f := math.Pi / 32.0
	t := f * math.Pi * float64(j)
	r := uint16(0x7fff + 0x7fff*math.Sin(o+t))
	g := uint16(0x7fff + 0x7fff*math.Sin(o+0.25*math.Pi+t))
	b := uint16(0x7fff + 0x7fff*math.Cos(o+t))

	l := color.NRGBA64{r, g, b, 0xffff}
	c.Image.Set(x, y, l)
}

func ColorBands(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.Set(x, y, c.Parameters.MemberColor)
		return
	}

	o := math.Pi
	f := float64(max_i) / 16.0
	t := f * math.Pi * (float64(i) / float64(max_i))
	r := uint16(0x7fff + 0x7fff*math.Sin(o+t))
	g := uint16(0x7fff + 0x7fff*math.Sin(o+0.25*math.Pi+t))
	b := uint16(0x7fff + 0x7fff*math.Cos(o+t))

	l := color.NRGBA64{r, g, b, 0xffff}
	c.Image.Set(x, y, l)
}

func ColorMono(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.Set(x, y, c.Parameters.MemberColor)
		return
	}

	k := uint16(0xffff * (i & 1))
	l := color.NRGBA64{k, k, k, 0xffff}
	c.Image.Set(x, y, l)
}

func ColorMonoStripe(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.Set(x, y, c.Parameters.MemberColor)
		return
	}

	if (i-1)%9 == 0 {
		c.Image.Set(x, y, color.NRGBA64{0, 0xa000, 0xc000, 0xffff})
		return
	}

	k := uint16(0xffff * (i & 1))
	l := color.NRGBA64{k, k, k, 0xffff}
	c.Image.Set(x, y, l)
}
