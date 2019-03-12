package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"

	"github.com/nfnt/resize"
)

var (
	inputFilename = flag.String("aaa", "", "it's the thing")
	count         = flag.Int("count", 0, "it's the number")
)

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	inputFilename := os.Args[1]
	outputFilename := os.Args[2]

	in, err := os.Open(inputFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't open", inputFilename, "for reading:", err)
		return
	}

	out, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't open", outputFilename, "for writing:", err)
		return
	}

	err = process(in, out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to process", inputFilename, ":", err)
		return
	}

	log.Println("Written", outputFilename)
}

func process(in io.Reader, out io.Writer) error {
	src, format, err := image.Decode(in)
	if err != nil {
		return err
	}
	log.Println("Decoded a", format)

	W, H := src.Bounds().Max.X, src.Bounds().Max.Y
	dst := image.NewRGBA(image.Rect(0, 0, W, H))

	draw.Draw(dst, dst.Bounds(), src, image.ZP, draw.Src)

	for x := 0; x < W; x += w {
		for y := 0; y < H; y += h {
			c := src.At(x+w/2, y+h/2)
			colorcat := colorizeCat(c)
			// log.Println("Drawing at", x, y)
			dstR := image.Rect(x, y, x+w, y+h)
			draw.Draw(dst, dstR, colorcat, image.ZP, draw.Over)
		}
	}

	_ = smallcat
	err = png.Encode(out, dst)
	return err
}

func avg(img image.Image, rect image.Rectangle) color.Color {
	var r, g, b, a uint32
	for x := rect.Bounds().Min.X; x < rect.Bounds().Max.X; x++ {
		for y := rect.Bounds().Min.Y; y < rect.Bounds().Max.Y; y++ {
			rr, gg, bb, aa := img.At(x, y).RGBA()
			r += rr
			g += gg
			b += bb
			a += aa
		}
	}
	n := uint32((rect.Bounds().Max.X - rect.Bounds().Min.X) * (rect.Bounds().Max.Y - rect.Bounds().Min.Y))

	return color.RGBA{
		R: uint8(r / n >> 8),
		G: uint8(r / n >> 8),
		B: uint8(r / n >> 8),
		A: uint8(r / n >> 8),
	}
}

func colorizeCat(c color.Color) image.Image {
	// Replace 1 pixel at a time
	cc := image.NewRGBA(image.Rect(0, 0, w, h))
	back := color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	blue := color.RGBA{
		R: 109,
		G: 159,
		B: 208,
		A: 255,
	}
	total, bluecount, backcount := 0, 0, 0
	draw.Draw(cc, cc.Bounds(), smallcat, image.ZP, draw.Src)

	hi, lo := colorPair(c)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			total++
			if similar(cc.At(x, y), blue) {
				bluecount++
				cc.Set(x, y, lo)
			}
			if isTransparent(cc.At(x, y)) {
				backcount++
				cc.Set(x, y, hi)
			}
		}
	}
	// log.Println("Blue area =", bluecount, "/", total)
	// log.Println("Back area =", backcount, "/", total)
	_ = back
	return cc
}

func similar(c1, c2 color.Color) bool {
	r, g, b, a := c1.RGBA()
	rr, gg, bb, aa := c2.RGBA()
	dr := int64(r) - int64(rr)
	dg := int64(g) - int64(gg)
	db := int64(b) - int64(bb)
	da := int64(a) - int64(aa)
	return dr*dr+dg*dg+db*db+da*da < 2000000000
}

func isTransparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a < 1000
}

// Produce 1 bright color + 1 dark color
func colorPair(c color.Color) (hi, lo color.Color) {
	r, g, b, a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	maxd := uint32(40)
	if r < maxd {
		maxd = r
	}
	if g < maxd {
		maxd = g
	}
	if b < maxd {
		maxd = b
	}
	if r > 128 && (255-r) < maxd {
		maxd = (255 - r)
	}
	if g > 128 && (255-g) < maxd {
		maxd = (255 - g)
	}
	if b > 128 && (255-b) < maxd {
		maxd = (255 - b)
	}
	hi = color.RGBA{
		R: uint8(r + maxd),
		G: uint8(g + maxd),
		B: uint8(b + maxd),
		A: uint8(a),
	}
	lo = color.RGBA{
		R: uint8(r - maxd),
		G: uint8(g - maxd),
		B: uint8(b - maxd),
		A: uint8(a),
	}
	return hi, lo
}

var cat, smallcat image.Image
var w, h = 24, 24

func init() {
	// f, err := os.Open("cat.png")
	f, err := os.Open("cat_transp.png")
	if err != nil {
		panic(err)
	}
	cat, _, err = image.Decode(f)
	if err != nil {
		panic(err)
	}
	smallcat = resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "\t", os.Args[0], "input", "output")
	flag.PrintDefaults()
	os.Exit(1)
}
