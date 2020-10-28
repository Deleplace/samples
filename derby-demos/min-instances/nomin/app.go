package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func init() {
	log.Println("App init start...")
	time.Sleep(5 * time.Second)
	log.Println("App init done.")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request processing start...")
		time.Sleep(1 * time.Second)
		log.Println("Request processing done.")
		fmt.Fprintln(w, "ok")
	})
	err := http.ListenAndServe(":8080", nil)
	log.Fatalln(err)
}
