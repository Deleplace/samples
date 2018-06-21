package main

import (
	"image"
	"image/draw"
)

func flipAllG(pics []image.Image) []image.Image {
	result := make([]image.Image, 0, len(pics))
	ch := make(chan image.Image)
	for _, pic := range pics {
		go func(pic image.Image) {
			ch <- flipG(pic)
		}(pic)
	}
	for flipped := range ch {
		result = append(result, flipped)
	}
	// What is the gotcha here?
	return result
}

func flipG(img image.Image) image.Image {
	flipped := image.NewNRGBA(img.Bounds())
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	draw.Draw(flipped, img.Bounds(), img, img.Bounds().Min, draw.Src)
	for y := 0; y < h/2; y++ {
		for x := 0; x < w; x++ {
			c1 := img.At(x, y)
			c2 := img.At(x, h-y-1)
			flipped.Set(x, y, c2)
			flipped.Set(x, h-y-1, c1)
		}
	}
	return flipped
}
