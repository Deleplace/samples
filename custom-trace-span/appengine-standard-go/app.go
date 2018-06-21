package customtracespan

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)

		// A HTTPS request to Wikipedia, to fetch some text
		def, err1 := fetchWikipedia(c)
		if err1 != nil {
			log.Errorf(c, "%v", err1)
		}
		fmt.Fprintf(w, "%v <br/><br/>\n\n", def)

		// A CPU-intensive task
		n := uint64(rand.Int63()) // Make sure n fits in a int64 as well
		log.Infof(c, "Computing Collatz number of steps for %d", n)
		steps := compute(c, n)
		log.Infof(c, "Done: %d steps", steps)
		fmt.Fprintf(w, "Collatz number of steps for %d is %d", n, steps)

		// Save the computation in Datastore
		err2 := save(c, n, steps)
		if err2 != nil {
			log.Errorf(c, "%v", err2)
		}

	})
}

func fetchWikipedia(c context.Context) (string, error) {
	const wikipediaDefinitionURL = "https://en.wikipedia.org/api/rest_v1/page/summary/Collatz_conjecture"

	client := urlfetch.Client(c)
	resp, err := client.Get(wikipediaDefinitionURL)
	if err != nil {
		return "", err
	}
	jsonDecoder := json.NewDecoder(resp.Body)
	data := &struct {
		Extract string `json:"extract"`
	}{}
	err = jsonDecoder.Decode(data)
	if err != nil {
		return "", err
	}
	return data.Extract, nil
}

func compute(c context.Context, n uint64) (steps uint64) {
	for n != 1 {
		if n&1 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		steps++
		time.Sleep(5 * time.Millisecond)
	}
	return steps
}

func save(c context.Context, n, steps uint64) error {
	entity := struct {
		N     int64
		Steps int64
		Date  time.Time
	}{
		int64(n),
		int64(steps),
		time.Now(),
	}

	k := datastore.NewKey(c, "CollatzSteps", "", int64(n), nil)
	_, err := datastore.Put(c, k, &entity)
	return err
}
