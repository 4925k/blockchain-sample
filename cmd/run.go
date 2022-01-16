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

			port, _ := cmd.Flags().GetUint64(flagPort)
			ip, _ := cmd.Flags().GetString(flagIP)

			bootstrap := node.NewPeerNode("40.71.208.186", 8080, true, true)
			n := node.New(dataDir, ip, port, *bootstrap)
			err := n.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	addDefaultRequiredFlags(runCMD)
	runCMD.Flags().Uint64(flagPort, node.DefaultHttpPort, "port to run the node on")
	runCMD.Flags().String(flagIP, node.DefaultIP, "ip to run the node on")
	return runCMD
}
