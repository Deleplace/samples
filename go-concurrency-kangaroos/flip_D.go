package main

import (
	"image"
	"sync"
)

func flipAllD(pics []image.Image) []image.Image {
	result := make([]image.Image, len(pics))
	var wg sync.WaitGroup
	wg.Add(len(pics))
	for i, pic := range pics {
		go func(i int, pic image.Image) {
			result[i] = flipD(pic)
			wg.Done()
		}(i, pic)
	}
	wg.Wait()
	return result
}

func flipD(img image.Image) image.Image {
	flipped := image.NewNRGBA(img.Bounds())
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	var wg sync.WaitGroup
	wg.Add(w)
	for x := 0; x < w; x++ {
		go func(x int) {
			for y := 0; y < h; y++ {
				c := img.At(x, y)
				flipped.Set(x, h-y-1, c)
			}
			wg.Done()
		}(x)
	}
	wg.Wait()
	return flipped
}
