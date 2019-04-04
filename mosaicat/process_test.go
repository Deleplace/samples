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
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/nfnt/resize"
)

func BenchmarkProcess(b *testing.B) {
	catwidth, ncats := 32, 20
	w, h := catwidth, catwidth
	smallcat := resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inputFilename := "testdata/monalisa.jpg"
		outputFilename := "testdata/out.png"

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

		err = Process(in, out, ncats, catwidth, smallcat)
		if err != nil {
			b.Fatal("Failed to process", inputFilename, ":", err)
			return
		}
	}
}

func BenchmarkProcessInMemory(b *testing.B) {
	catwidth, ncats := 32, 20
	w, h := catwidth, catwidth
	smallcat := resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)
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

		err = Process(in, out, ncats, catwidth, smallcat)
		if err != nil {
			b.Fatal("Failed to process", inputFilename, ":", err)
			return
		}
	}
}

func BenchmarkProcessLargeInMemory(b *testing.B) {
	catwidth, ncats := 128, 50
	w, h := catwidth, catwidth
	smallcat := resize.Resize(uint(w), uint(h), cat, resize.Lanczos3)
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

		err = Process(in, out, ncats, catwidth, smallcat)
		if err != nil {
			b.Fatal("Failed to process", inputFilename, ":", err)
			return
		}
	}
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
