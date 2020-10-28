package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		fmt.Fprintln(w, homepage)
	})
	http.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		if r.Method != "POST" {
			http.Error(w, `{"message": "POST only"}`, http.StatusBadRequest)
			log.Printf("Unexpected HTTP method %q\n", r.Method)
			return
		}

		feedback := r.FormValue("feedback")
		log.Printf("User feedback: %q\n", feedback)

		// Contractually mandated to notify the external service
		u := "https://externalservice202010.oa.r.appspot.com/submit"
		resp, err := http.PostForm(u, url.Values{
			"customer": []string{"customer1"},
		})
		if err != nil {
			log.Println("Notifying external service:", err)
		}
		if resp.StatusCode != 200 {
			msg, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Notifying external service: %d %q\n", resp.StatusCode, string(msg))
		}

		// Other important work to process the feedback
		time.Sleep(time.Duration(1000+rand.Intn(1000)) * time.Millisecond)

		fmt.Fprintf(w, `{"message": "Thank you for your feedback"}`)
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

const homepage = `
Please give us your feedback on Company 1!
<form action="/feedback" method="POST">
	<textarea name="feedback"></textarea>
	<input type="submit">
</form>
`
