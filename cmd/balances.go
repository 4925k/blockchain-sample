package main

import (
	"blockchain-sample/database"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func balancesCMD() *cobra.Command {
	var balancesCMD = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balancesCMD.AddCommand(balancesListCMD)

	return balancesCMD
}

var balancesListCMD = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer state.Close()

		fmt.Println("Accounts Balances:")
		for account, balance := range state.Balances {
			fmt.Printf("%s: %d", account, balance)
		}
	},
}
