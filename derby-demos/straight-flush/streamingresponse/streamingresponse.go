package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		flush := func() {
			if ww, ok := w.(http.Flusher); ok {
				ww.Flush()
			}
		}

		result1 := processing1(ctx)
		fmt.Fprintln(w, result1)
		flush()

		result2 := processing2(ctx)
		fmt.Fprintln(w, result2)
		flush()

		result3 := processing3(ctx)
		fmt.Fprintln(w, result3)
	})

	http.HandleFunc("/nostream", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		result1 := processing1(ctx)
		fmt.Fprintln(w, result1)

		result2 := processing2(ctx)
		fmt.Fprintln(w, result2)

		result3 := processing3(ctx)
		fmt.Fprintln(w, result3)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

func processing1(ctx context.Context) string {
	time.Sleep(1 * time.Second)
	return "Processing 1 complete"
}

func processing2(ctx context.Context) string {
	time.Sleep(2 * time.Second)
	return "Processing 2 complete"
}

func processing3(ctx context.Context) string {
	time.Sleep(3 * time.Second)
	return "Processing 3 complete"
}
