package condors

import (
	"context"
	"html/template"
	"net/http"
	"sort"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", frontPage)
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var facade facade
	var err1, err2, err3 error
	facade.Year = 2017

	facade.Definition, err1 = fetchDefinition(c)
	facade.takeError(err1)

	facade.HasWinner, facade.WinningObservation, err2 = computeWinner(c, facade.Year)
	facade.takeError(err2)

	facade.PageViews, err3 = getAndIncrementHitCount(c, "/")
	facade.takeError(err3)

	tmpl.Execute(w, &facade)
}

type facade struct {
	Definition string
	Year       int

	HasWinner          bool
	WinningObservation Observation

	PageViews int
	Errors    []error
}

func (facade *facade) takeError(err error) {
	if err != nil {
		facade.Errors = append(facade.Errors, err)
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
