package peer

import (
	"context"
	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/connection_pool"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
	"math/rand"
)

type ClientPeer struct {
	Peer
}

func NewClientPeer(tunnelNum int, endpoint string, cipher tunnel.Cipher) ClientPeer {
	peerID := rand.Uint32()
	return NewClientPeerWithID(peerID, tunnelNum, endpoint, cipher)
}

func NewClientPeerWithID(peerID uint32, tunnelNum int, endpoint string, cipher tunnel.Cipher) ClientPeer {
	peerCtx, removePeerFunc := context.WithCancel(context.Background())

	poolManager := tunnel_pool.NewClientManager(tunnelNum, endpoint, peerID, cipher)
	tunnelPool := tunnel_pool.NewTunnelPool(peerID, &poolManager, peerCtx)
	connectionPool := connection_pool.NewConnectionPool(&tunnelPool, peerCtx)

	return ClientPeer{
		Peer: Peer{
			peerID:         peerID,
			connectionPool: connectionPool,
			tunnelPool:     tunnelPool,
			ctx:            peerCtx,
			cancel:         removePeerFunc,
		},
	}
}

func (cp *ClientPeer) Dial(address string) connection.Connection {
	conn := cp.connectionPool.NewPooledInboundConnection()
	conn.SendConnect(address)
	return conn
}
