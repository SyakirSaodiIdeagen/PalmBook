package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"net/http"
)

type BulkDeleteRequest struct {
	Delete struct {
		Index string `json:"_index"`
		ID    string `json:"_id"`
	} `json:"delete"`
}

func GetEsClient() *elasticsearch.Client {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Username: "elastic",
		Password: "password1!",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
}

func BulkInsert(es *elasticsearch.Client, buf bytes.Buffer) {
	if buf.Len() <= 0 {
		return
	}
	res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex("golang-bulk-index"), es.Bulk.WithRefresh("true"))
	if err != nil {
		log.Fatalf("Failure indexing batch: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Bulk request failed: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response: %s", err)
	}

	if result["errors"].(bool) {
		for _, item := range result["items"].([]interface{}) {
			for action, info := range item.(map[string]interface{}) {
				if status := info.(map[string]interface{})["status"].(float64); status >= 400 {
					fmt.Printf("Error [%s]: %v\n", action, info)
				}
			}
		}
	} else {
		fmt.Println("Bulk insert successful!")
	}
}

func GetAll() []string {
	client := GetEsClient()

	searchQuery := `{
		"_source": false, 
		"fields": ["_id"],
		"query": {
			"match_all": {}
		}
	}`

	req := esapi.SearchRequest{
		Body:  bytes.NewReader([]byte(searchQuery)),
		Index: []string{"golang-bulk-index3"},
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error executing the search request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		log.Println("Not Found")
		return nil
	}
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the search response: %s", err)
	}
	alldocuments := []string{}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		doc := hit.(map[string]interface{})
		id := doc["_id"].(string)
		alldocuments = append(alldocuments, id)
		fmt.Println("Document ID:", id)
	}

	return alldocuments
}

func DeleteDocumentsBulk(ids []string) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	var buf bytes.Buffer

	if len(ids) == 0 {
		return
	}

	for _, id := range ids {
		req := BulkDeleteRequest{
			Delete: struct {
				Index string `json:"_index"`
				ID    string `json:"_id"`
			}{
				Index: "golang-bulk-index3",
				ID:    id,
			},
		}
		if err := json.NewEncoder(&buf).Encode(req); err != nil {
			log.Fatalf("Error encoding bulk delete request: %s", err)
		}
	}

	res, err := es.Bulk(&buf)
	if err != nil {
		log.Fatalf("Error executing bulk delete: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error deleting documents: %s", res.String())
	} else {
		fmt.Println("Successfully deleted documents")
	}
}
