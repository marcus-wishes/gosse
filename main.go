package main

var testWebHookConfig = []WebHookIncomingConfig{
	{
		Method: "POST",
		Path:   "/webhook",
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
