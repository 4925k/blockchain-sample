package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Major  = "0"
	Minor  = "6"
	Fix    = "0"
	Verbal = "HTTP API"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "describes version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s.%s.%s %s\n", Major, Minor, Fix, Verbal)
	},
}
