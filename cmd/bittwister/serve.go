package bittwister

import (
	"fmt"

	api "github.com/celestiaorg/bittwister/api/v1"
	"github.com/spf13/cobra"
)

const (
	flagServeAddr = "serve-addr"
)

var flagsServe struct {
	serveAddr      string
	originAllowed  string
	logLevel       string
	productionMode bool
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVar(&flagsServe.serveAddr, flagServeAddr, ":9007", "address to serve on")
	serveCmd.PersistentFlags().StringVar(&flagsServe.originAllowed, "origin-allowed", "*", "origin allowed for CORS")

	serveCmd.PersistentFlags().StringVar(&flagsServe.logLevel, flagLogLevel, "info", "log level (e.g. debug, info, warn, error, dpanic, panic, fatal)")
	serveCmd.PersistentFlags().BoolVar(&flagsServe.productionMode, flagProductionMode, false, "production mode (e.g. disable debug logs)")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the Bit Twister API server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := getLogger(flagsServe.logLevel, flagsServe.productionMode)
		if err != nil {
			return err
		}
		defer func() {
			// The error is ignored because of this issue: https://github.com/uber-go/zap/issues/328
			_ = logger.Sync()
		}()

		logger.Info("Starting the API server...")

		restAPI := api.NewRESTApiV1(flagsServe.productionMode, logger)
		logger.Fatal(fmt.Sprintf("REST API server: %v", restAPI.Serve(flagsServe.serveAddr, flagsServe.originAllowed)))

		return nil
	},
}
