package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/Deleplace/samples/mosaicat"
	"github.com/nfnt/resize"
)

var (
	catwidth = flag.Int("catwidth", 32, "size of the side of a cat square")
	ncats    = flag.Int("ncats", 20, "number of cats per line")
)

func main() {
	flag.Parse()
	w, h := *catwidth, *catwidth
	smallcat := resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)

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

	err = mosaicat.Process(in, out, *ncats, *catwidth, smallcat)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to process", inputFilename, ":", err)
		os.Exit(3)
	}

	log.Println("Written", outputFilename)
}

var cat image.Image

func init() {
	// PNG cat
	// f, err := os.Open("cat.png")
	f, err := os.Open("cat_transp.png")
	if err != nil {
		panic(err)
	}
	cat, _, err = image.Decode(f)
	if err != nil {
		panic(err)
	}

	// GIF cat
	// f, err := os.Open("cat_transp.gif")
	// if err != nil {
	// 	panic(err)
	// }
	// g, err := gif.DecodeAll(f)
	// if err != nil {
	// 	panic(err)
	// }
	// cat = g.Image[0]
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "\t", os.Args[0], "input", "output")
	flag.PrintDefaults()
	os.Exit(1)
}
