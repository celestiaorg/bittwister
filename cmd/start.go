package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/celestiaorg/bittwister/xdp/packetloss"
	"github.com/spf13/cobra"
)

const (
	flagPacketLossRate       = "packet-loss-rate"
	flagNetworkInterfaceName = "network-device-name"
	flagLogLevel             = "log-level"
	flagProductionMode       = "production-mode"
)

var flagsStart struct {
	packetLossRate       string
	networkInterfaceName string
	logLevel             string
	productionMode       bool
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVarP(&flagsStart.packetLossRate, flagPacketLossRate, "p", "0", "packet loss rate (e.g. 10 for 10% packet loss)")
	startCmd.PersistentFlags().StringVarP(&flagsStart.networkInterfaceName, flagNetworkInterfaceName, "d", "", "network interface name")
	startCmd.PersistentFlags().StringVarP(&flagsStart.logLevel, flagLogLevel, "l", "info", "log level (e.g. debug, info, warn, error, dpanic, panic, fatal)")
	startCmd.PersistentFlags().BoolVarP(&flagsStart.productionMode, flagProductionMode, "m", false, "production mode (e.g. disable debug logs)")
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

		lossRate, err := strconv.ParseInt(flagsStart.packetLossRate, 10, 32)
		if err != nil {
			return fmt.Errorf("parse packet loss rate %q: %v", flagsStart.packetLossRate, err)
		}

		if lossRate > 0 {

			pl := packetloss.PacketLoss{
				PacketLossRate:   int32(lossRate),
				NetworkInterface: iface,
			}

			go pl.Start(ctx, logger)
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
