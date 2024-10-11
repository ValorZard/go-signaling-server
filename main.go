package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pion/webrtc/v4"
)

var sessions = []webrtc.SessionDescription{}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/offer/get", offerGet)
	http.HandleFunc("/offer/post", offerPost)
	http.HandleFunc("/answer", answer)
	http.HandleFunc("/ice", ice)
	http.ListenAndServe(":8080", nil)
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

	sessions = append(sessions, sdp)

	fmt.Println(sessions)
}

func offerGet(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the offer handler
	w.Header().Set("Content-Type", "application/json")

}

func answer(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the answer handler
	w.Header().Set("Content-Type", "application/json")
}

func ice(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the ice handler
	w.Header().Set("Content-Type", "application/json")
}
