package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var state struct {
	lock        sync.Mutex
	unsentNotif int
}

func main() {

	// These 11 lines register to listen to SIGTERM and flush the instance
	// state by sending a batch request to the external service.
	sigs := make(chan os.Signal, 1)
	go func() {
		sig := <-sigs
		log.Println("catched signal", sig)
		ok := flushState(context.Background(), "graceful termination")
		if !ok {
			os.Exit(1)
		}
		os.Exit(0)
	}()
	signal.Notify(sigs, syscall.SIGTERM)

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
		// Let's do it every 20 requests
		state.lock.Lock()
		state.unsentNotif++
		flushNow := (state.unsentNotif >= 20)
		state.lock.Unlock()
		if flushNow {
			flushState(r.Context(), "full buffer")
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

func notifyExternalService(ctx context.Context, batch int) error {
	log.Println("Notifying external service with batch of", batch)
	u := "https://externalservice202010.oa.r.appspot.com/submit"
	err := retry(5, time.Second, func() error {
		resp, err := http.PostForm(u, url.Values{
			"customer":   []string{"customer3"},
			"nbmetadata": []string{strconv.Itoa(batch)},
		})
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			msg, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("external service response: %d %q", resp.StatusCode, string(msg))
		}
		return nil
	})
	return err
}

func flushState(ctx context.Context, motive string) bool {
	log.Printf("Flushing state (%s)", motive)
	var batch int
	state.lock.Lock()
	batch = state.unsentNotif
	state.unsentNotif = 0
	state.lock.Unlock()
	if batch == 0 {
		log.Println("State already clean, no need to call the external service")
	} else {
		err := notifyExternalService(ctx, batch)
		if err != nil {
			log.Println("Notifying external service:", err)
			// Remember N unsent volume
			state.lock.Lock()
			state.unsentNotif += batch
			state.lock.Unlock()
			return false
		}
	}
	return true
}

const homepage = `
Please give us your feedback on Company 3!
<form action="/feedback" method="POST">
	<textarea name="feedback"></textarea>
	<input type="submit">
</form>
`

func retry(attempts int, delay time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(delay)

		log.Println("retrying after attempt", (1 + i), " error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
