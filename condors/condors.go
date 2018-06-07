package condors

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", frontPage)
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var facade facade
	var err1, err2 error
	facade.Definition, err1 = fetchDefinition(c)
	facade.takeError(err1)
	facade.PageViews, err2 = getAndIncrementHitCount(c, "/")
	facade.takeError(err2)
	tmpl.Execute(w, &facade)
}

type facade struct {
	Definition string
	PageViews  int

	Errors []error
}

func (facade *facade) takeError(err error) {
	if err != nil {
		facade.Errors = append(facade.Errors, err)
	}
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
	<h1>Condor observation 2017</h1>
	
	<div class="definition">
		{{.Definition}}
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
