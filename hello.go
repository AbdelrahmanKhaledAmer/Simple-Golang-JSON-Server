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
	mux.HandleFunc("/echo", serveAndLog(echo))
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

// "echo" route server
func echo(res http.ResponseWriter, req *http.Request) {
	// Only listen to POST requests
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(res, JSON{
			"message": "Only POST requests are allowed on this route.",
		})
		return
	}

	// Make sure the data sent is in a JSON format.
	data := JSON{}
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		writeJSON(res, JSON{
			"message": "I could not understand what you said because it wasn't written in a JSON format!",
		})
		return
	}

	// Colse body when out of scope.
	defer req.Body.Close()

	// Make sure the data sent is in "data"
	reqData, received := data["data"]
	if !received {
		res.WriteHeader(http.StatusBadRequest)
		writeJSON(res, JSON{
			"message": "No data received.",
		})
		return
	}
	res.WriteHeader(http.StatusOK)
	writeJSON(res, JSON{
		"message": reqData,
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
