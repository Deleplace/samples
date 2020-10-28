package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	go func() {
		sig := <-sigs
		fmt.Println("catched signal", sig)
	}()
	signal.Notify(sigs, syscall.SIGTERM)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving request")
		fmt.Fprintln(w, "Hello")
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
