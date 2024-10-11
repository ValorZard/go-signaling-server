package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pion/webrtc/v4"
)

var offers = []webrtc.SessionDescription{}
var answers = []webrtc.SessionDescription{}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/offer/get", offerGet)
	http.HandleFunc("/offer/post", offerPost)
	http.HandleFunc("/answer/get", answerGet)
	http.HandleFunc("/answer/post", answerPost)
	http.HandleFunc("/ice", ice)

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
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
