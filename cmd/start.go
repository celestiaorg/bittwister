package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celestiaorg/bittwister/xdp/bandwidth"
	"github.com/celestiaorg/bittwister/xdp/latency"
	"github.com/celestiaorg/bittwister/xdp/packetloss"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	flagPacketLossRate       = "packet-loss-rate"
	flagNetworkInterfaceName = "network-device-name"
	flagLogLevel             = "log-level"
	flagProductionMode       = "production-mode"
	flagBandwidth            = "bandwidth"
	flagLatency              = "latency"
	flagJitter               = "jitter"
	flagTcBinPath            = "tc-path"
)

var flagsStart struct {
	networkInterfaceName string
	packetLossRate       int32
	bandwidth            int64
	latency              int64
	jitter               int64
	tcBinPath            string

	logLevel       string
	productionMode bool
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().Int32VarP(&flagsStart.packetLossRate, flagPacketLossRate, "p", 0, "packet loss rate (e.g. 10 for 10% packet loss)")
	startCmd.PersistentFlags().StringVarP(&flagsStart.networkInterfaceName, flagNetworkInterfaceName, "d", "", "network interface name")
	startCmd.PersistentFlags().Int64VarP(&flagsStart.bandwidth, flagBandwidth, "b", 0, "bandwidth limit in bps (e.g. 1000 for 1Kbps)")
	startCmd.PersistentFlags().Int64VarP(&flagsStart.latency, flagLatency, "l", 0, "latency in milliseconds (e.g. 100 for 100ms)")
	startCmd.PersistentFlags().Int64VarP(&flagsStart.jitter, flagJitter, "j", 0, "jitter in milliseconds (e.g. 10 for 10ms)")
	startCmd.PersistentFlags().StringVar(&flagsStart.tcBinPath, flagTcBinPath, "tc", "path to tc binary")

	startCmd.PersistentFlags().StringVar(&flagsStart.logLevel, flagLogLevel, "info", "log level (e.g. debug, info, warn, error, dpanic, panic, fatal)")
	startCmd.PersistentFlags().BoolVar(&flagsStart.productionMode, flagProductionMode, false, "production mode (e.g. disable debug logs)")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the Bit Twister",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := getLogger(flagsStart.logLevel, flagsStart.productionMode)
		if err != nil {
			return err
		}
		defer func() {
			// The error is ignored because of this issue: https://github.com/uber-go/zap/issues/328
			_ = logger.Sync()
		}()

		logger.Info("Starting Bit Twister...")

		iface, err := net.InterfaceByName(flagsStart.networkInterfaceName)
		if err != nil {
			return fmt.Errorf("lookup network device %q: %v", flagsStart.networkInterfaceName, err)
		}

		/*---------*/

		if flagsStart.packetLossRate > 0 {
			pl := packetloss.PacketLoss{
				PacketLossRate:   flagsStart.packetLossRate,
				NetworkInterface: iface,
			}
			cancel, err := pl.Start()
			if err != nil {
				return err
			}
			logger.Info("Packetloss started", zap.Int32("rate (%)", flagsStart.packetLossRate), zap.String("device", flagsStart.networkInterfaceName))
			defer func() {
				if err := cancel(); err != nil {
					logger.Error("cancel packetloss", zap.Error(err))
				}
				logger.Info("Packetloss stopped", zap.String("device", flagsStart.networkInterfaceName))
			}()
		}

		/*---------*/

		if flagsStart.bandwidth > 0 {
			b := bandwidth.Bandwidth{
				Limit:            flagsStart.bandwidth,
				NetworkInterface: iface,
			}
			cancel, err := b.Start()
			if err != nil {
				return err
			}
			logger.Info("Bandwidth started", zap.Int64("limit (bps)", flagsStart.bandwidth), zap.String("device", flagsStart.networkInterfaceName))
			defer func() {
				if err := cancel(); err != nil {
					logger.Error("cancel bandwidth", zap.Error(err))
				}
				logger.Info("Bandwidth stopped", zap.String("device", flagsStart.networkInterfaceName))
			}()
		}

		/*---------*/

		if flagsStart.latency > 0 || flagsStart.jitter > 0 {
			l := latency.Latency{
				Latency:          time.Duration(flagsStart.latency) * time.Millisecond,
				Jitter:           time.Duration(flagsStart.jitter) * time.Millisecond,
				NetworkInterface: iface,
				TcBinPath:        flagsStart.tcBinPath,
			}
			cancel, err := l.Start()
			if err != nil {
				return err
			}
			logger.Info("Latency/Jitter started",
				zap.Int64("latency (ms)", l.Latency.Milliseconds()),
				zap.Int64("jitter (ms)", l.Jitter.Milliseconds()),
				zap.String("device", flagsStart.networkInterfaceName))
			defer func() {
				if err := cancel(); err != nil {
					logger.Error("cancel latency", zap.Error(err))
				}
				logger.Info("Latency/Jitter stopped", zap.String("device", flagsStart.networkInterfaceName))
			}()
		}

		/*---------*/

		// Handle interrupt signal (Ctrl+C)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGINT)

		<-signalChan
		logger.Info("Received interrupt signal. Shutting down...")

		return nil
	},
}
