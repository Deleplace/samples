package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 1) Read the files from file system, decode them into in-memory images.
	pics := readPictures()

	// 2) Flip the images vertically.
	// Choose your weapon!
	flipped := flipAllA(pics)
	// flipped := flipAllB(pics)
	// flipped := flipAllC(pics)
	// flipped := flipAllD(pics)
	// flipped := flipAllE(pics)
	// flipped := flipAllF(pics)
	// flipped := flipAllG(pics)

	// 3) Encode and save the flipped images to the file system.
	for i, f := range flipped {
		dstpath := fmt.Sprintf("result/result_%d.jpg", i)
		_ = save(f, dstpath)
	}
}

func readPictures() (pics []image.Image) {
	err := filepath.Walk("./original", func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".jpg") {
			// Not a JPG
			return nil
		}

		img, err := load(path)
		if err != nil {
			return err
		}
		pics = append(pics, img)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return pics
}

//
// Flip vertically all JPG files in current folder and subfolders.
//

func load(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	img, _, err := image.Decode(r)
	return img, err
}

func save(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, nil)
}
