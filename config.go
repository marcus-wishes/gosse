package main

import (
	"flag"
	"os"
	"strings"
)

type WebHookConfig struct {
	/* Method is the method used for outgoing requests
	@default: POST
	*/
	Method string

	/* Port is the port where the service sends the events
	@default: 5000
	*/
	Port int

	/* Path is the path where the service sends the events
	@default: ["ticket_created", "ticket_updated", "ticket_deleted"]
	*/
	Paths []string
}

type SSEClientConnectConfig struct {
	/* PathClient is the path where the client connects to receive SSE events
	@default: /events
	*/
	Path string

	/* MethodClient is the method used for client connections to receive SSE events
	@default: GET
	*/
	Method string

	/* PortClient is the port where the client connects to receive SSE events
	@default: 5001
	*/
	Port int
}

type Config struct {
	WebHooks WebHookConfig

	SSEClient SSEClientConnectConfig
}

func paseSSEClientConfig(args []string) SSEClientConnectConfig {
	flagSet := flag.NewFlagSet("sse-client", flag.ExitOnError)

	path := flagSet.String("sse-path", "/events", "Path for the SSE client")
	method := flagSet.String("sse-method", "GET", "Method for the SSE client")
	port := flagSet.Int("sse-port", 5001, "Port for the SSE client")

	flagSet.Parse(args)

	return SSEClientConnectConfig{
		Path:   *path,
		Port:   *port,
		Method: *method,
	}
}

func parseWebhookConfig(args []string) (WebHookConfig, []string) {
	flagSet := flag.NewFlagSet("webhooks", flag.ExitOnError)

	method := flagSet.String("method", "POST", "Method for the webhook server")
	port := flagSet.Int("port", 5000, "Port for the webhook server")
	paths := flagSet.String("paths", "ticket_created,ticket_updated,ticket_deleted", "Paths for the webhook server")

	err := flagSet.Parse(args)
	if err != nil {
		panic(err)
	}

	return WebHookConfig{
		Method: *method,
		Port:   *port,
		Paths:  strings.Split(*paths, ","),
	}, args[:flagSet.NArg()]
}

var config Config

func ParseConfigFromCommandLine() {

	args := os.Args[1:]
	webhookConfig, args := parseWebhookConfig(args)
	sseClientConfig := paseSSEClientConfig(args)

	config = Config{
		WebHooks:  webhookConfig,
		SSEClient: sseClientConfig,
	}

}
