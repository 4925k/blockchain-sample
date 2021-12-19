package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// State stores the current state of blockchain
// It stores the balances of all individuals,
// a list of all transactions and a pointer to dbFile
type State struct {
	Balances        map[Account]uint
	txnMempool      []Txn
	dbFile          *os.File
	latestBlockHash Hash
}

func NewStateFromDisk(path string) (*State, error) {
	// get current working directory
	err := initDataDirIfNotExists(path)
	if err != nil {
		return nil, err
	}

	// forge the filepath and load data
	gen, err := loadGenesis(getGenesisJsonFilePath(path))
	if err != nil {
		return nil, err
	}

	// update balances
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// load txns from txn.db to update the blockchain
	f, err := os.OpenFile(getBlocksDbFilePath(path), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Txn, 0), f, Hash{}}

	// iterate over the txns
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockFs BlockFs
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		if err := state.applyBlock(blockFs.Value); err != nil {
			return nil, err
		}

		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}

// Add adds a block to the current state
func (s *State) AddBlock(b Block) error {
	for _, txn := range b.Txns {
		fmt.Println(txn)
		if err := s.AddTxn(txn); err != nil {
			return err
		}
	}
	return nil
}

// applyBlock adds all the txns in the block to the state
func (s *State) applyBlock(b Block) error {
	for _, txn := range b.Txns {
		if err := s.apply(txn); err != nil {
			return err
		}
	}
	return nil
}

// AddTxn processes the given txn
func (s *State) AddTxn(txn Txn) error {
	if err := s.apply(txn); err != nil {
		return err
	}

	s.txnMempool = append(s.txnMempool, txn)
	return nil
}

// apply validates the txn and makes changes to the State
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

// Persist adds the transactions to the block
func (s *State) Persist() (Hash, error) {
	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txnMempool)
	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFs{blockHash, block}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new block to disk:\n%s\n", blockFsJson)

	if _, err = s.dbFile.Write(append(blockFsJson, '\n')); err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash
	s.txnMempool = []Txn{}

	return s.latestBlockHash, nil
}

// Latest Snapshot returns the latest snapshot of the current state
func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

// Close closes the db file
func (s *State) Close() error {
	return s.dbFile.Close()
}
