# Bit Twister

Bit Twister: A CLI tool for precise network traffic shaping. Simulate latency, bandwith limitation, drop packets, impose jitter with ease. Enhance your network testing capabilities for development.

## Build

```bash
make all
```

## Usage

```bash
sudo ./bin/bittwister start [flags]

Flags:
  -b, --bandwidth int                bandwidth limit in bps (e.g. 1000 for 1Kbps)
  -h, --help                         help for start
  -j, --jitter int                   jitter in milliseconds (e.g. 10 for 10ms)
  -l, --latency int                  latency in milliseconds (e.g. 100 for 100ms)
      --log-level string             log level (e.g. debug, info, warn, error, dpanic, panic, fatal) (default "info")
  -d, --network-device-name string   network interface name
  -p, --packet-loss-rate int32       packet loss rate (e.g. 10 for 10% packet loss)
      --production-mode              production mode (e.g. disable debug logs)
      --tc-path string               path to tc binary (default "tc")
```

### Example

```bash
# Apply 25 percent packet loss to eth0
sudo ./bin/bittwister start -d eth0 -p 25
```

```bash
# Apply 1 Mbps bandwidth limit to eth0
sudo ./bin/bittwister start -d eth0 -b 1048576
```

```bash
# Apply 100 ms latency to eth0
sudo ./bin/bittwister start -d eth0 -l 100
```

```bash
# Apply 10 ms jitter to eth0
sudo ./bin/bittwister start -d eth0 -j 10
```

### Start the API server

```bash
sudo ./bin/bittwister serve [flags]

Flags:
  -h, --help                    help for serve
      --log-level string        log level (e.g. debug, info, warn, error, dpanic, panic, fatal) (default "info")
      --origin-allowed string   origin allowed for CORS (default "*")
      --production-mode         production mode (e.g. disable debug logs)
      --serve-addr string       address to serve on (default "localhost:9007")
```

### API Endpoints

Please note that all the endpoints have to be prefixed with `/api/v1`.

#### Packet Loss

- **Endpoint:** `/packetloss`
  - `/start`
    - **Method:** POST
    - **Data**: `{"network_interface":"eth0","packet_loss_rate":30}`
    - **Description:** Start packetloss service.
  - `/status`
    - **Method:** GET
    - **Description:** Get packetloss status.
  - `/stop`
    - **Method:** POST
    - **Description:** Stop packetloss service.

**example:**

```bash
curl -iX POST http://localhost:9007/api/v1/packetloss/start --data '{"network_interface":"eth0","packet_loss_rate":30}'
```

#### Bandwidth

- **Endpoint:** `/bandwidth`
  - `/start`
    - **Method:** POST
    - **Data**: `{"network_interface":"eth0","bandwidth":1048576}`
    - **Description:** Start bandwidth service.
  - `/status`
    - **Method:** GET
    - **Description:** Get bandwidth status.
  - `/stop`
    - **Method:** POST
    - **Description:** Stop bandwidth service.

#### Latency

- **Endpoint:** `/latency`
  - `/start`
    - **Method:** POST
    - **Data**: `{"network_interface":"eth0","latency":100,"jitter":10}`
    - **Description:** Start latency service.
  - `/status`
    - **Method:** GET
    - **Description:** Get latency status.
  - `/stop`
    - **Method:** POST
    - **Description:** Stop latency service.

#### Services

- **Endpoint:** `/services`
  - `/status`
    - **Method:** GET
    - **Description:** Get all network restriction services statuses and their configured parameters.

### SDK for Go

The BitTwister SDK for Go provides a convenient interface to interact with the BitTwister tool, which applies network restrictions on a network interface, including bandwidth limitation, packet loss, latency, and jitter.
More details about the SDK and how to use it can be found [here](./sdk/README.md).

### Using Bittwister in Kubernetes

To utilize Bittwister within a Kubernetes environment, specific configurations must be added to the container.

For simulating latency and jitter, the container needs to be granted additional capabilities. This can be achieved by adding the `NET_ADMIN` capability to the container's security context:

```yaml
securityContext:
  capabilities:
    add:
      - NET_ADMIN
```

