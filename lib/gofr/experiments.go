package gofr

import (
	"math"
	"math/cmplx"
)

// Experimental is
func Experimental(c *Context, cancel chan bool) int {
	maxI := c.MaxI
	fn := func(x, y int, z complex128) {
		i, zn := Escape(c, z, maxI)
		c.ColorFunc(c, zn, x, y, i, maxI)
	}
	c.EachPoint(fn, cancel)
	return 0
}

// ExperimentalEscape is
func ExperimentalEscape(c *Context, z complex128, maxI int) (int, complex128) {
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

		// Rotate z about the complex origin.
		r, theta := cmplx.Polar(z)
		theta += 0.25 * math.Pi
		z = cmplx.Rect(r, theta)

		z += z0

		if zn == z {
			return maxI, z
		}
		zn = z

		// Benchmark Polar() vs doing the math ourselves.
		d := math.Sqrt(real(z)*real(z) + imag(z)*imag(z))
		if d >= c.EscapeRadius || i == maxI {
			return i, z
		}

		i++
	}

	return i, z0
}
