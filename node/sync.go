package node

import (
	"blockchain-sample/database"
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
			n.doSync()
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *Node) doSync() {
	for _, peer := range n.knownPeers {
		if n.ip == peer.IP && n.port == peer.Port {
			return
		}

		fmt.Printf("Searching for new peers and their blocks: %s", peer.TcpAddress())
		status, err := queryPeerStatus(peer)
		if err != nil {
			fmt.Println("ERROR =>", err)
			fmt.Printf("Peer %s was removed from KnownPeers", peer.TcpAddress())
			n.RemovePeer(peer)
			continue
		}
		err = n.joinKnownPeers(peer)
		if err != nil {
			fmt.Println("ERROR =>", err)
			continue
		}

		err = n.syncBlocks(peer, status)
		if err != nil {
			fmt.Println("ERROR =>", err)
			continue
		}

		n.syncKnownPeers(peer, status)
	}
}

func (n *Node) syncBlocks(peer PeerNode, status StatusRes) error {
	localBlockNumber := n.state.LatestBlock().Header.Number

	// check if peer has no blocks
	if status.Hash.IsEmpty() {
		return nil
	}

	// check if peer has lesser blocks than us
	if status.Number < localBlockNumber {
		return nil
	}

	// check if it is genesis block and we already synced it
	if status.Number == 0 && !n.state.LatestBlockHash().IsEmpty() {
		return nil
	}

	newBlockCount := status.Number - localBlockNumber

	// display found 1 new block if we sync genesis block
	if localBlockNumber == 0 && status.Number == 0 {
		newBlockCount = 1
	}

	fmt.Printf("Found %d new blockd from %s\n", newBlockCount, peer.TcpAddress())

	blocks, err := fetchBlocksFromPeer(peer, n.state.LatestBlockHash())
	if err != nil {
		return err
	}

	return n.state.AddBlocks(blocks)
}

func (n *Node) syncKnownPeers(peer PeerNode, status StatusRes) {
	for _, statusPeer := range status.KnownPeers {
		if !n.IsKnownPeer(statusPeer) {
			fmt.Printf("Found new peer %s/n", statusPeer.TcpAddress())
			n.AddPeer(statusPeer)
		}
	}
}

func (n *Node) joinKnownPeers(peer PeerNode) error {
	if peer.connected {
		return nil
	}

	url := fmt.Sprintf("http://%s%s?%s=%s&%s=%d", peer.TcpAddress(), endpointAddPeer, endpointAddPeerQueryKeyIP, n.ip, endpointAddPeerQueryKeyPort, n.port)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	addPeerRes := AddPeerRes{}
	err = readRes(res, &addPeerRes)
	if err != nil {
		return fmt.Errorf(addPeerRes.Error)
	}

	knownPeer := n.knownPeers[peer.TcpAddress()]
	knownPeer.connected = addPeerRes.Success

	n.AddPeer(knownPeer)

	if !addPeerRes.Success {
		return fmt.Errorf("unable to join known peer %s", peer.TcpAddress())
	}

	return nil
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

func fetchBlocksFromPeer(peer PeerNode, fromBlock database.Hash) ([]database.Block, error) {
	fmt.Println("Importing Blocks from Peer %s\n", peer.TcpAddress())

	url := fmt.Sprintf("http://%s%s?%s=%s", peer.TcpAddress(), endpointSync, endpointSyncQueryKeyFromBlock, fromBlock)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	syncRes := SyncRes{}
	err = readRes(res, &syncRes)
	if err != nil {
		return nil, err
	}

	return syncRes.Blocks, nil
}
func (n *Node) AddPeer(peer PeerNode) {
	n.knownPeers[peer.TcpAddress()] = peer
}

func (n *Node) RemovePeer(peer PeerNode) {
	delete(n.knownPeers, peer.TcpAddress())
}

func (n *Node) IsKnownPeer(peer PeerNode) bool {
	if peer.IP == n.ip && peer.Port == n.port {
		return true
	}

	_, isKnownPeer := n.knownPeers[peer.TcpAddress()]

	return isKnownPeer
}
