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

## Test

The tests require docker to be installed. To run all the tests, execute the following command:
  
```bash
make test
```

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
Number of parallel connections per test: 100
Duration per test: 60 seconds
Bandwidth: 0 bps          Expected: 64 Kbps
Bandwidth: 95 Kbps        Expected: 128 Kbps
Bandwidth: 150 Kbps       Expected: 256 Kbps
Bandwidth: 191 Kbps       Expected: 512 Kbps
Bandwidth: 737 Kbps       Expected: 1024 Kbps
Bandwidth: 1 Mbps         Expected: 2 Mbps
Bandwidth: 3 Mbps         Expected: 4 Mbps
Bandwidth: 6 Mbps         Expected: 8 Mbps
Bandwidth: 13 Mbps        Expected: 16 Mbps
Bandwidth: 27 Mbps        Expected: 32 Mbps
Bandwidth: 50 Mbps        Expected: 64 Mbps
Bandwidth: 115 Mbps       Expected: 128 Mbps
Bandwidth: 221 Mbps       Expected: 256 Mbps
Bandwidth: 407 Mbps       Expected: 512 Mbps
Bandwidth: 934 Mbps       Expected: 1024 Mbps
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
