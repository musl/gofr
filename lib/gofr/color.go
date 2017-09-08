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

func centeredUint16(n float64) uint16 {
	return uint16(n*0x7fff + 0x7fff + 0.5)
}

func fullUint16(n float64) uint16 {
	return uint16(n*0xffff + 0.5)
}

func HclaToNRGBA64(h, c, l, a float64) color.NRGBA64 {
	k := colorful.Hcl(h*360.0, c, l).Clamped()

	return color.NRGBA64{
		fullUint16(k.R),
		fullUint16(k.G),
		fullUint16(k.B),
		fullUint16(a),
	}
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
	case "softspectrum":
		return ColorSoftSpectrum, nil
	case "fire":
		return ColorFire, nil
	case "ice":
		return ColorIce, nil
	case "unicornrainbow":
		return ColorUnicornRainbow, nil
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
		centeredUint16(math.Sin(math.Pi + t)),
		centeredUint16(math.Sin(math.Pi + 0.25*math.Pi + t)),
		centeredUint16(math.Cos(math.Pi + t)),
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
		centeredUint16(math.Sin(math.Pi + t)),
		centeredUint16(math.Sin(math.Pi + 0.25*math.Pi + t)),
		centeredUint16(math.Cos(math.Pi + t)),
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

func ColorSoftSpectrum(ctx *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		ctx.Image.SetNRGBA64(x, y, ctx.MemberColor)
		return
	}

	log_zn := math.Log(real(z)*real(z)+imag(z)*imag(z)) / 2.0
	nu := math.Log(log_zn/math.Log(float64(ctx.Power))) / math.Log(float64(ctx.Power))
	j := float64(i) + 1.0 - nu

	h := 0.5 + 0.5*math.Sin(0.125*math.Pi*j)
	c := 0.5 + 0.333*math.Sin(0.0625*math.Pi*j)
	l := 0.5 + 0.333*math.Sin(0.03125*math.Pi*j)

	ctx.Image.SetNRGBA64(x, y, HclaToNRGBA64(h, c, l, 1.0))
}

func colorSmoothMono(ctx *Context, z complex128, x, y, i, max_i int, hue float64) {
	if i == max_i {
		ctx.Image.SetNRGBA64(x, y, ctx.MemberColor)
		return
	}

	log_zn := math.Log(real(z)*real(z)+imag(z)*imag(z)) / 2.0
	nu := math.Log(log_zn/math.Log(float64(ctx.Power))) / math.Log(float64(ctx.Power))
	j := float64(i) + 1.0 - nu

	c := 0.5 + 0.5*math.Sin(0.0625*math.Pi*j)
	l := 0.5 + 0.5*math.Sin(0.03125*math.Pi*j)

	ctx.Image.SetNRGBA64(x, y, HclaToNRGBA64(hue, c, l, 1.0))
}

func ColorFire(ctx *Context, z complex128, x, y, i, max_i int) {
	colorSmoothMono(ctx, z, x, y, i, max_i, 0.15)
}

func ColorIce(ctx *Context, z complex128, x, y, i, max_i int) {
	colorSmoothMono(ctx, z, x, y, i, max_i, 0.6)
}

func ColorUnicornRainbow(ctx *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		ctx.Image.SetNRGBA64(x, y, ctx.MemberColor)
		return
	}

	log_zn := math.Log(real(z)*real(z)+imag(z)*imag(z)) / 2.0
	nu := math.Log(log_zn/math.Log(float64(ctx.Power))) / math.Log(float64(ctx.Power))
	j := float64(i) + 1.0 - nu

	h := 0.5 + 0.5*math.Sin(0.125*math.Pi*j)
	c := 1.0
	l := 0.8 + 0.2*math.Pow(math.Sin(0.5*math.Pi*j), 8)

	ctx.Image.SetNRGBA64(x, y, HclaToNRGBA64(h, c, l, 1.0))
}

func ColorExperiment1(ctx *Context, z complex128, x, y, i, max_i int) {
	if i == max_i {
		ctx.Image.SetNRGBA64(x, y, ctx.MemberColor)
		return
	}

	if (i-1)%3 == 0 {
		ctx.Image.SetNRGBA64(x, y, Black)
		return
	}

	log_zn := math.Log(real(z)*real(z)+imag(z)*imag(z)) / 2.0
	nu := math.Log(log_zn/math.Log(float64(ctx.Power))) / math.Log(float64(ctx.Power))
	j := float64(i) + 1.0 - nu

	h := 0.5 + 0.5*math.Sin(0.125*math.Pi*j)
	c := 1.0
	l := 0.6

	ctx.Image.SetNRGBA64(x, y, HclaToNRGBA64(h, c, l, 1.0))
}
