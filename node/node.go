package node

import (
	"blockchain-sample/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var httpPort = 8080

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash    database.Hash             `json:"block_hash"`
	Balance map[database.Account]uint `json:"balances"`
}

type TxnAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxnAddRes struct {
	Hash database.Hash `json:"hash"`
}

func Run(path string) error {
	state, err := database.NewStateFromDisk(path)
	if err != nil {
		return err
	}
	defer state.Close()

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})
	http.HandleFunc("/txn/add", func(w http.ResponseWriter, r *http.Request) {
		txnAddHandler(w, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)

}

func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}

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

func writeErrRes(w http.ResponseWriter, err error) {
	jsonErrRes, _ := json.Marshal(ErrRes{err.Error()})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonErrRes)
}

func writeRes(w http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		writeErrRes(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(contentJson)
}

func readReq(r *http.Request, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read req body: %s", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request content: %s", err)
	}

	return nil
}
