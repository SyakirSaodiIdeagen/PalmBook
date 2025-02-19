package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"palmsearch/elasticsearch"
	"palmsearch/googledrive"
	"palmsearch/sharepoint"
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

	fmt.Println("execution done")
	if err := http.ListenAndServe(":5555", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}

func sync(w http.ResponseWriter, r *http.Request) {
	fmt.Println("execution sync")
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
	fmt.Println("execution sync")
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
}
