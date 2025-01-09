package main

var testWebHookConfig = []WebHookIncomingConfig{
	{
		Method: "POST",
		Path:   "/ticket_created",
	},

	{
		Method: "POST",
		Path:   "/ticket_updated",
	},

	{
		Method: "POST",
		Path:   "/ticket_deleted",
	},
}

var testSEEClientConfig = SSEClientConnectConfig{
	Path:   "/events",
	Method: "GET",
	Port:   "5001",
}

var testConfig = Config{
	WebHook:     testWebHookConfig,
	WebHookPort: "5000",
	SSEClient:   testSEEClientConfig,
}

var config Config = testConfig

func main() {
	// start the webhook server
	go WebHookServer()

	// start the sse server
	SEEClientServer()
}
