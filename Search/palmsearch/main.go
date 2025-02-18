package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"palmsearch/elasticsearch"
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

	fmt.Println("execution done")
	if err := http.ListenAndServe(":5555", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}

func sync(w http.ResponseWriter, r *http.Request) {
	fmt.Println("execution sync")

	sharepoint.IndexSharepoint()
}

func getAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("execution sync")

	elasticsearch.GetAll()
}
