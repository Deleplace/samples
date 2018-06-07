package condors

import (
	"context"
	"encoding/json"

	"google.golang.org/appengine/urlfetch"
)

const wikipediaDefinitionURL = "https://en.wikipedia.org/api/rest_v1/page/summary/Condor"

func fetchDefinition(c context.Context) (string, error) {
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
