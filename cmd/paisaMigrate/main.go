package main

import (
	"blockchain-sample/database"
	"fmt"
	"os"
	"time"
)

func main() {

	state, err := database.NewStateFromDisk()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer state.Close()

	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Txn{
			database.NewTxn("dibek", "dibek", 3, ""),
			database.NewTxn("dibek", "dibek", 700, "reward"),
		},
	)

	state.AddBlock(block0)
	block0hash, _ := state.Persist()

	block1 := database.NewBlock(
		block0hash,
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
	state.Persist()
}
