package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
)

const projectID = "externalservice202010"

var fsClient *firestore.Client

type CustomerData struct {
	// The total number of end user requests reported
	Volume int `firestore:"volume"`
	// The number of calls to this API
	APICalls int `firestore:"api_calls"`
	// The total price for this customer
	TotalPrice float64 `firestore:"total_price"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving request")
		fmt.Fprintln(w, "Hello")
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		if r.Method != "POST" {
			http.Error(w, `{"message": "POST only"}`, http.StatusBadRequest)
			log.Printf("Unexpected HTTP method %q\n", r.Method)
			return
		}
		ctx := r.Context()
		customerID := r.FormValue("customer")
		if customerID == "" {
			http.Error(w, `{"message": "mandatory parameter customer"}`, http.StatusBadRequest)
			log.Printf("No parameter customer provided\n")
			return
		}
		nbMetadata := 1
		nbMetadataStr := r.FormValue("nbmetadata")
		if nbMetadataStr == "" {
			log.Println("Reporting default volume 1 for", customerID)
		} else {
			log.Println("Reporting volume", nbMetadataStr, "for", customerID)
			var err error
			nbMetadata, err = strconv.Atoi(nbMetadataStr)
			if err != nil {
				http.Error(w, `{"message": "invalid parameter nbmetadata"}`, http.StatusBadRequest)
				log.Printf("invalid parameter nbmetadata %q\n", nbMetadataStr)
				return
			}
		}
		metadata := r.FormValue("metadata")
		// The customer sends us stuff but we don't really care
		_ = metadata
		// We care about charging them, however
		docRef := fsClient.Collection("customers").Doc(customerID)

		err := fsClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			doc, err := tx.Get(docRef)
			if err != nil {
				return fmt.Errorf("loading customer data: %v", err)
			}
			if !doc.Exists() {
				return fmt.Errorf("no such customer %q", customerID)
			}
			var data CustomerData
			err = doc.DataTo(&data)
			if err != nil {
				return err
			}
			return tx.Update(docRef, []firestore.Update{
				{Path: "volume", Value: data.Volume + nbMetadata},
				{Path: "total_price", Value: data.TotalPrice + 0.01},
				{Path: "api_calls", Value: data.APICalls + 1},
			})
		})
		if err != nil {
			http.Error(w, `{"message": "transaction error :("}`, http.StatusInternalServerError)
			log.Println("updating customer data in transaction:", err)
			return
		}

		fmt.Fprintln(w, `{"message": "ok"}`)
	})

	http.HandleFunc("/dashboard", dashboard)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

func init() {
	ctx := context.Background()
	var err error
	fsClient, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

}

func dashboard(w http.ResponseWriter, r *http.Request) {
	var customers [3]struct {
		ID   string
		Data CustomerData
	}

	for i, customerID := range []string{
		"customer1",
		"customer2",
		"customer3",
	} {
		docRef := fsClient.Collection("customers").Doc(customerID)
		doc, err := docRef.Get(r.Context())
		if err != nil {
			log.Printf("loading customer data: %v\n", err)
		}
		if !doc.Exists() {
			log.Printf("no such customer %q\n", customerID)
		}
		var data CustomerData
		err = doc.DataTo(&data)
		if err != nil {
			log.Printf("loading customer data: %v\n", err)
		}
		customers[i].ID = customerID
		customers[i].Data = data
	}
	w.Header().Set("Content-Type", "text/html")
	err := dashboardTmpl.Execute(w, customers)
	if err != nil {
		log.Println("Templating:", err)
	}
}

var dashboardTmpl = template.Must(template.New("dashboard").
	Funcs(template.FuncMap{
		"currency": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
	}).Parse(dashboardHTML))

var dashboardHTML = `
	<meta charset="utf-8">
	<style>
		table {
			margin-left: 1em;
			border-collapse: collapse;
		}
		tr {
			border: 1px solid #BBF;
		}
		th {
			padding: 0.5em;
			border: 1px solid #BBF;
		}
		td {
			padding: 0.5em;
			border: 1px solid #BBF;
		}
	</style>
	<h2>External Service accounts dashboard<h1>

	<table>
		<tr>
			<th>ID</th>
			<th>Volume</th>
			<th>API calls</th>
			<th>Total price</th>
		</tr>
		{{range .}}
		<tr>
			<td>{{.ID}}</td>
			<td>{{.Data.Volume}}</td>
			<td>{{.Data.APICalls}}</td>
			<td>{{currency .Data.TotalPrice}} €</td>
		</tr>
		{{end}}
	</table>
`
