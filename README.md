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
      --log-level string             log level (e.g. debug, info, warn, error, dpanic, panic, fatal) (default "info")
  -d, --network-device-name string   network interface name
  -p, --packet-loss-rate int32       packet loss rate (e.g. 10 for 10% packet loss)
      --production-mode              production mode (e.g. disable debug logs)

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

## Test

The tests require docker to be installed. To run all the tests, execute the following command:
  
```bash
make test
```

### Test Bandwidth Limitation

The bandwidth limitation function can be tested by running the following command:

```bash
make test-bandwidth
```

```yaml
Results:
Number of parallel connections per test: 100
Duration per test: 60 seconds
expected bandwidth: 64 Kbps 	actual bandwidth: 0 bps
expected bandwidth: 128 Kbps 	actual bandwidth: 95 Kbps
expected bandwidth: 256 Kbps 	actual bandwidth: 150 Kbps
expected bandwidth: 512 Kbps 	actual bandwidth: 191 Kbps
expected bandwidth: 1024 Kbps 	actual bandwidth: 737 Kbps
expected bandwidth: 2 Mbps 	    actual bandwidth: 1 Mbps
expected bandwidth: 4 Mbps 	    actual bandwidth: 3 Mbps
expected bandwidth: 8 Mbps 	    actual bandwidth: 6 Mbps
expected bandwidth: 16 Mbps 	actual bandwidth: 13 Mbps
expected bandwidth: 32 Mbps 	actual bandwidth: 27 Mbps
expected bandwidth: 64 Mbps 	actual bandwidth: 50 Mbps
expected bandwidth: 128 Mbps 	actual bandwidth: 115 Mbps
expected bandwidth: 256 Mbps 	actual bandwidth: 221 Mbps
expected bandwidth: 512 Mbps 	actual bandwidth: 407 Mbps
expected bandwidth: 1024 Mbps 	actual bandwidth: 934 Mbps
```

