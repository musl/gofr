package gofr

import (
	"errors"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
	"math/cmplx"
	"strconv"
)

var (
	Accent  = color.NRGBA64{0, 0xa000, 0xc000, 0xffff}
	White   = color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	Black   = color.NRGBA64{0, 0, 0, 0xffff}
	Red     = color.NRGBA64{0xffff, 0, 0, 0xffff}
	Yellow  = color.NRGBA64{0xffff, 0xffff, 0, 0xffff}
	Green   = color.NRGBA64{0, 0xffff, 0, 0xffff}
	Cyan    = color.NRGBA64{0, 0xffff, 0xffff, 0xffff}
	Blue    = color.NRGBA64{0, 0, 0xffff, 0xffff}
	Magenta = color.NRGBA64{0xffff, 0, 0xffff, 0xffff}
)

func ftoui16(n float64) uint16 {
	return uint16(0x7fff + 0x7fff*n)
}

func word(n float64) uint16 {
	return uint16(0xffff * n)
}

// value ranges:
// h: 0-1
// c: 0-1
// l: 0-1
func HclToNRGBA64(h, c, l float64) color.NRGBA64 {
	r, g, b, a := colorful.Hcl(h*360.0, c, l).RGBA()
	return color.NRGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func ColorFuncFromString(name string) (ColorFunc, error) {
	switch name {
	case "smooth":
		return ColorSmooth, nil
	case "bands":
		return ColorBands, nil
	case "mono":
		return ColorMono, nil
	case "stripe":
		return ColorMonoStripe, nil
	case "parti":
		return ColorParti, nil
	case "superparti":
		return ColorSuperParti, nil
	case "check":
		return ColorCheck, nil
	case "e1":
		return ColorExperiment1, nil
	default:
		return nil, errors.New("Invalid ColorFunc name.")
	}
}

func MemberColorFromString(hex string) (color.NRGBA64, error) {
	mc, err := strconv.ParseInt(hex[1:len(hex)], 16, 32)
	if err != nil {
		return color.NRGBA64{0, 0, 0, 0}, err
	}

	member_color := color.NRGBA64{
		uint16(((mc >> 16) & 0xff) * 0x101),
		uint16(((mc >> 8) & 0xff) * 0x101),
		uint16((mc & 0xff) * 0x101),
		0xffff,
	}

	return member_color, nil
}

type ColorFunc func(*Context, complex128, int, int, int, int)

func ColorSmooth(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	log_zn := math.Log(real(z)*real(z)+imag(z)*imag(z)) / 2.0
	nu := math.Log(log_zn/math.Log(float64(c.Power))) / math.Log(float64(c.Power))
	j := float64(i) + 1.0 - nu

	// TODO: this kinda looks like the bands coloring algorithm, but
	// doesn't match. the 4.75 factor is a guess.
	t := (math.Pi / (4.75 * float64(c.Power))) * j

	k := color.NRGBA64{
		ftoui16(math.Sin(math.Pi + t)),
		ftoui16(math.Sin(math.Pi + 0.25*math.Pi + t)),
		ftoui16(math.Cos(math.Pi + t)),
		0xffff,
	}

	c.Image.SetNRGBA64(x, y, k)
}

func ColorBands(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	t := (float64(max_i) / math.Pi) * (float64(i) / float64(max_i))

	k := color.NRGBA64{
		ftoui16(math.Sin(math.Pi + t)),
		ftoui16(math.Sin(math.Pi + 0.25*math.Pi + t)),
		ftoui16(math.Cos(math.Pi + t)),
		0xffff,
	}

	c.Image.SetNRGBA64(x, y, k)
}

func ColorMono(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	if i&1 == 0 {
		c.Image.SetNRGBA64(x, y, White)
	} else {
		c.Image.SetNRGBA64(x, y, Black)
	}
}

func ColorMonoStripe(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	if (i-1)%9 == 0 {
		c.Image.SetNRGBA64(x, y, Accent)
		return
	}

	if i&1 == 0 {
		c.Image.SetNRGBA64(x, y, White)
	} else {
		c.Image.SetNRGBA64(x, y, Black)
	}
}

func ColorCheck(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	p := cmplx.Phase(z)

	if p >= 0 {
		c.Image.SetNRGBA64(x, y, White)
	} else {
		c.Image.SetNRGBA64(x, y, Black)
	}
}

func ColorParti(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	p := cmplx.Phase(z)
	if p > math.Pi/2.0 {
		c.Image.SetNRGBA64(x, y, White)
	} else if p >= 0 {
		c.Image.SetNRGBA64(x, y, Blue)
	} else if p > -1.0*math.Pi/2.0 {
		c.Image.SetNRGBA64(x, y, Red)
	} else if p > -1.0*math.Pi {
		c.Image.SetNRGBA64(x, y, Black)
	}
}

func ColorSuperParti(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	p := cmplx.Phase(z)
	if p > 3.0*math.Pi/4.0 {
		c.Image.SetNRGBA64(x, y, White)
	} else if p > math.Pi/2.0 {
		c.Image.SetNRGBA64(x, y, Red)
	} else if p > math.Pi/4.0 {
		c.Image.SetNRGBA64(x, y, Yellow)
	} else if p >= 0 {
		c.Image.SetNRGBA64(x, y, Green)
	} else if p > math.Pi/-4.0 {
		c.Image.SetNRGBA64(x, y, Cyan)
	} else if p > math.Pi/-2.0 {
		c.Image.SetNRGBA64(x, y, Blue)
	} else if p > 3.0*math.Pi/-4.0 {
		c.Image.SetNRGBA64(x, y, Magenta)
	} else if p > -1*math.Pi {
		c.Image.SetNRGBA64(x, y, Black)
	}
}

func ColorExperiment1(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	log_zn := math.Log(real(z)*real(z)+imag(z)*imag(z)) / 2.0
	nu := math.Log(log_zn/math.Log(float64(c.Power))) / math.Log(float64(c.Power))
	j := float64(i) + 1.0 - nu

	c.Image.SetNRGBA64(x, y, HclToNRGBA64(math.Sin(j), 1.0, 0.0))
}
