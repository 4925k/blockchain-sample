package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var paisaCMD = &cobra.Command{
		Use:   "paisa",
		Short: "nefoli blockchain",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	paisaCMD.AddCommand(versionCMD)
	paisaCMD.AddCommand(balancesCMD())
	paisaCMD.AddCommand(txnCMD())

	err := paisaCMD.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
