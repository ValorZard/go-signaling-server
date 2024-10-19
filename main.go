package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/pion/webrtc/v4"
	"github.com/rs/cors"
)

type ClientConnection struct {
	IsHost bool
	Offer webrtc.SessionDescription
	Answer webrtc.SessionDescription
}

type Lobby struct {
	// host is first client in lobby.Clients
	Clients []ClientConnection
}

type PlayerData struct {
	// player id is index in lobby.Clients
	Id int
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
	fmt.Println("lobbyJoin")
	w.Header().Set("Content-Type", "application/json")
	// https://freshman.tech/snippets/go/extract-url-query-params/
	// get lobby id from query params
	lobby_id := r.URL.Query().Get("id")
	fmt.Printf("lobby_id: %s\n", lobby_id)

	// only continue with connection if lobby exists
	lobby, ok := lobby_list[lobby_id]
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

	// send player id once generated
	player_id := len(lobby.Clients)
	fmt.Printf("player_id: %d\n", player_id)
	lobby.Clients = append(lobby.Clients, ClientConnection{IsHost: false})
	fmt.Println(lobby.Clients)
	player_data := PlayerData{Id: player_id}
	jsonValue, _ := json.Marshal(player_data)
	io.Writer.Write(w, jsonValue)
}

func validatePlayer(w http.ResponseWriter, r *http.Request) (string, int, error) {
	fmt.Println("validatePlayer")
	lobby_id := r.URL.Query().Get("lobby_id")
	//fmt.Printf("lobby_id: %s\n", lobby_id)

	// only continue with connection if lobby exists
	_, ok := lobby_list[lobby_id]
	// If the key doesn't exist, return error
	if !ok {
    	w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Lobby not found"))
		return "", 0, errors.New("Lobby not found")
	}

	player_id_string := r.URL.Query().Get("player_id")
	player_id, err := strconv.Atoi(player_id_string)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Player not found"))
		return "", 0, errors.New("Player not found")
    }
	//fmt.Printf("player_id: %d\n", player_id)
	//fmt.Printf("length of lobby.Clients: %d\n", len(lobby_list[lobby_id].Clients))
	fmt.Println(lobby_list[lobby_id].Clients)
	// check if player actually exists
	if player_id < 0 || player_id >= len(lobby_list[lobby_id].Clients) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Player not found"))
		return "", 0, errors.New("Player not found")
	}
	return lobby_id, player_id, nil
}

func offerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lobby_id, player_id, err := validatePlayer(w, r)
	if err != nil {
		return
	}

	offer := lobby_list[lobby_id].Clients[player_id].Offer

	jsonValue, _ := json.Marshal(offer)

	io.Writer.Write(w, jsonValue)

	/*
	fmt.Println("offerGet")
	fmt.Println(jsonValue)
	*/
}

func offerPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("offerPost")
	
	lobby_id, player_id, err := validatePlayer(w, r)
	if err != nil {
		return
	}

	var sdp webrtc.SessionDescription

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&sdp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby_list[lobby_id].Clients[player_id].Offer = sdp
	fmt.Println(lobby_list[lobby_id])
}

func answerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lobby_id, player_id, err := validatePlayer(w, r)
	if err != nil {
		return
	}

	answer := lobby_list[lobby_id].Clients[player_id].Answer

	jsonValue, _ := json.Marshal(answer)

	io.Writer.Write(w, jsonValue)
}

func answerPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lobby_id, player_id, err := validatePlayer(w, r)
	if err != nil {
		return
	}

	var sdp webrtc.SessionDescription

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&sdp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby_list[lobby_id].Clients[player_id].Answer = sdp
	fmt.Println(lobby_list[lobby_id])
}

func ice(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the ice handler
	w.Header().Set("Content-Type", "application/json")
}
