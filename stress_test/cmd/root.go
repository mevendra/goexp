package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"stress_test/internal"
)

var rootCmd = &cobra.Command{
	Use:   "[-u url -r requests -c concurrency]",
	Short: "Execute stress test with given parameters",
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			panic(err)
		}

		requests, err := cmd.Flags().GetInt("requests")
		if err != nil {
			panic(err)
		}

		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			panic(err)
		}

		internal.Run(url, requests, concurrency)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "url to test")
	if err := rootCmd.MarkFlagRequired("url"); err != nil {
		panic(err)
	}

	rootCmd.Flags().IntP("requests", "r", 0, "number of requests to test")
	if err := rootCmd.MarkFlagRequired("requests"); err != nil {
		panic(err)
	}

	rootCmd.Flags().IntP("concurrency", "c", 0, "number of concurrent requests to test")
	if err := rootCmd.MarkFlagRequired("concurrency"); err != nil {
		panic(err)
	}
}
