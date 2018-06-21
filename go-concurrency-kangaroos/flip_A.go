package main

import (
	"image"
)

func flipAllA(pics []image.Image) []image.Image {
	result := make([]image.Image, len(pics))
	for i, pic := range pics {
		result[i] = flipA(pic)
	}
	return result
}

func flipA(img image.Image) image.Image {
	flipped := image.NewNRGBA(img.Bounds())
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := img.At(x, y)
			flipped.Set(x, h-y-1, c)
		}
	}
	return flipped
}
