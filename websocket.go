package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/coder/websocket"
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		log.Printf("websocket initialization error: %s", err)
		return
	}

	log.Printf("websocket connection accepted")

	// Do something with conn C.
	defer conn.CloseNow()

	if conn.Subprotocol() != "echo" {
		conn.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
		return
	}

	limiter := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = echo(conn, limiter)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			fmt.Printf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	// Debugging, read from reader and print out what we're echoing
	buf := new(strings.Builder)
	_, err = io.Copy(buf, r)
	// check errors
	if err != nil {
		return err
	}
	log.Println(buf.String())

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}