package gofr

import (
	"math"
	"math/cmplx"
)

func Experimental(c *Context, cancel chan bool) int {
	max_i := c.MaxI
	fn := func(x, y int, z complex128) {
		i, zn := ExperimentalEscape(c, z, max_i)
		c.ColorFunc(c, zn, x, y, i, max_i)
	}
	c.EachPoint(fn, cancel)
	return 0
}

func ExperimentalEscape(c *Context, z complex128, max_i int) (int, complex128) {
	i := 0
	z0 := z
	zn := complex(0, 0)
	p := c.Power
	//e := complex(math.E, 0.0)

	if p <= 0 {
		p = 2
	}

	for {

		t := z
		for j := 0; j < p-1; j++ {
			z = z * t
		}
		z += 1 / cmplx.Log(-z0)

		if zn == z {
			return max_i, z
		}
		zn = z

		d := math.Sqrt(real(z)*real(z) + imag(z)*imag(z))
		if d >= c.EscapeRadius || i == max_i {
			return i, z
		}

		i++
	}

	return i, z0
}
