package gofr

import (
	"math"
)

func Mandelbrot(c *Context) int {
	max_i := c.Parameters.MaxI
	fn := func(x, y int, z complex128) {
		i, zn := Escape(c, z, max_i)
		c.Parameters.ColorFunc(c, zn, x, y, i, max_i)
	}
	c.EachPoint(fn)
	return 0
}

func Escape(c *Context, z complex128, max_i int) (int, complex128) {
	var d float64
	var i int

	z0 := z
	zn := complex(0, 0)

	for {
		z = z*z + z0
		if zn == z {
			return max_i, z
		}
		zn = z

		d = math.Sqrt(real(z)*real(z) + imag(z)*imag(z))
		if d >= c.Parameters.EscapeRadius || i == max_i {
			return i, z
		}

		i++
	}

	return i, z0
}
