package main

import (
	"./handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


func main() {

	router := mux.NewRouter()


	router.HandleFunc("/nodes", handlers.NodesHandler)
	router.HandleFunc("/nodes/{nodeId}/measurements", handlers.MeasurementsByNodeHandler)
	router.HandleFunc("/nodes/{nodeId}/api-credentials", handlers.GenerateApiCredentialsHandler)
	err := http.ListenAndServe(":8080", logRequest(router))
	if err != nil {
		log.Fatal(err)
	}
}



func logRequest(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w,r)
	})
}