package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var envPath string

var rootCmd = &cobra.Command{
	Use:   "sso",
	Short: "golang sso app.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&envPath, "env", "", "--env=/path/to/.env")
}
