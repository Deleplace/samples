package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

const N = 100

func main() {
	bench("https://nomin-gfcmoiomza-uc.a.run.app")
}

func bench(url string) {
	for i := 0; i < N; i++ {
		t := time.Now()
		resp, err := http.Get(url)
		d := time.Since(t)
		switch {
		case err != nil:
			log.Printf("getting %s: %v\n", url, err)
		case resp.StatusCode != 200:
			log.Printf("getting %s expected 200, got: %v\n", url, resp.StatusCode)
		default:
			ms := int(d / time.Millisecond)
			log.Println(ms)
		}
		time.Sleep(time.Duration(rand.Intn(1_500_000)) * time.Millisecond)
	}
}
