package gofr

import (
	"errors"
	"image/color"
	"math"
	"strconv"
)

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
	nu := math.Log(log_zn/math.Log(2.0)) / math.Log(2.0)
	j := float64(i) + 1.0 - nu

	o := math.Pi
	f := math.Pi / 32.0
	t := f * math.Pi * float64(j)
	r := uint16(0x7fff + 0x7fff*math.Sin(o+t))
	g := uint16(0x7fff + 0x7fff*math.Sin(o+0.25*math.Pi+t))
	b := uint16(0x7fff + 0x7fff*math.Cos(o+t))

	l := color.NRGBA64{r, g, b, 0xffff}
	c.Image.SetNRGBA64(x, y, l)
}

func ColorBands(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	o := math.Pi
	f := float64(max_i) / 16.0
	t := f * math.Pi * (float64(i) / float64(max_i))
	r := uint16(0x7fff + 0x7fff*math.Sin(o+t))
	g := uint16(0x7fff + 0x7fff*math.Sin(o+0.25*math.Pi+t))
	b := uint16(0x7fff + 0x7fff*math.Cos(o+t))

	l := color.NRGBA64{r, g, b, 0xffff}
	c.Image.SetNRGBA64(x, y, l)
}

func ColorMono(c *Context, z complex128, x, y, i, max_i int) {
	white := color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	black := color.NRGBA64{0, 0, 0, 0xffff}

	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	if i&1 == 0 {
		c.Image.SetNRGBA64(x, y, white)
	} else {
		c.Image.SetNRGBA64(x, y, black)
	}
}

func ColorMonoStripe(c *Context, z complex128, x, y, i, max_i int) {
	white := color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	black := color.NRGBA64{0, 0, 0, 0xffff}
	accent := color.NRGBA64{0, 0xa000, 0xc000, 0xffff}

	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	if (i-1)%9 == 0 {
		c.Image.SetNRGBA64(x, y, accent)
		return
	}

	if i&1 == 0 {
		c.Image.SetNRGBA64(x, y, white)
	} else {
		c.Image.SetNRGBA64(x, y, black)
	}
}
