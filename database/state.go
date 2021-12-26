package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"
)

// State stores the current state of blockchain
// It stores the balances of all individuals,
// a list of all transactions and a pointer to dbFile
type State struct {
	Balances        map[Account]uint
	txnMempool      []Txn
	dbFile          *os.File
	latestBlock     Block
	latestBlockHash Hash
	hasGenesisBlock bool
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

	f, err := os.OpenFile(getBlocksDbFilePath(path), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Txn, 0), f, Block{}, Hash{}, false}

	// iterate over the txns
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		if len(blockFsJson) == 0 {
			break
		}
		var blockFs BlockFs
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		err = applyTxns(blockFs.Value.Txns, state)
		if err != nil {
			return nil, err
		}

		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockFs.Key
		state.hasGenesisBlock = true
	}

	return state, nil
}

//adds collection of blocks to the current state
func (s *State) AddBlocks(blocks []Block) error {
	for _, b := range blocks {
		_, err := s.AddBlock(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add adds a block to the current state
func (s *State) AddBlock(b Block) (Hash, error) {
	pendingState := s.copy()

	err := applyBlock(b, pendingState)
	if err != nil {
		return Hash{}, err
	}

	blockHash, err := b.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFs{blockHash, b}
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}

	s.Balances = pendingState.Balances
	s.latestBlockHash = blockHash
	s.latestBlock = b
	s.hasGenesisBlock = true

	return blockHash, nil
}

// applyBlock adds all the txns in the block to the state
func applyBlock(b Block, s State) error {
	nextExpectedBlockNumber := s.latestBlock.Header.Number + 1

	if s.hasGenesisBlock && b.Header.Number != nextExpectedBlockNumber {
		return fmt.Errorf("next expected block must be %d, not %d", nextExpectedBlockNumber, b.Header.Number)
	}

	if s.hasGenesisBlock && s.latestBlock.Header.Number > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf("next block parent hash must be %x not %x", s.latestBlockHash, b.Header.Parent)
	}

	return applyTxns(b.Txns, &s)
}

// applyTxns completes the given transactions on the state
func applyTxns(txns []Txn, s *State) error {
	for _, txn := range txns {
		err := ApplyTxn(txn, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// ApplyTxn completes the given transaction on the state
func ApplyTxn(txn Txn, s *State) error {
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

// Latest Snapshot returns the latest snapshot of the current state
func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}
func (s *State) LatestBlock() Block {
	return s.latestBlock
}

// Close closes the db file
func (s *State) Close() error {
	return s.dbFile.Close()
}

// NextBlockNumber returns 0 if its the first block
// otherwise returns the next block number by incrementing the lastest block
func (s *State) NextBlockNumber() uint64 {
	if !s.hasGenesisBlock {
		return uint64(0)
	}
	return s.LatestBlock().Header.Number + 1
}

// copy copies the
func (s *State) copy() State {
	c := State{}
	c.hasGenesisBlock = s.hasGenesisBlock
	c.latestBlock = s.latestBlock
	c.latestBlockHash = s.latestBlockHash
	c.txnMempool = make([]Txn, len(s.txnMempool))
	c.Balances = make(map[Account]uint)

	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}

	c.txnMempool = append(c.txnMempool, s.txnMempool...)

	return c
}

// Persist adds the transactions to the block
func (s *State) Persist() (Hash, error) {
	latestBlockHash, err := s.latestBlock.Hash()
	if err != nil {
		return Hash{}, err
	}
	block := NewBlock(
		latestBlockHash,
		s.latestBlock.Header.Number+1,
		uint64(time.Now().Unix()),
		s.txnMempool,
	)
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

	s.latestBlockHash = latestBlockHash
	s.latestBlock = block
	s.txnMempool = []Txn{}

	return latestBlockHash, nil
}
