package main

import (
	"blockchain-sample/node"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	var runCMD = &cobra.Command{
		Use:   "run",
		Short: "Launches the nefoli node and its HTTP API",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			fmt.Println("Launching nefole node and its HTTP API")

			err := node.Run(dataDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	addDefaultRequiredFlags(runCMD)

	return runCMD
}
