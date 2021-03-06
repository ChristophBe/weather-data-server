package main

import (
	"flag"
	"fmt"
	"github.com/ChristophBe/weather-data-server/config"
	"github.com/ChristophBe/weather-data-server/handlers"
	"github.com/ChristophBe/weather-data-server/handlers/httpHandler"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	configFilePtr := flag.String("config", "config.json", "Path to the Config File")
	flag.Parse()

	initializeConfiguration(configFilePtr)

	log.Printf("Init Server")

	router := mux.NewRouter()

	rand.Seed(time.Now().Unix())

	userHandlers := handlers.GetUserHandlers()
	nodeHandlers := handlers.GetNodeHandlers()
	measurementHandlers := handlers.GetMeasurementHandlers()

	router.Path("/nodes").Handler(nodeHandlers.GetFetchNodesHandler()).Methods(http.MethodGet)
	router.Path("/nodes").Handler(nodeHandlers.GetSaveNodeHandler()).Methods(http.MethodPost)
	router.Path("/nodes/{nodeId}/measurements").Handler(measurementHandlers.GetAddMeasurementHandler()).Methods(http.MethodPost)
	router.Path("/nodes/{nodeId}/measurements").Handler(measurementHandlers.GetMeasurementsByNodeHandler()).Methods(http.MethodGet)
	router.Path("/nodes/{nodeId}/api-token").Handler(nodeHandlers.GetNodeAuthTokenHandler()).Methods(http.MethodGet)
	router.Path("/nodes/{nodeId}/share").Handler(nodeHandlers.GetShareNodeHandler()).Methods(http.MethodPost)

	router.Path("/users").Handler(userHandlers.GetCreateUserHandler()).Methods(http.MethodPost)
	router.Path("/users/login").Handler(httpHandler.JsonHandler(handlers.UserAuthenticationHandler)).Methods(http.MethodPost)
	router.Path("/users/enable").Handler(userHandlers.GetUserEnableHandler()).Methods(http.MethodPost)
	router.Path("/users/me").Handler(userHandlers.GetUserMeHandler()).Methods(http.MethodGet)
	router.Path("/users/{userId}/nodes").Handler(nodeHandlers.GetFetchNodesByOwnerHandler()).Methods(http.MethodGet)

	conf, err := config.GetConfigManager().GetConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("can not load config; cause: %w", err))
	}

	port := ":8080"
	if conf.ServerPort != 0 {
		port = fmt.Sprintf(":%d", conf.ServerPort)
	}

	log.Printf("Server started")
	log.Printf("You can access the Api at http://localhost%s", port)

	err = http.ListenAndServe(port, corsHandler(logRequest(router)))

	if err != nil {
		err = fmt.Errorf("can not start server cause: %w", err)
		log.Fatal(err)
	}
}

func logRequest(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
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
			handler.ServeHTTP(w, r)
		}
	})
}

func initializeConfiguration(configFilePtr *string) *config.Configuration {
	configManager := config.GetConfigManager()
	err := configManager.LoadConfig(*configFilePtr)
	if err != nil {
		log.Fatal("Can not load configuration File", err)
	}
	conf, err := configManager.GetConfig()
	if err != nil {
		log.Fatal("Failed to Read Config.", err)
	}
	return conf
}
