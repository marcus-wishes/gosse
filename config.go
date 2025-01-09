package main

type WebHookIncomingConfig struct {
	/* Method is the method used for incoming requests
	@default: POST
	*/
	Method string

	/* Path is the path where the service listens to incoming requests
	@default: ""
	*/
	Path string
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
	Port string
}

type Config struct {
	// Incoming WebHook Config
	WebHook []WebHookIncomingConfig

	/* WebHookPort is the port where the webhooks sends the events
	@default: 5000
	*/
	WebHookPort string

	// Clients SSE Config
	SSEClient SSEClientConnectConfig
}
