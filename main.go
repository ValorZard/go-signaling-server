package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/pion/webrtc/v4"
	"github.com/rs/cors"
)

var offers = []webrtc.SessionDescription{}
var answers = []webrtc.SessionDescription{}

type ClientConnection struct {
	IsHost bool
	Offer webrtc.SessionDescription
	Answer webrtc.SessionDescription
}

type Lobby struct {
	Clients []ClientConnection
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateNewLobbyId() string {
	// have random size for lobby id
	size := 6
    buffer := make([]rune, size)
    for i := range buffer {
        buffer[i] = letters[rand.Intn(len(letters))]
    }
    id := string(buffer)

	// check if room id is already in lobby_list
	_, ok := lobby_list[id]
	if ok {
		// if it already exists, call function again
		return generateNewLobbyId()
	}
	return id
}

func makeLobby() string {
	lobby := Lobby{}
	lobby.Clients = []ClientConnection{}
	// first client is always host
	lobby.Clients = append(lobby.Clients, ClientConnection{IsHost: true})
	lobby_id := generateNewLobbyId()
	lobby_list[lobby_id] = lobby
	return lobby_id
}

var lobby_list = map[string]Lobby{}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/lobby/host", lobbyHost)
	mux.HandleFunc("/lobby/join", lobbyJoin)
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

func lobbyHost(w http.ResponseWriter, r *http.Request) {
	lobby_id := makeLobby()
	// return lobby id to host
	io.Writer.Write(w, []byte(lobby_id))
	fmt.Println("lobbyHost")
	fmt.Println(lobby_id)
	fmt.Println(lobby_list)
}

// call "/lobby?id={lobby_id}" to connect to lobby
func lobbyJoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// https://freshman.tech/snippets/go/extract-url-query-params/
	// get lobby id from query params
	lobby_id := r.URL.Query().Get("id")
	fmt.Printf("lobby_id: %s", lobby_id)

	// only continue with connection if lobby exists
	_, ok := lobby_list[lobby_id]
	// If the key doesn't exist, return error
	if !ok {
    	w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Lobby not found"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %s", err)
		return
	}

	fmt.Printf("body: %s", body)
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
