package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State stores the current state of blockchain
// It stores the balances of all individuals,
// a list of all transactions and a pointer to dbFile
type State struct {
	Balances   map[Account]uint
	txnMempool []Txn
	dbFile     *os.File
}

func NewStateFromDisk() (*State, error) {
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// forge the filepath and load data
	genesisFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genesisFilePath)
	if err != nil {
		return nil, err
	}

	// update balances
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// load txns from txn.db to update the blockchain
	txnDbPath := filepath.Join(cwd, "database", "txn.db")
	f, err := os.OpenFile(txnDbPath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Txn, 0), f}

	// iterate over the txns
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// loading json encoded txns into struct
		var txn Txn
		json.Unmarshal(scanner.Bytes(), &txn)
		if err := state.apply(txn); err != nil {
			return nil, err
		}
	}

	return state, nil
}

// Add adds a txn to the current state
func (s *State) Add(txn Txn) error {
	if err := s.apply(txn); err != nil {
		return err
	}
	s.txnMempool = append(s.txnMempool, txn)
	return nil
}

// Persist adds the transactions to disk
func (s *State) Persist() error {
	// creating a temp mempool as current mempool will be modified
	tempMemPool := make([]Txn, len(s.txnMempool))
	copy(tempMemPool, s.txnMempool)

	// loop over the mempool
	// add txn to db
	// and remove txn from mempool
	for i := 0; i < len(tempMemPool); i++ {
		txnJson, err := json.Marshal(tempMemPool[i])
		if err != nil {
			return err
		}

		if _, err := s.dbFile.Write(append(txnJson, '\n')); err != nil {
			return err
		}

		s.txnMempool = s.txnMempool[1:]
	}
	return nil
}

// apply validates the txn against the current State
// and makes changes to the State
func (s *State) apply(txn Txn) error {
	// check is txn is block reward
	if txn.IsReward() {
		s.Balances[txn.To] += txn.Value
		return nil
	}

	// check if account has enough funds
	if txn.Value > s.Balances[txn.From] {
		return fmt.Errorf("insufficient funds")
	}

	// complete txn
	s.Balances[txn.From] -= txn.Value
	s.Balances[txn.To] += txn.Value
	return nil
}

func (s *State) Close() {
	s.dbFile.Close()
}
