package node

import (
	"blockchain-sample/database"
	"net/http"
	"time"
)

// statusHandler responds with the latest block hash and height
func statusHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	res := StatusRes{node.state.LatestBlockHash(), node.state.LatestBlock().Header.Number, node.knownPeers}
	writeRes(w, res)
}

// listBalanceHandler responds with the
// latest block hash and the current balances
func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}

// txnAddHandler adds the given valid transaction to the current state
func txnAddHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	req := TxnAddReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	txn := database.Txn{
		From:  database.NewAccount(req.From),
		To:    database.NewAccount(req.To),
		Value: req.Value,
		Data:  req.Data}

	block := database.NewBlock(
		state.LatestBlockHash(),
		state.NextBlockNumber(),
		uint64(time.Now().Unix()),
		[]database.Txn{txn},
	)
	hash, err := state.AddBlock(block)
	if err != nil {
		writeErrRes(w, err)
	}
	writeRes(w, TxnAddRes{hash})
}

// syncHandler fetches newer block if present
func syncHandler(w http.ResponseWriter, r *http.Request, dataDir string) {
	//get target node's latest block hash
	reqHash := r.URL.Query().Get(endpointSyncQueryKeyFromBlock)

	hash := database.Hash{}
	err := hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		writeErrRes(w, err)
		return
	}

	blocks, err := database.GetBlocksAfter(hash, dataDir)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, SyncRes{Blocks: blocks})
}
