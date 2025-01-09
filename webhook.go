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
		whLogger.Println("start sending message to clients")
		SendSSE(req)
	}
	whLogger.Println("incoming requests channel closed")
}

func WebHookServer() {

	for _, config := range config.WebHook {
		http.HandleFunc(config.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != config.Method {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			messageData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			if len(messageData) == 0 {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			whLogger.Printf("received message: %s\n", messageData)

			incomingRequests <- sseMessage{Event: config.Path, Data: messageData, Retry: retryDuration, ID: time.Now().String()}
		})
		whLogger.Printf("path: %s\n", config.Path)
	}

	go processRequests()

	whLogger.Printf("listener is running on :%s\n", config.WebHookPort)
	if err := http.ListenAndServe(":"+config.WebHookPort, nil); err != nil {
		fmt.Println(err.Error())
		close(incomingRequests)
	}
}
