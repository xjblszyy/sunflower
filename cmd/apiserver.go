package cmd

import (
	"sunflower/cmd/apiserver"
	"sunflower/config"
	metricsMlwr "sunflower/pkg/middleware/metrics"

	"github.com/spf13/cobra"
)

func apiServerCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use: "apiserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			setups := []func() error{
				setupLogger,
				setupDB,
			}

			for _, setup := range setups {
				if err := setup(); err != nil {
					panic(err)
				}
			}

			var metrics *metricsMlwr.Prometheus = metricsMlwr.NewPrometheus(config.C.ServiceName, nil, nil)

			apiserver.RunServer(config.C, metrics)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			return nil
		},
	}

	serverCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	return serverCmd
}

func init() {
	RootCmd.AddCommand(apiServerCmd())
}
