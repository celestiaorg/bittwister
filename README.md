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
  -h, --help                         help for start
  -l, --log-level string             log level (e.g. debug, info, warn, error, dpanic, panic, fatal) (default "info")
  -d, --network-device-name string   network interface name
  -p, --packet-loss-rate string      packet loss rate (e.g. 10 for 10% packet loss) (default "0")
  -m, --production-mode              production mode (e.g. disable debug logs)

```

### Example

```bash
sudo ./bin/bittwister start -d eth0 -p 25
```
