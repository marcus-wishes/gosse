package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var sseLogger = log.New(os.Stdout, "SSE - ", log.LstdFlags)

type sseClient struct {
	w    http.ResponseWriter
	cont *http.ResponseController
	r    *http.Request
}

// https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events
type sseMessage struct {
	Data  string        `json:"data"`
	ID    string        `json:"id"`
	Event string        `json:"event"`
	Retry time.Duration `json:"retry"`
}

func (msg sseMessage) ToString() string {
	ret := ""
	if msg.ID != "" {
		ret += fmt.Sprintf("id: %s\n", msg.ID)
	}
	if msg.Event != "" {
		ret += fmt.Sprintf("event: %s\n", msg.Event)
	}
	if msg.Retry != 0 {
		ret += fmt.Sprintf("retry: %d\n", msg.Retry)
	}
	ret += fmt.Sprintf("data: %s\n", msg.Data)
	return ret + "\n"
}

var sseClients = make(map[string]sseClient, 0)

// the heartBeat keeps the connection alive
func heartBeat() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		SendSSE(sseMessage{
			Data:  "ping",
			Event: "sse_heartbeat",
		})
	}
}

/* SendSSE sends a message to all connected clients
 */
func SendSSE(sseMsg sseMessage) {
	buf := sseMsg.ToString()
	for _, client := range sseClients {
		_, err := fmt.Fprint(client.w, buf)
		if err != nil {
			return
		}
		err = client.cont.Flush()
		if err != nil {
			return
		}
	}
}

// the sseHandle adds the connection to a connected clients map and keeps the connection alive
func sseHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != config.SSEClient.Method {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, exists := sseClients[r.RemoteAddr]; exists {
		delete(sseClients, r.RemoteAddr)
		fmt.Println("client already connected, reconnecting", r.RemoteAddr)
	}

	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rc := http.NewResponseController(w)
	sseClients[r.RemoteAddr] = sseClient{w, rc, r}
	sseLogger.Println("client connected", r.RemoteAddr)

	// Create a channel for client disconnection
	clientGone := r.Context().Done()
	<-clientGone
	sseLogger.Println("client disconnected", r.RemoteAddr)
	delete(sseClients, r.RemoteAddr)
}

func SEEClientServer() {
	go heartBeat()

	http.HandleFunc(config.SSEClient.Path, sseHandler)

	sseLogger.Printf("client server is running on :%s%s using %s\n", config.SSEClient.Port, config.SSEClient.Path, config.SSEClient.Method)
	if err := http.ListenAndServe(":"+config.SSEClient.Port, nil); err != nil {
		sseLogger.Println(err.Error())
	}
}
