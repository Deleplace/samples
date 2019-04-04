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

package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Deleplace/samples/mosaicat"
	"github.com/nfnt/resize"
)

func main() {
	http.HandleFunc("/", form)
	http.HandleFunc("/process", process)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func process(w http.ResponseWriter, r *http.Request) {
	catwidth := 32
	ncats := 20
	var err error
	if s := r.FormValue("catwidth"); s != "" {
		catwidth, err = strconv.Atoi(s)
		if err != nil {
			http.Error(w, "catwidth", 400)
			return
		}
	}
	if s := r.FormValue("ncats"); s != "" {
		ncats, err = strconv.Atoi(s)
		if err != nil {
			http.Error(w, "ncats", 400)
			return
		}
	}
	log.Printf("Settings ncats=%d, catwidth=%d\n", ncats, catwidth)

	var file multipart.File
	file, _, err = r.FormFile("pic")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var pic bytes.Buffer
	var n int64
	n, err = pic.ReadFrom(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("Processing picture of size %d\n", n)
	start := time.Now()
	ww, hh := catwidth, catwidth
	smallcat := resize.Resize(uint(ww), uint(hh), cat, resize.Lanczos3)

	var result bytes.Buffer
	err = mosaicat.Process(&pic, &result, ncats, catwidth, smallcat)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("Processing successful in", time.Since(start))
	result.WriteTo(w)
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

func form(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<html>
<h1>Mosaicat</h1>
<form action="/process" method="POST" enctype="multipart/form-data">
	<div>
		<label for="ncats">Cat tiles (horizontal)</label>
		<input name="ncats" value="40"></input>
	</div>
	<div>
		<label for="catwidth">Cat width in pixels</label>
		<input name="catwidth" value="64"></input>
	</div>
	<div>
	<input type="file" name="pic"></input>
	</div>
	<button type="submit">Generate</button>
</form>
<style>
	input:not([type]), input[type='text'] {
		width: 4rem;
	}
	input {
		margin-bottom: 1rem;
	}
	label {
		display: inline-block;
		width: 9rem;
	}
</style>
</html>
	`)
}
