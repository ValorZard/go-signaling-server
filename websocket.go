package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("websocket initialization error: %s", err)
		return
	}

	log.Printf("websocket connection accepted")

	// Do something with conn C.
	defer conn.CloseNow()

	ctx := context.Background()

    for {
        typ, r, err := conn.Reader(ctx)
        if err != nil {
            log.Println("Failed to get reader:", err)
            return
        }

        w, err := conn.Writer(ctx, typ)
        if err != nil {
            log.Println("Failed to get writer:", err)
            return
        }

        buf := new(bytes.Buffer)
        buf.ReadFrom(r)
        log.Println(buf.String())

        _, err = io.Copy(w, r)
        if err != nil {
            log.Println("Failed to io.Copy:", err)
            return
        }

        err = w.Close()
        if err != nil {
            log.Println("Failed to close writer:", err)
            return
        }
    }
}