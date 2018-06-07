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
	facade.Definition, _ = fetchDefinition(c)
	tmpl.Execute(w, &facade)
}

type facade struct {
	Definition string
}

var tmpl = template.Must(template.New("foobar").Parse(`
<html>
  <head>
    <title>The Condor observation fanclub</title>
  </head>
  <body>
	<h1>Condor observation 2017</h1>
	
	<div class="definition">
		{{.Definition}}
	</div>
	
  <body>
</html>
`))
