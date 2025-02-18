package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func GetEsClient() *elasticsearch.Client {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
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
