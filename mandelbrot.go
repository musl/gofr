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
	z0 := z
	var zl complex128
	var i int

	if CheckCardioid(z) {
		return max_i, z
	}

	d := math.Sqrt(real(z)*real(z) + imag(z)*imag(z))
	for i = 0; d < c.Parameters.EscapeRadius && i < max_i; i++ {
		z = z*z + z0
		d = math.Sqrt(real(z)*real(z) + imag(z)*imag(z))

		// Periodicity Check
		if z == zl {
			return max_i, z
		}
		zl = z
	}

	return i, z
}

func CheckCardioid(z complex128) bool {
	q := (real(z)-0.25)*(real(z)-0.25) + imag(z)*imag(z)
	q = q * (q + (real(z) - 0.25))
	return q < imag(z)*imag(z)*0.25
}
