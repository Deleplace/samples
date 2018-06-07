package condors

import (
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", frontPage)
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

var tmpl = template.Must(template.New("foobar").Parse(`
<html>
  <head>
    <title>The Condor observation fanclub</title>
  </head>
  <body>
    <h1>Condor observation 2017</h1>
	
  <body>
</html>
`))
