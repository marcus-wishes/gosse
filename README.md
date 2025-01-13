# GoSSE

GoSSE is a simple REST interface that allows a server to send webhooks to notify connected clients. It supports both webhook and SSE (Server-Sent Events) configurations.

## Configuration

GoSSE can be configured using either command line flags or a configuration file in TOML format.

### Using Command Line Flags

You can configure GoSSE by passing the following flags:

- `-method`: Method for the webhook server (default: "POST")
- `-port`: Port for the webhook server (default: 5000)
- `-paths`: Paths for the webhook server (default: "ticket_created,ticket_updated,ticket_deleted")
- `-sse-path`: Path for the SSE client (default: "/events")
- `-sse-method`: Method for the SSE client (default: "GET")
- `-sse-port`: Port for the SSE client (default: 5001)

Example:
```sh
gosse -port=6000 -paths=ticket_created,ticket_updated,ticket_deleted -method=PUT -sse-path=/events -sse-method=GET -sse-port=5001
```

### Using a Configuration File

You can also configure GoSSE using a `config.toml` file. The file should have the following structure:

```toml
[WebHooks]
Method = "PUT"
Port = 6000
Paths = ["ticket_created", "ticket_updated", "ticket_deleted"]

[SSEClient]
Path = "/events"
Method = "GET"
Port = 5001
```

To use the configuration file, pass the `-config` flag with the path to the file:

```sh
gosse -config=config.toml
```

### Help

To display the help message, use the `--help` flag:

```sh
gosse --help
```

## Running GoSSE

You can run GoSSE using the provided configurations. The server will start and listen for incoming webhook requests and SSE client connections.

Example:
```sh
gosse -config=config.toml
```

Or using command line flags:
```sh
gosse -port=6000 -paths=ticket_created,ticket_updated,ticket_deleted -method=PUT -sse-path=/events -sse-method=GET -sse-port=5001
```

GoSSE will print the starting information and begin handling requests based on the provided configuration.
