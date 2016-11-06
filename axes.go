package gofr

import (
	"image/color"
	"math"
)

func DrawAxesInv(c *Context) {
	var x, y int
	context_min := c.Image.Bounds().Min
	context_max := c.Image.Bounds().Max

	zmin := c.Parameters.Min
	zmax := c.Parameters.Max

	zw := real(zmax) - real(zmin)
	zh := imag(zmax) - imag(zmin)

	dx := float64(c.Parameters.ImageWidth) / zw
	dy := float64(c.Parameters.ImageHeight) / zh

	x = int(math.Max(math.Abs(real(zmin))*dx, 0))
	y = int(math.Max(math.Abs(imag(zmin))*dy, 0.0))

	for ay := context_min.Y; ay < context_max.Y; ay++ {
		inv := c.Image.At(x, ay).(color.NRGBA64)
		c.Image.Set(x, ay, color.NRGBA64{^inv.R, ^inv.G, ^inv.B, 0xffff})
	}

	for ax := context_min.X; ax < context_max.X; ax++ {
		inv := c.Image.At(ax, y).(color.NRGBA64)
		c.Image.Set(ax, y, color.NRGBA64{^inv.R, ^inv.G, ^inv.B, 0xffff})
	}
}

func DrawAxesColor(c *Context, k color.NRGBA64) {
	var x, y int
	context_min := c.Image.Bounds().Min
	context_max := c.Image.Bounds().Max

	zmin := c.Parameters.Min
	zmax := c.Parameters.Max

	zw := real(zmax) - real(zmin)
	zh := imag(zmax) - imag(zmin)

	dx := float64(c.Parameters.ImageWidth) / zw
	dy := float64(c.Parameters.ImageHeight) / zh

	x = int(math.Max(math.Abs(real(zmin))*dx, 0))
	y = int(math.Max(math.Abs(imag(zmin))*dy, 0.0))

	for ay := context_min.Y; ay < context_max.Y; ay++ {
		c.Image.Set(x, ay, k)
	}

	for ax := context_min.X; ax < context_max.X; ax++ {
		c.Image.Set(ax, y, k)
	}
}

func DrawTicksColor(c *Context, tl int, unit float64, k color.NRGBA64) {
	var x, y, i int
	var n float64

	zmin := c.Parameters.Min
	zmax := c.Parameters.Max

	zw := real(zmax) - real(zmin)
	zh := imag(zmax) - imag(zmin)

	dx := float64(c.Parameters.ImageWidth) / zw
	dy := float64(c.Parameters.ImageHeight) / zh

	x = int(math.Max(math.Abs(real(zmin))*dx, 0))
	y = int(math.Max(math.Abs(imag(zmin))*dy, 0.0))

	for n = 0.0; n <= real(zmax); n += unit {
		tx := x + int(n*dx)
		for i = -1 * tl; i <= tl; i++ {
			ty := y + i
			c.Image.Set(tx, ty, k)
		}
	}

	for n = 0.0; n >= real(zmin); n -= unit {
		tx := x + int(n*dx)
		for i = -1 * tl; i <= tl; i++ {
			ty := y + i
			c.Image.Set(tx, ty, k)
		}
	}

	for n = 0.0; n <= imag(zmax); n += unit {
		ty := y + int(n*dy)
		for i = -1 * tl; i <= tl; i++ {
			tx := x + i
			c.Image.Set(tx, ty, k)
		}
	}

	for n = 0.0; n >= imag(zmin); n -= unit {
		ty := y + int(n*dy)
		for i = -1 * tl; i <= tl; i++ {
			tx := x + i
			c.Image.Set(tx, ty, k)
		}
	}
}

func DrawTicksInv(c *Context, tl int, unit float64) {
	var x, y, i int
	var n float64

	zmin := c.Parameters.Min
	zmax := c.Parameters.Max

	zw := real(zmax) - real(zmin)
	zh := imag(zmax) - imag(zmin)

	dx := float64(c.Parameters.ImageWidth) / zw
	dy := float64(c.Parameters.ImageHeight) / zh

	x = int(math.Max(math.Abs(real(zmin))*dx, 0))
	y = int(math.Max(math.Abs(imag(zmin))*dy, 0.0))

	for n = 0.0; n <= real(zmax); n += unit {
		tx := x + int(n*dx)
		for i = -1 * tl; i <= tl; i++ {
			ty := y + i

			inv := c.Image.At(tx, ty).(color.NRGBA64)
			c.Image.Set(tx, ty, color.NRGBA64{^inv.R, ^inv.G, ^inv.B, 0xffff})
		}
	}

	for n = 0.0; n >= real(zmin); n -= unit {
		tx := x + int(n*dx)
		for i = -1 * tl; i <= tl; i++ {
			ty := y + i

			inv := c.Image.At(tx, ty).(color.NRGBA64)
			c.Image.Set(tx, ty, color.NRGBA64{^inv.R, ^inv.G, ^inv.B, 0xffff})
		}
	}

	for n = 0.0; n <= imag(zmax); n += unit {
		ty := y + int(n*dy)
		for i = -1 * tl; i <= tl; i++ {
			tx := x + i

			inv := c.Image.At(tx, ty).(color.NRGBA64)
			c.Image.Set(tx, ty, color.NRGBA64{^inv.R, ^inv.G, ^inv.B, 0xffff})
		}
	}

	for n = 0.0; n >= imag(zmin); n -= unit {
		ty := y + int(n*dy)
		for i = -1 * tl; i <= tl; i++ {
			tx := x + i

			inv := c.Image.At(tx, ty).(color.NRGBA64)
			c.Image.Set(tx, ty, color.NRGBA64{^inv.R, ^inv.G, ^inv.B, 0xffff})
		}
	}
}
