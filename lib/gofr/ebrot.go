package gofr

import (
	"math"
	"math/cmplx"
)

func Ebrot(c *Context, cancel chan bool) int {
	max_i := c.MaxI
	fn := func(x, y int, z complex128) {
		i, zn := EBrotEscape(c, z, max_i)
		c.ColorFunc(c, zn, x, y, i, max_i)
	}
	c.EachPoint(fn, cancel)
	return 0
}

func EBrotEscape(c *Context, z complex128, max_i int) (int, complex128) {
	i := 0
	z0 := z
	zn := complex(0, 0)
	e := complex(math.E, 0)
	p := c.Power
	dx, _ := c.Delta()

	if p <= 0 {
		p = 2
	}

	if Round(real(z0), dx) == 0.5 {
		return 0, complex(0, 0)
	}

	for {

		t := z
		for j := 0; j < p-1; j++ {
			z = z * t
		}
		z += z0
		z = cmplx.Pow(e, z)

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
