package gofr

import (
	"math"
	"math/cmplx"
)

func Ebrot(c *Context, cancel chan bool) int {
	maxI := c.MaxI
	fn := func(x, y int, z complex128) {
		i, zn := EBrotEscape(c, z, maxI)
		c.ColorFunc(c, zn, x, y, i, maxI)
	}
	c.EachPoint(fn, cancel)
	return 0
}

func EBrotEscape(c *Context, z complex128, maxI int) (int, complex128) {
	i := 0
	z0 := z
	zn := complex(0, 0)
	e := complex(math.E, math.E)
	//p := complex(float64(c.Power), 0)

	for {
		z = cmplx.Pow(z, e) + z0

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
