package cmd

import (
	"context"
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
)

const (
	flagPacketLossRate       = "packet-loss-rate"
	flagNetworkInterfaceName = "network-device-name"
	flagLogLevel             = "log-level"
	flagProductionMode       = "production-mode"
	flagBandwidth            = "bandwidth"
	flagLatency              = "latency"
	flagTcBinPath            = "tc-path"
)

var flagsStart struct {
	networkInterfaceName string
	packetLossRate       int32
	bandwidth            int64
	latency              int64
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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		/*---------*/

		if flagsStart.packetLossRate > 0 {
			pl := packetloss.PacketLoss{
				PacketLossRate:   flagsStart.packetLossRate,
				NetworkInterface: iface,
			}
			go pl.Start(ctx, logger)
		}

		/*---------*/

		if flagsStart.bandwidth > 0 {
			b := bandwidth.Bandwidth{
				Limit:            flagsStart.bandwidth,
				NetworkInterface: iface,
			}
			go b.Start(ctx, logger)
		}

		/*---------*/

		if flagsStart.latency > 0 {
			l := latency.Latency{
				Latency:          time.Duration(flagsStart.latency) * time.Millisecond,
				NetworkInterface: iface,
				TcBinPath:        flagsStart.tcBinPath,
			}
			go l.Start(ctx, logger)
		}

		/*---------*/

		// Handle interrupt signal (Ctrl+C)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGINT)

		<-signalChan
		logger.Info("Received interrupt signal. Shutting down...")
		cancel()

		return nil
	},
}
