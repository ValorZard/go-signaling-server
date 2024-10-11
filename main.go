package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pion/webrtc/v4"
	"github.com/rs/cors"
)

var offers = []webrtc.SessionDescription{}
var answers = []webrtc.SessionDescription{}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/offer/get", offerGet)
	mux.HandleFunc("/offer/post", offerPost)
	mux.HandleFunc("/answer/get", answerGet)
	mux.HandleFunc("/answer/post", answerPost)
	mux.HandleFunc("/ice", ice)

	fmt.Println("Server started on port 3000")
	// cors.Default() setup the middleware with default options being
    // all origins accepted with simple methods (GET, POST). See
    // documentation below for more options.
    handler := cors.Default().Handler(mux)
    http.ListenAndServe(":3000", handler)
}

func offerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get last element of offer slice
	if len(offers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Offer not found"))
		return
	}
	offer := offers[len(offers)-1]

	jsonValue, _ := json.Marshal(offer)

	io.Writer.Write(w, jsonValue)

	fmt.Println("offerGet")
	fmt.Println(jsonValue)
}

func offerPost(w http.ResponseWriter, r *http.Request) {
	var sdp webrtc.SessionDescription

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&sdp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offers = append(offers, sdp)

	fmt.Println("offerPost")
	fmt.Println(offers)
}

func answerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get last element of offer slice
	if len(answers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Answer not found"))
		return
	}

	answer := answers[len(answers)-1]

	jsonValue, _ := json.Marshal(answer)

	io.Writer.Write(w, jsonValue)

	fmt.Println("answerGet")
	fmt.Println(jsonValue)
}

func answerPost(w http.ResponseWriter, r *http.Request) {
	var sdp webrtc.SessionDescription

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&sdp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answers = append(answers, sdp)

	fmt.Println("answerPost")
	fmt.Println(answers)
}

func ice(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the ice handler
	w.Header().Set("Content-Type", "application/json")
}
