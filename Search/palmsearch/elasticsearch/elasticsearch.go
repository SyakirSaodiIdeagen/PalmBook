package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func GetEsClient() *elasticsearch.Client {

	//caCert, err := os.ReadFile("./elasticsearch/http_ca.crt") // Replace with your .crt file
	//if err != nil {
	//	log.Fatalf("Error reading CA certificate: %s", err)
	//}
	//
	//// Create a CA certificate pool and add the CA certificate
	//caCertPool := x509.NewCertPool()
	//if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
	//	log.Fatalf("Failed to append CA certificate to pool")
	//}

	// TLS configuration
	//tlsConfig := &tls.Config{
	//	RootCAs:    caCertPool, // Trust Elasticsearch's certificate
	//	MinVersion: tls.VersionTLS12,
	//}

	// Elasticsearch client configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200", // Ensure HTTPS
		},
		Username: "elastic",    // Optional: Basic auth if enabled
		Password: "password1!", // Optional: Basic auth if enabled
		//Transport: &http.Transport{
		//	TLSClientConfig: tlsConfig,
		//},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
	// Call bulk insertion
}

func BulkInsert(es *elasticsearch.Client, buf bytes.Buffer) {
	//docs := []map[string]interface{}{
	//	{"id": 1, "user": "alice", "message": "First message", "timestamp": "2025-02-02"},
	//	{"id": 2, "user": "bob", "message": "Second message", "timestamp": "2025-02-03"},
	//	{"id": 3, "user": "carol", "message": "Third message", "timestamp": "2025-02-04"},
	//}
	//
	//var buf bytes.Buffer
	//
	//for _, doc := range data.Value {
	//	// Action metadata (index operation)
	//	meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "golang-bulk-index", "_id" : "%d" } }%s`, doc.Id, "\n"))
	//	// Document body
	//	data, err := json.Marshal(doc)
	//	if err != nil {
	//		log.Fatalf("Cannot encode document %d: %s", doc.Id, err)
	//	}
	//
	//	buf.Grow(len(meta) + len(data) + 1)
	//	buf.Write(meta)
	//	buf.Write(data)
	//	buf.WriteByte('\n')
	//}

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
