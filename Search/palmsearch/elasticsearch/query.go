package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticsearchHit struct {
	Name                 string `json:"Name"`
	CreatedDateTime      string `json:"createdDateTime"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	WebUrl               string `json:"webUrl"`
	Source               string `json:"source"`
}

func Search(param string) []ElasticsearchHit {
	es := GetEsClient()

	// Escape the param to make it safe for wildcard query
	//param = strings.ReplaceAll(param, ".", "\\.")
	//param = strings.ReplaceAll(param, "?", "\\?")

	query := fmt.Sprintf(`{
		"query": {
			"wildcard": {
				"Name": {
					"value": "*%v*"
				}
			}
		}
	}`, param)

	req := esapi.SearchRequest{
		Index: []string{"golang-bulk-index3"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error executing search query: %s", err)
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing response: %s", err)
	}

	// Check if hits exist
	hits, found := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !found {
		log.Println("No hits found")
		return nil
	}

	eshits := make([]ElasticsearchHit, 0)

	fmt.Println("Search Results:")
	for _, hit := range hits {
		// Assert the hit to map[string]interface{}
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			fmt.Println("Error: hit is not a valid map")
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: _source field is missing or not of expected type")
			continue
		}

		var esHit ElasticsearchHit

		// Map the data into the struct
		esHit.Name = source["Name"].(string)
		esHit.CreatedDateTime = source["createdDateTime"].(string)
		esHit.LastModifiedDateTime = source["lastModifiedDateTime"].(string)
		esHit.WebUrl = source["webUrl"].(string)
		esHit.Source = source["source"].(string)

		// Print the structured data
		fmt.Printf("Elasticsearch Hit: %+v\n", esHit)
		eshits = append(eshits, esHit)
	}

	return eshits
}
