package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var flagDataDir = "dataDir"

func main() {
	var paisaCMD = &cobra.Command{
		Use:   "paisa",
		Short: "nefoli blockchain",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	paisaCMD.AddCommand(versionCMD)
	paisaCMD.AddCommand(balancesCMD())
	paisaCMD.AddCommand(txnCMD())
	paisaCMD.AddCommand(runCmd())

	err := paisaCMD.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "absolute path where all data is stored")
	cmd.MarkFlagRequired(flagDataDir)
}
