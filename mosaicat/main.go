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
	catwidth = flag.Int("catwidth", 32, "size of the side of a cat square")
	ncats    = flag.Int("ncats", 20, "number of cats per line")
)

func main() {
	flag.Parse()
	w, h = *catwidth, *catwidth
	smallcat = resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)

	if flag.NArg() != 2 {
		log.Println(flag.NArg())
		usage()
	}

	inputFilename := flag.Arg(0)
	outputFilename := flag.Arg(1)

	in, err := os.Open(inputFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't open", inputFilename, "for reading:", err)
		os.Exit(1)
	}

	out, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't open", outputFilename, "for writing:", err)
		os.Exit(2)
	}

	err = process(in, out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to process", inputFilename, ":", err)
		os.Exit(3)
	}

	log.Println("Written", outputFilename)
}

func process(in io.Reader, out io.Writer) error {
	src, format, err := image.Decode(in)
	if err != nil {
		return err
	}
	log.Println("Decoded a", format)

	if srcW := src.Bounds().Max.X; srcW < *ncats {
		return fmt.Errorf("Can't have input width (%d) less than ncats (%d)", srcW, *ncats)
	}

	W, H := src.Bounds().Max.X, src.Bounds().Max.Y

	WW := *ncats * *catwidth
	HH := (WW * H) / W
	dst := image.NewRGBA(image.Rect(0, 0, WW, HH))

	// draw.Draw(dst, dst.Bounds(), src, image.ZP, draw.Src)

	for x := 0; x < WW; x += w {
		for y := 0; y < HH; y += h {
			c := src.At((x*W+w/2)/WW, (y*H+h/2)/HH)
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
			current := cc.At(x, y)
			if similar(current, blue) {
				bluecount++
				cc.Set(x, y, lo)
			}
			if isTransparent(current) {
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
var w, h int

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
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "\t", os.Args[0], "input", "output")
	flag.PrintDefaults()
	os.Exit(1)
}
