// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mosaicat

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
)

func Process(in io.Reader, out io.Writer, ncats int, catwidth int, smallcat image.Image) error {
	w, h := catwidth, catwidth
	src, format, err := image.Decode(in)
	if err != nil {
		return err
	}
	// log.Println("Decoded a", format)
	_ = format

	if srcW := src.Bounds().Max.X; srcW < ncats {
		return fmt.Errorf("Can't have input width (%d) less than ncats (%d)", srcW, ncats)
	}

	W, H := src.Bounds().Max.X, src.Bounds().Max.Y

	WW := ncats * catwidth
	HH := (WW * H) / W
	dst := image.NewRGBA(image.Rect(0, 0, WW, HH))

	// draw.Draw(dst, dst.Bounds(), src, image.ZP, draw.Src)

	for x := 0; x < WW; x += w {
		for y := 0; y < HH; y += h {
			c := src.At((x*W+w/2)/WW, (y*H+h/2)/HH)
			colorcat := colorizeCat(c, w, h, smallcat)
			// log.Println("Drawing at", x, y)
			dstR := image.Rect(x, y, x+w, y+h)
			draw.Draw(dst, dstR, colorcat, image.ZP, draw.Over)
		}
	}

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

func colorizeCat(c color.Color, w, h int, smallcat image.Image) image.Image {
	// Replace 1 pixel at a time
	cc := image.NewRGBA(image.Rect(0, 0, w, h))
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
			current := cc.RGBAAt(x, y)
			if similarRGBA(current, blue) {
				bluecount++
				cc.Set(x, y, lo)
			}
			if isTransparentRGBA(current) {
				backcount++
				cc.Set(x, y, hi)
			}
		}
	}
	// log.Println("Blue area =", bluecount, "/", total)
	// log.Println("Back area =", backcount, "/", total)
	return cc
}

func colorizeCatPaletted(c color.Color) *image.Paletted {
	var p *image.Paletted
	blue := color.RGBA{
		R: 109,
		G: 159,
		B: 208,
		A: 255,
	}

	hi, lo := colorPair(c)

	altPalette := make(color.Palette, len(p.Palette))
	copy(altPalette, p.Palette)
	for i, x := range altPalette {
		if similar(x, blue) {
			altPalette[i] = lo
		}
		if isTransparent(x) {
			altPalette[i] = hi
		}
	}

	pp := image.Paletted{
		Pix:     p.Pix,
		Stride:  p.Stride,
		Rect:    p.Rect,
		Palette: altPalette,
	}
	return &pp
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

func similarRGBA(c1, c2 color.RGBA) bool {
	dr := int(c1.R) - int(c2.R)
	dg := int(c1.G) - int(c2.G)
	db := int(c1.B) - int(c2.B)
	da := int(c1.A) - int(c2.A)
	return dr*dr+dg*dg+db*db+da*da < 30500
}

func isTransparentRGBA(c color.RGBA) bool {
	return c.A <= 10
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
