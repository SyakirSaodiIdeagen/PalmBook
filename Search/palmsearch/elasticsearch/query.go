package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func search() {
	es := GetEsClient()

	// Define the query in JSON format
	query := `{
		"query": {
			"match": {
				"title": "golang"
			}
		}
	}`

	// Create the search request
	req := esapi.SearchRequest{
		Index: []string{"my_index"}, // Change to your index name
		Body:  strings.NewReader(query),
	}

	// Execute the request
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error executing search query: %s", err)
	}
	defer res.Body.Close()

	// Read and parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing response: %s", err)
	}

	// Print the response
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	fmt.Println("Search Results:")
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})["_source"]
		fmt.Println(hitMap)
	}
}
