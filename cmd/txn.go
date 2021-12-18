package main

import (
	"blockchain-sample/database"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var flagFrom = "from"
var flagTo = "to"
var flagValue = "value"

func txnCMD() *cobra.Command {
	var txnsCMD = &cobra.Command{
		Use:   "txn",
		Short: "Interact with txns",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	txnsCMD.AddCommand(txnAddCMD())

	return txnsCMD
}

func txnAddCMD() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "Adds new txn to database",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			txn := database.NewTxn(fromAcc, toAcc, value, "")

			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer state.Close()

			err = state.Add(txn)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = state.Persist()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Txn added successfully to nefoli")
		},
	}
	cmd.Flags().String(flagFrom, "", "Sender account")
	cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "Reciever account")
	cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	cmd.MarkFlagRequired(flagValue)

	return cmd
}
