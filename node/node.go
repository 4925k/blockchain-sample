package node

import (
	"blockchain-sample/database"
	"context"
	"fmt"
	"net/http"
)

const (
	DefaultHttpPort               = 8080
	endpointStatus                = "/node/status"
	endpointSync                  = "/node/sync"
	endpointSyncQueryKeyFromBlock = "fromBlock"
)

type Node struct {
	dataDir    string
	port       uint64
	state      *database.State
	knownPeers map[string]PeerNode
}

// BalanceRes stores the block hash and balances
type BalancesRes struct {
	Hash    database.Hash             `json:"block_hash"`
	Balance map[database.Account]uint `json:"balances"`
}

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootStrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

// NewPeerNode returns a new peer node
func NewPeerNode(ip string, port uint64, isbootstrap, isactive bool) *PeerNode {
	return &PeerNode{ip, port, isbootstrap, isactive}
}

// New returns a new node
func New(dataDir string, port uint64, bootstrap PeerNode) *Node {
	knownPeers := make(map[string]PeerNode)
	knownPeers[bootstrap.TcpAddress()] = bootstrap
	return &Node{
		dataDir:    dataDir,
		port:       port,
		knownPeers: knownPeers,
	}
}

// Run starts the HTTP server and APIs
func (n *Node) Run() error {
	ctx := context.Background()
	fmt.Printf("Listening on HTTP Port: %d\n", n.port)
	state, err := database.NewStateFromDisk(n.dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.state = state
	go n.sync(ctx)

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})
	http.HandleFunc("/txn/add", func(w http.ResponseWriter, r *http.Request) {
		txnAddHandler(w, r, state)
	})
	http.HandleFunc(endpointStatus, func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, n)
	})
	http.HandleFunc(endpointSync, func(w http.ResponseWriter, r *http.Request) {
		syncHandler(w, r, n.dataDir)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", n.port), nil)

}

// TcpAddress returns "a.b.c.d:port" format ip address
func (pn PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}
