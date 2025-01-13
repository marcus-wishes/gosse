package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
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

var helpMessage = `
Usage: 
gosse -config <path-to-config-toml-file> OR
gosse <flags>

Flags:
  -method string
    	Method for the webhook server (default "POST")
  -port int
    	Port for the webhook server (default 5000)
  -paths string
    	Paths for the webhook server (default "ticket_created,ticket_updated,ticket_deleted")
  -sse-path string
    	Path for the SSE client (default "/events")
  -sse-method string
    	Method for the SSE client (default "GET")
  -sse-port int
    	Port for the SSE client (default 5001)
`

func parseSSEClientConfig(args []string) SSEClientConnectConfig {
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
	if len(os.Args) == 2 {
		if (os.Args[1] == "-h") || (os.Args[1] == "--help") {
			fmt.Println(
				helpMessage,
			)
			os.Exit(0)
		}
		configFile := flag.String("config", "", "Path to the config.toml file")
		flag.Parse()

		var err error
		config, err = parseConfigFromToml(*configFile)
		if err != nil {
			fmt.Println("Failed to parse config file:", err)
			os.Exit(1)
		}
	} else {

		args := os.Args[1:]
		webhookConfig, args := parseWebhookConfig(args)
		sseClientConfig := parseSSEClientConfig(args)

		config = Config{
			WebHooks:  webhookConfig,
			SSEClient: sseClientConfig,
		}
	}
}

func parseConfigFromToml(filePath string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
