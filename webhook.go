package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var whLogger = log.New(os.Stdout, "webhook - ", log.LstdFlags)

var incomingRequestBuffer = 100 // x incoming requests can be buffered
var incomingRequests = make(chan sseMessage, incomingRequestBuffer)

var retryDuration = time.Second * 5

func processRequests() {
	for req := range incomingRequests {
		whLogger.Printf("sending message of type %s", req.Event)
		SendSSE(req)
	}
	whLogger.Println("incoming requests channel closed")
}

func WebHookServer() {

	whLogger.Printf("listener is running on :%s\n", config.WebHookPort)

	for _, cfg := range config.WebHook {
		handleFunc := func(config WebHookIncomingConfig) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.Method != config.Method {
					http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
					return
				}

				messageData, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))

				event := config.Path
				if event[0] == '/' {
					event = event[1:]
				}

				incomingRequests <- sseMessage{Event: event, Data: string(messageData), Retry: retryDuration, ID: time.Now().String()}
			}
		}
		http.HandleFunc(cfg.Path, handleFunc(cfg))
		whLogger.Printf("path: %s\n", cfg.Path)
	}

	go processRequests()

	if err := http.ListenAndServe(":"+config.WebHookPort, nil); err != nil {
		fmt.Println(err.Error())
		close(incomingRequests)
	}
}
