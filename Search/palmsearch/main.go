package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"palmsearch/elasticsearch"
	"palmsearch/googledrive"
	"palmsearch/sharepoint"
	"strings"
)

func main() {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow specific origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})
	handler := c.Handler(r)
	r.Handle("/sync", http.HandlerFunc(sync)).Methods("POST")
	r.Handle("/getAll", http.HandlerFunc(getAll)).Methods("GET")
	r.Handle("/search", http.HandlerFunc(search)).Methods("GET")
	r.Handle("/googledrive", http.HandlerFunc(runGd)).Methods("GET")
	//syncInit()

	if err := http.ListenAndServe(":5555", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}

func syncInit() {
	fmt.Println("sync executed")
	googledrive.IndexGoogleDrive()
	sharepoint.IndexSharepoint()
	fmt.Println("sync completed")

}

func sync(w http.ResponseWriter, r *http.Request) {
	googledrive.IndexGoogleDrive()
	sharepoint.IndexSharepoint()
}

func getAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("execution sync")

	elasticsearch.GetAll()
}

func runGd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("execution sync")

	googledrive.IndexGoogleDrive()
}

func search(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	payload, err := decodeJWTPayload(authHeader)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Username does not match", http.StatusUnauthorized)
		return
	}
	username := payload["preferred_username"].(string)
	usernamesplit := strings.Split(username, "@")

	userTenantName := usernamesplit[1]
	tenantName := os.Getenv("TENANT_NAME")

	if userTenantName == tenantName {
		queryParam := r.URL.Query().Get("query")
		searchHits := elasticsearch.Search(queryParam)

		responseData, err := json.Marshal(searchHits)
		if err != nil {
			http.Error(w, "Error marshaling response", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(responseData)
		if err != nil {
			return
		}
	} else {
		http.Error(w, "Username does not match", http.StatusUnauthorized)
		return
	}
}

func decodeJWTPayload(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return payload, nil
}
