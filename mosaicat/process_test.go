package main

import (
	"bytes"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/nfnt/resize"
)

func BenchmarkProcess(b *testing.B) {
	*catwidth, *ncats = 32, 20
	w, h = *catwidth, *catwidth
	smallcat = resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)

	for i := 0; i < b.N; i++ {
		inputFilename := "testdata/monalisa.jpg"
		outputFilename := "out.png"

		in, err := os.Open(inputFilename)
		if err != nil {
			b.Fatal("Couldn't open", inputFilename, "for reading:", err)
			return
		}

		out, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			b.Fatal("Couldn't open", outputFilename, "for writing:", err)
			return
		}

		err = process(in, out)
		if err != nil {
			b.Fatal("Failed to process", inputFilename, ":", err)
			return
		}
	}
}

func BenchmarkProcessInMemory(b *testing.B) {
	*catwidth, *ncats = 32, 20
	w, h = *catwidth, *catwidth
	smallcat = resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)
	inputFilename := "testdata/monalisa.jpg"

	indata, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		b.Fatal("Couldn't open", inputFilename, "for reading:", err)
		return
	}
	outmem := make([]byte, 20*1024*1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := bytes.NewBuffer(indata)
		out := bytes.NewBuffer(outmem)

		err = process(in, out)
		if err != nil {
			b.Fatal("Failed to process", inputFilename, ":", err)
			return
		}
	}
}
