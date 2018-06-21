package main

import (
	"image"
	"image/draw"
	"sync"
)

func flipAllF(pics []image.Image) []image.Image {
	result := make([]image.Image, len(pics))
	var wg sync.WaitGroup
	wg.Add(len(pics))
	for i, pic := range pics {
		go func(i int, pic image.Image) {
			result[i] = flipF(pic)
			wg.Done()
		}(i, pic)
	}
	wg.Wait()
	return result
}

func flipF(img image.Image) image.Image {
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
