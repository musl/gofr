package gofr

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"runtime"
)

func PrintByteSize(s uint64) string {
	key := "BKMGTPEZY"
	l := int(math.Floor(math.Log(float64(s)) / math.Log(1024)))
	if l < 0 {
		l = 0
	}
	if l > 8 {
		l = 8
	}
	u := string(key[l])
	t := float64(s) / math.Pow(1024.0, float64(l))
	return fmt.Sprintf("%10.4f %s", t, u)
}

func logMemStats(l *log.Logger) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	l.Printf("memory: %s heap in use %s heap idle %s alloc %s total alloc\n",
		PrintByteSize(stats.HeapInuse),
		PrintByteSize(stats.HeapIdle),
		PrintByteSize(stats.Alloc),
		PrintByteSize(stats.TotalAlloc))
}

func TestBlocks(c *Context) int {
	colors := []color.NRGBA64{
		color.NRGBA64{0, 0, 0, 0xffff},
		color.NRGBA64{0x7fff, 0x7fff, 0x7fff, 0xffff},
		color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff},
		color.NRGBA64{0xffff, 0, 0, 0xffff},
		color.NRGBA64{0xffff, 0xffff, 0, 0xffff},
		color.NRGBA64{0, 0xffff, 0, 0xffff},
		color.NRGBA64{0, 0xffff, 0xffff, 0xffff},
		color.NRGBA64{0, 0, 0xffff, 0xffff},
		color.NRGBA64{0xffff, 0, 0xffff, 0xffff},
	}
	Fill(c, colors[c.Id%len(colors)])
	return 0
}

func Fill(c *Context, k color.NRGBA64) {
	min := c.Image.Bounds().Min
	max := c.Image.Bounds().Max

	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			c.Image.Set(x, y, k)
		}
	}
}
