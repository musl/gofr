package gofr

import (
	"errors"
	"image/color"
	"math"
	"math/cmplx"
	"strconv"
)

func ftoui16(n float64) uint16 {
	return uint16(0x7fff + 0x7fff*n)
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
	case "e00":
		return ColorExperiment00, nil
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

	o := math.Pi
	f := math.Pi / 32.0 * float64(c.Power)
	t := f * math.Pi * float64(j)
	r := ftoui16(math.Sin(o + t))
	g := ftoui16(math.Sin(o + 0.25*math.Pi + t))
	b := ftoui16(math.Cos(o + t))

	l := color.NRGBA64{r, g, b, 0xffff}
	c.Image.SetNRGBA64(x, y, l)
}

func ColorBands(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	o := math.Pi
	f := float64(max_i) / 16.0 * float64(c.Power)
	t := f * math.Pi * (float64(i) / float64(max_i))
	r := ftoui16(math.Sin(o + t))
	g := ftoui16(math.Sin(o + 0.25*math.Pi + t))
	b := ftoui16(math.Cos(o + t))

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

func ColorCheck(c *Context, z complex128, x, y, i, max_i int) {
	white := color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	black := color.NRGBA64{0, 0, 0, 0xffff}

	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	p := cmplx.Phase(z)

	if p >= 0 {
		c.Image.SetNRGBA64(x, y, white)
	} else {
		c.Image.SetNRGBA64(x, y, black)
	}
}

func ColorParti(c *Context, z complex128, x, y, i, max_i int) {
	white := color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	black := color.NRGBA64{0, 0, 0, 0xffff}
	red := color.NRGBA64{0xffff, 0, 0, 0xffff}
	blue := color.NRGBA64{0, 0, 0xffff, 0xffff}

	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	p := cmplx.Phase(z)
	if p > math.Pi/2.0 {
		c.Image.SetNRGBA64(x, y, white)
	} else if p >= 0 {
		c.Image.SetNRGBA64(x, y, blue)
	} else if p > -1.0*math.Pi/2.0 {
		c.Image.SetNRGBA64(x, y, red)
	} else if p > -1.0*math.Pi {
		c.Image.SetNRGBA64(x, y, black)
	}
}

func ColorSuperParti(c *Context, z complex128, x, y, i, max_i int) {
	white := color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	black := color.NRGBA64{0, 0, 0, 0xffff}
	red := color.NRGBA64{0xffff, 0, 0, 0xffff}
	yellow := color.NRGBA64{0xffff, 0xffff, 0, 0xffff}
	green := color.NRGBA64{0, 0xffff, 0, 0xffff}
	cyan := color.NRGBA64{0, 0xffff, 0xffff, 0xffff}
	blue := color.NRGBA64{0, 0, 0xffff, 0xffff}
	magenta := color.NRGBA64{0xffff, 0, 0xffff, 0xffff}

	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	p := cmplx.Phase(z)
	if p > 3.0*math.Pi/4.0 {
		c.Image.SetNRGBA64(x, y, white)
	} else if p > math.Pi/2.0 {
		c.Image.SetNRGBA64(x, y, red)
	} else if p > math.Pi/4.0 {
		c.Image.SetNRGBA64(x, y, yellow)
	} else if p >= 0 {
		c.Image.SetNRGBA64(x, y, green)
	} else if p > math.Pi/-4.0 {
		c.Image.SetNRGBA64(x, y, cyan)
	} else if p > math.Pi/-2.0 {
		c.Image.SetNRGBA64(x, y, blue)
	} else if p > 3.0*math.Pi/-4.0 {
		c.Image.SetNRGBA64(x, y, magenta)
	} else if p > math.Pi {
		c.Image.SetNRGBA64(x, y, black)
	}
}

func ColorExperiment00(c *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		c.Image.SetNRGBA64(x, y, c.MemberColor)
		return
	}

	r := ftoui16(math.Sin(cmplx.Abs(z)/c.EscapeRadius + cmplx.Phase(z)/math.Pi))
	g := ftoui16(math.Sin(cmplx.Abs(z)/c.EscapeRadius + cmplx.Phase(z)/math.Pi))
	b := ftoui16(math.Sin(cmplx.Abs(z)/c.EscapeRadius + cmplx.Phase(z)/math.Pi))

	l := color.NRGBA64{r, g, b, 0xffff}
	c.Image.SetNRGBA64(x, y, l)
}
