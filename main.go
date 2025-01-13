package main

import (
	"fmt"
)

func main() {
	// parse the command line arguments
	ParseConfigFromCommandLine()

	fmt.Printf("Starting WebHook server on port %d expecting method %s\n", config.WebHooks.Port, config.WebHooks.Method)

	// start the webhook server
	go WebHookServer()

	// start the sse server
	SEEClientServer()
}
