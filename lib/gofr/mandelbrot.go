package gofr

import "math"

func Mandelbrot(c *Context, cancel chan bool) int {
	maxI := c.MaxI
	fn := func(x, y int, z complex128) {
		i, zn := Escape(c, z, maxI)
		c.ColorFunc(c, zn, x, y, i, maxI)
	}
	c.EachPoint(fn, cancel)
	return 0
}

func Escape(c *Context, z complex128, maxI int) (int, complex128) {
	i := 0
	z0 := z
	zn := complex(0, 0)
	p := c.Power

	if p <= 0 {
		p = 2
	}

	for {
		// inflexible!
		//z = z*z + z0
		// slow!
		//z = cmplx.Pow(z, c.Power) + z0

		t := z
		for j := 0; j < p-1; j++ {
			z = z * t
		}
		z += z0

		if zn == z {
			return maxI, z
		}
		zn = z

		d := math.Sqrt(real(z)*real(z) + imag(z)*imag(z))
		if d >= c.EscapeRadius || i == maxI {
			return i, z
		}

		i++
	}

	return i, z0
}
