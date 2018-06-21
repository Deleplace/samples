package main

import (
	"image"
	"sync"
)

func flipAllE(pics []image.Image) []image.Image {
	result := make([]image.Image, len(pics))
	var wg sync.WaitGroup
	wg.Add(len(pics))
	for i, pic := range pics {
		go func(i int, pic image.Image) {
			result[i] = flipE(pic)
			wg.Done()
		}(i, pic)
	}
	wg.Wait()
	return result
}

func flipE(img image.Image) image.Image {
	flipped := image.NewNRGBA(img.Bounds())
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	var wg sync.WaitGroup
	wg.Add(h)
	for y := 0; y < h; y++ {
		go func(y int) {
			for x := 0; x < w; x++ {
				c := img.At(x, y)
				flipped.Set(x, h-y-1, c)
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
	return flipped
}
