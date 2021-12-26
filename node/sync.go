package node

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// sync queries known peers for new blocks and peers
// every given interval.
func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Searching for new Peers and Blocks")
			n.fetchNewBlocksAndPeers()
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

// fetchNewBlocksAndPeers() gets newer blocks and peers info
// from current known peer list
func (n *Node) fetchNewBlocksAndPeers() {
	for _, knownPeer := range n.knownPeers {
		status, err := queryPeerStatus(knownPeer)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}

		// check if new blocks are there
		localBlockNumber := n.state.LatestBlock().Header.Number
		if localBlockNumber < status.Number {
			fmt.Printf("Found %d new blocks from Peer %s\n", status.Number-localBlockNumber, knownPeer.IP)
			//prolly need to fetch new data
		}

		// check for new peer nodes, adds to current peer list if found
		for _, maybeNewPeer := range status.KnownPeers {
			newPeer, isKnownPeer := n.knownPeers[maybeNewPeer.TcpAddress()]
			if !isKnownPeer {
				fmt.Printf("Found new Peer %s\n", knownPeer.TcpAddress())
				n.knownPeers[maybeNewPeer.TcpAddress()] = newPeer
			}

		}
	}
}

// queryPeerStatus queries enpointStatus of the given PeerNode
func queryPeerStatus(peer PeerNode) (StatusRes, error) {
	url := fmt.Sprintf("http://%s/%s", peer.TcpAddress(), endpointStatus)
	res, err := http.Get(url)
	if err != nil {
		return StatusRes{}, err
	}

	statusRes := StatusRes{}
	err = readRes(res, &statusRes)
	if err != nil {
		return statusRes, err
	}

	return statusRes, nil
}
