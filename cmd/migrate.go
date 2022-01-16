package main

import (
	"blockchain-sample/database"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var migrateCMD = func() *cobra.Command {
	var migrateCMD = &cobra.Command{
		Use:   "migrate",
		Short: "Migrates the blockchain databse according to new business rule.",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer state.Close()

			block0 := database.NewBlock(
				database.Hash{},
				0,
				uint64(time.Now().Unix()),
				[]database.Txn{
					database.NewTxn("dibek", "dibek", 3, ""),
					database.NewTxn("dibek", "dibek", 700, "reward"),
				},
			)

			state.AddBlock(block0)
			block0hash, err := state.Persist()
			if err != nil {
				fmt.Println(err)
				return
			}

			block1 := database.NewBlock(
				block0hash,
				1,
				uint64(time.Now().Unix()),
				[]database.Txn{
					database.NewTxn("dibek", "nishan", 2000, ""),
					database.NewTxn("dibek", "dibek", 100, "reward"),
					database.NewTxn("nishan", "dibek", 1, ""),
					database.NewTxn("nishan", "sasim", 1000, ""),
					database.NewTxn("nishan", "dibek", 50, ""),
					database.NewTxn("dibek", "dibek", 600, "reward"),
				},
			)

			state.AddBlock(block1)
			block1hash, err := state.Persist()
			if err != nil {
				fmt.Println(err)
				return
			}

			block2 := database.NewBlock(
				block1hash,
				2,
				uint64(time.Now().Unix()),
				[]database.Txn{
					database.NewTxn("dibek", "dibek", 2400, "reward"),
				},
			)

			state.AddBlock(block2)
			_, err = state.Persist()
			if err != nil {
				fmt.Println(err)
				return
			}

		},
	}
	addDefaultRequiredFlags(migrateCMD)

	return migrateCMD
}
