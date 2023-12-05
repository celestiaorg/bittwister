# BitTwister SDK for Go

The BitTwister SDK for Go provides a convenient interface to interact with the BitTwister tool, which applies network restrictions on a network interface, including bandwidth limitation, packet loss, latency, and jitter.

## Installation

To use this SDK, import it into your Go project:

```bash
go get -u github.com/celestiaorg/bittwister/sdk
```

## Usage

### Initialization

First, import the SDK into your Go project:

```go
import "github.com/celestiaorg/bittwister/sdk"
```

Next, create a client by specifying the base URL of where BitTwister is running e.g. _a container running on the same machine_:

```go
func main() {
  baseURL := "http://localhost:9007/api/v1"
  client := sdk.NewClient(baseURL)
  // Use the client for API requests
}
```

### Examples

Bandwidth Service

Start the Bandwidth service with a specified network interface and bandwidth limit:

```go
req := sdk.BandwidthStartRequest{
    NetworkInterfaceName: "eth0",
    Limit:                100,
}

err := client.BandwidthStart(req)
if err != nil {
    // Handle error
}
```

Stop the Bandwidth service:

```go
err := client.BandwidthStop()
if err != nil {
    // Handle error
}
```

Retrieve the status of the Bandwidth service:
  
```go
status, err := client.BandwidthStatus()
if err != nil {
    // Handle error
}
// Use status for further processing
```

Similarly, you can use PacketlossStart, PacketlossStop, PacketlossStatus, LatencyStart, LatencyStop, LatencyStatus, and other functions provided by the SDK following similar usage patterns.
