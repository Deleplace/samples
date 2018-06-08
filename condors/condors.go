package condors

import (
	"context"
	"html/template"
	"net/http"
	"sort"
	"sync"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", frontPage)
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var facade facade
	var err1, err2, err3 error
	facade.Year = 2017

	RunConcurrent(
		func() {
			log.Infof(c, "Fetching Wikipedia definition: start...")
			facade.Definition, err1 = fetchDefinition(c)
			facade.takeError(c, err1)
			log.Infof(c, "Fetching Wikipedia definition: done.")
		},

		func() {
			log.Infof(c, "Querying observations winner: start...")
			facade.HasWinner, facade.WinningObservation, err2 = computeWinner(c, facade.Year)
			facade.takeError(c, err2)
			log.Infof(c, "Querying observations winner: done.")
		},

		func() {
			log.Infof(c, "Incrementing pageviews counter: start...")
			facade.PageViews, err3 = getAndIncrementHitCount(c, "/")
			facade.takeError(c, err3)
			log.Infof(c, "Incrementing pageviews counter: done.")
		},
	)

	tmpl.Execute(w, &facade)
}

type facade struct {
	Definition string
	Year       int

	HasWinner          bool
	WinningObservation Observation

	PageViews   int
	errorsMutex sync.Mutex
	Errors      []error
}

func (facade *facade) takeError(c context.Context, err error) {
	if err != nil {
		// Error will be displayed on rendered page
		facade.errorsMutex.Lock()
		facade.Errors = append(facade.Errors, err)
		facade.errorsMutex.Unlock()

		if c != nil {
			// Error will also be logged server-side
			log.Errorf(c, "%v", err)
		}
	}
}

// The winner is the user who observed the greatest number of condors.
// In case of a tie, the one who observed first wins.
func computeWinner(c context.Context, year int) (found bool, obs Observation, err error) {
	var candidates []Observation
	candidates, err = queryObservations(c, year)
	if err != nil {
		return
	}
	if len(candidates) == 0 {
		return
	}
	sort.Slice(candidates, func(i, j int) bool {
		// We want NbCondors DESC, date ASC
		switch {
		case candidates[i].NbCondors > candidates[j].NbCondors:
			return true
		case candidates[j].NbCondors > candidates[i].NbCondors:
			return false
		default:
			return candidates[j].Date.Before(candidates[i].Date)
		}
	})
	return true, candidates[0], nil
}

var tmpl = template.Must(template.New("foobar").Parse(`
<html>
  <head>
    <title>The Condor observation fanclub</title>
	<link rel="stylesheet" href="/static/condors.css" />
	<link rel="SHORTCUT ICON" href="/static/favicon.png" />
  </head>
  <body>
	<img src="static/condor.jpg" />
	<h1>Condor observation {{.Year}}</h1>
	
	<div class="definition">
		{{.Definition}}
	</div>

	<div class="winner">
		{{if .HasWinner}}
			{{with .WinningObservation}}
				The winner so far is: <b>{{.Username}}</b>, who observed <b>{{.NbCondors}} condors</b> on {{.Date.Format "2006-01-02"}} in {{.Region}}!
			{{end}}
		{{else}}
			<i>There are no known observations in {{.Year}}, so far.</i>
		{{end}}
	</div>	

	<div class="hit-count">
		This page was viewed {{.PageViews}} times.
	</div>

	{{range .Errors}}
		<div class="error">
		  ⚠ {{.}}
		</div>
	{{end}}
  <body>
</html>
`))
