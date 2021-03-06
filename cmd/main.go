package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var flagDataDir = "dataDir"
var flagPort = "port"
var flagIP = "ip"

func main() {
	var paisaCMD = &cobra.Command{
		Use:   "paisa",
		Short: "nefoli blockchain",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	paisaCMD.AddCommand(runCmd())
	paisaCMD.AddCommand(versionCMD)
	paisaCMD.AddCommand(migrateCMD())
	paisaCMD.AddCommand(balancesCMD())

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
