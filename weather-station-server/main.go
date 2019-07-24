package main

import (
	"./handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)


func main() {

	log.Printf("Init Server")
	router := mux.NewRouter()

	rand.Seed(time.Now().Unix())

	router.Path("/nodes").HandlerFunc(handlers.FetchNodesHandler).Methods(http.MethodGet)
	router.Path("/nodes/{nodeId}/measurements").HandlerFunc(handlers.PostMeasurementForNodeHandler).Methods(http.MethodPost)
	router.Path("/nodes/{nodeId}/measurements").HandlerFunc(handlers.GetLastMeasurementsByNodeHandler).Methods(http.MethodGet).Queries("limit", "{[0-9]*?}")
	router.Path("/nodes/{nodeId}/measurements").HandlerFunc(handlers.GetAllMeasurementsByNodeHandler).Methods(http.MethodGet)
	router.Path("/nodes/{nodeId}/api-token").HandlerFunc(handlers.GenerateApiCredentialsHandler).Methods(http.MethodGet)
	router.Path("/users").HandlerFunc(handlers.CreateUserHandler).Methods(http.MethodPost)
	router.Path("/users/login").HandlerFunc(handlers.AuthenticationHandler).Methods(http.MethodPost)
	router.Path("/users/enable").HandlerFunc(handlers.EnableUserHandler).Methods(http.MethodPost)
	router.Path("/users/me").HandlerFunc(handlers.UsersMe).Methods(http.MethodGet)



	log.Printf("Server started")
	log.Printf("You can access the Api at http://localhost:8080")

	err := http.ListenAndServe(":8080", corsHandler(logRequest(router)))

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


func corsHandler(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")


		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			handler.ServeHTTP(w,r)
		}
	})
}

