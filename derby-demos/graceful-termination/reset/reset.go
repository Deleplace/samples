package main

import (
	"context"
	"log"

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
	ctx := context.Background()
	for _, customerID := range []string{
		"customer1",
		"customer2",
		"customer3",
	} {
		docRef := fsClient.Collection("customers").Doc(customerID)
		var data CustomerData
		_, err := docRef.Set(ctx, data)
		if err != nil {
			log.Println("updating customer data in transaction:", err)
			return
		}
	}

	log.Println("Done.")
}

func init() {
	ctx := context.Background()
	var err error
	fsClient, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

}
