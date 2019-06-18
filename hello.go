package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	cors "github.com/heppu/simple-cors"
)

type (
	// JSON models a json for sending and recieving in requests and responses
	JSON map[string]interface{}
	// Session models the session of a user
	Session map[string]interface{}
)

func main() {
	// Get port from env, if it doesn't exist, set to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting server on port", port)

	mux := http.NewServeMux()
	// Route handlers
	mux.HandleFunc("/hello", serveAndLog(helloWorld))
	mux.HandleFunc("/", serveAndLog(defaultServe))
	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, cors.CORS(mux)))
}

// Default route handler.
func defaultServe(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	writeJSON(res, JSON{
		"message": "Try a different route",
	})
}

// "hello" route handler
func helloWorld(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	writeJSON(res, JSON{
		"message": "Hello World",
	})
}

// Write JSON data to response
func writeJSON(res http.ResponseWriter, data JSON) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(res).Encode(data)
}

// Intermediary function that logs the current request and the status code attached to the response.
func serveAndLog(handler http.HandlerFunc) http.HandlerFunc {
	return func(inRes http.ResponseWriter, req *http.Request) {
		res := httptest.NewRecorder()
		handler(res, req)
		log.Printf("[%d] %-4s %s\n", res.Code, req.Method, req.URL.Path)

		for k, v := range res.HeaderMap {
			inRes.Header()[k] = v
		}
		inRes.WriteHeader(res.Code)
		res.Body.WriteTo(inRes)
	}
}