For simulating packet loss and limiting bandwidth, the container needs to operate in privileged mode. This can be set in the container's security context as follows:

```yaml
securityContext:
  privileged: true
```

## Test

The tests require docker to be installed. To run all the tests, execute the following command:

```bash
make test
```

### Go unit tests

The Go unit tests can be run by executing the following command:

```bash
make test-go
```

**Note**: Root permission is required to run the unit tests. The tests are run on the loopback interface.

### Test Packet Loss

The packet loss function can be tested by running the following command:

```bash
make test-packetloss
```

```yaml
Results:
Number of packets per test: 50

Packet loss: 0%     Expected: 0%
Packet loss: 6%     Expected: 10%
Packet loss: 28%    Expected: 20%
Packet loss: 48%    Expected: 50%
Packet loss: 68%    Expected: 70%
Packet loss: 78%    Expected: 80%
Packet loss: +25    Expected: 100%
```

### Test Bandwidth Limitation

The bandwidth limitation function can be tested by running the following command:

```bash
make test-bandwidth
```

```yaml
Results:
Number of parallel connections per test: 128
Test duration per test: 600 seconds
expected bandwidth: 64 Kbps      actual bandwidth: 47 Kbps
expected bandwidth: 128 Kbps     actual bandwidth: 92 Kbps
expected bandwidth: 256 Kbps     actual bandwidth: 186 Kbps
expected bandwidth: 512 Kbps     actual bandwidth: 361 Kbps
expected bandwidth: 1024 Kbps    actual bandwidth: 0 bps
expected bandwidth: 2 Mbps       actual bandwidth: 1 Mbps
expected bandwidth: 4 Mbps       actual bandwidth: 2 Mbps
expected bandwidth: 8 Mbps       actual bandwidth: 6 Mbps
expected bandwidth: 16 Mbps      actual bandwidth: 13 Mbps
expected bandwidth: 32 Mbps      actual bandwidth: 24 Mbps
expected bandwidth: 64 Mbps      actual bandwidth: 49 Mbps
expected bandwidth: 128 Mbps     actual bandwidth: 100 Mbps
expected bandwidth: 256 Mbps     actual bandwidth: 187 Mbps
expected bandwidth: 512 Mbps     actual bandwidth: 459 Mbps
expected bandwidth: 1024 Mbps    actual bandwidth: 948 Mbps
```

### Test Latency

The latency function can be tested by running the following command:

```bash
make test-latency
```

```yaml
Results:
Number of packets per test: 10

Average latency: 11.203 ms       Expected: 10 ms
Average latency: 55.192 ms       Expected: 50 ms
Average latency: 110.250 ms      Expected: 100 ms
Average latency: 220.238 ms      Expected: 200 ms
Average latency: 550.342 ms      Expected: 500 ms
Average latency: 770.362 ms      Expected: 700 ms
Average latency: 1100.419 ms     Expected: 1000 ms
Average latency: 1700.255 ms     Expected: 1500 ms
Average latency: 2300.317 ms     Expected: 2000 ms
Average latency: 3600.249 ms     Expected: 3000 ms
Average latency: 5166.803 ms     Expected: 5000 ms
```

### Test Jitter

The jitter function can be tested by running the following command:

```bash
make test-jitter
```

```yaml
Number of packets per test: 50

Average jitter: 4.199 ms       Max expected jitter: 10 ms
Average jitter: 16.465 ms      Max expected jitter: 50 ms
Average jitter: 34.757 ms      Max expected jitter: 100 ms
Average jitter: 66.711 ms      Max expected jitter: 200 ms
Average jitter: 167.649 ms     Max expected jitter: 500 ms
Average jitter: 216.576 ms     Max expected jitter: 700 ms
Average jitter: 289.205 ms     Max expected jitter: 1000 ms
Average jitter: 466.368 ms     Max expected jitter: 1500 ms
Average jitter: 660.321 ms     Max expected jitter: 2000 ms
Average jitter: 669.862 ms     Max expected jitter: 3000 ms
Average jitter: 780.299 ms     Max expected jitter: 5000 ms
```
