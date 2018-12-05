package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	aelog "google.golang.org/appengine/log"
)

var sd *stackdriver.Exporter

func main() {
	rand.Seed(time.Now().UnixNano())

	initSD()
	defer sd.Flush()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		c := appengine.NewContext(r)

		// A HTTPS request to Wikipedia, to fetch some text
		def, err1 := fetchWikipedia(c)
		if err1 != nil {
			aelog.Errorf(c, "%v", err1)
		}
		fmt.Fprintf(w, "%v <br/><br/>\n\n", def)

		// A CPU-intensive task
		n := uint64(rand.Int63()) // Make sure n fits in a int64 as well
		aelog.Infof(c, "Computing Collatz number of steps for %d", n)
		c2, endSpan := startSpanfWRT(r, "Compute Collatz number of steps for %d", n)
		steps := compute(c2, n)
		endSpan()
		aelog.Infof(c, "Done: %d steps", steps)
		fmt.Fprintf(w, "Collatz number of steps for %d is %d", n, steps)

		// Save the computation in Datastore
		err2 := save(c, n, steps)
		if err2 != nil {
			aelog.Errorf(c, "%v", err2)
		}
	})

	appengine.Main()
}

func fetchWikipedia(c context.Context) (string, error) {
	const wikipediaDefinitionURL = "https://en.wikipedia.org/api/rest_v1/page/summary/Collatz_conjecture"

	resp, err := http.Get(wikipediaDefinitionURL)
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

func initSD() {
	var err error
	sd, err = stackdriver.NewExporter(stackdriver.Options{
		ProjectID:    os.Getenv("GOOGLE_CLOUD_PROJECT"),
		MetricPrefix: "demo-prefix",
	})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	}

	// Configure 100% sample rate
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Register it as a trace exporter
	trace.RegisterExporter(sd)
}

// Start a new span "With Remote Parent"
func startSpanfWRT(r *http.Request, msg string, args ...interface{}) (c2 context.Context, endSpan func()) {
	caption := fmt.Sprintf(msg, args...)
	c := r.Context()

	spanContext, ok := (&propagation.HTTPFormat{}).SpanContextFromRequest(r)
	if !ok {
		return c, func() {}
	}
	var span *trace.Span
	c2, span = trace.StartSpanWithRemoteParent(c, caption, spanContext)
	endSpan = func() {
		span.End()
	}
	return c2, endSpan
}
