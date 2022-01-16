package database

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

// Block header stores parent block metadata
// []Txn stores tsns in the new block
type Block struct {
	Header BlockHeader `json:"header"`
	Txns   []Txn       `json:"payload"`
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Number uint64 `json:"number"`
	Time   uint64 `json:"time"`
}

type BlockFs struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

// NewBlock returns a Block including the given parameters
func NewBlock(parent Hash, number, time uint64, txns []Txn) Block {
	return Block{BlockHeader{parent, number, time}, txns}
}

// Hash returns the sha2356 hash of given blcok
func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, nil
	}
	return sha256.Sum256(blockJson), nil
}

func GetBlocksAfter(blockHash Hash, dataDir string) ([]Block, error) {
	// open file and load blockchain
	f, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	// create blocks to add newer blocks if present
	blocks := make([]Block, 0)
	newBlock := false

	// loop over blockchain
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return blocks, err
		}

		// unmarshal each block into blockFs
		var blockFs BlockFs
		err = json.Unmarshal(scanner.Bytes(), &blockFs)
		if err != nil {
			return blocks, err
		}

		// if newBlock, add to blocks
		if newBlock {
			blocks = append(blocks, blockFs.Value)
			continue
		}

		// if blockHash matches the blockFs.Key, the node is at current block
		// the following blocks are newer blocks
		if blockHash == blockFs.Key {
			newBlock = true
		}
	}

	return blocks, nil
}

func (h Hash) IsEmpty() bool {
	emptyHash := Hash{}

	return bytes.Equal(emptyHash[:], h[:])
}
