package libp2p

import (
	"log/slog"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"golang.org/x/time/rate"
)

type blocker interface {
	isBlocked(peer.ID) bool
}

type gater struct {
	limiter *limiter
	blocker blocker
	logger  *slog.Logger
}

// NewGater returns a new libp2p connection gater.
func newGater(logger *slog.Logger) *gater {
	return &gater{
		limiter: newLimiter(),
		logger:  logger,
	}
}

func (g *gater) setBlocker(b blocker) {
	g.blocker = b
}

func (g *gater) InterceptPeerDial(p peer.ID) bool {
	if g.blocker != nil && g.blocker.isBlocked(p) {
		g.logger.Warn("blocked dial: peer is blocklisted", "peerID", p)
		return false
	}
	return true
}

// No filters atm for dialing multiaddrs.
func (*gater) InterceptAddrDial(_ peer.ID, _ multiaddr.Multiaddr) bool {
	return true
}

func (g *gater) InterceptAccept(n network.ConnMultiaddrs) bool {
	if !g.limiter.Allow(n.RemoteMultiaddr()) {
		g.logger.Warn(
			"blocked accept: rate limit exceeded",
			"remoteAddr", n.RemoteMultiaddr(),
		)
		return false
	}
	return true
}

func (g *gater) InterceptSecured(_ network.Direction, p peer.ID, _ network.ConnMultiaddrs) bool {
	if g.blocker != nil && g.blocker.isBlocked(p) {
		g.logger.Warn("blocked accept: peer is blocklisted", "peerID", p)
		return false
	}
	return true
}

func (*gater) InterceptUpgraded(_ network.Conn) (bool, control.DisconnectReason) {
	return true, 0
}

func newLimiter() *limiter {
	cache, _ := lru.New[string, *rate.Limiter](1000)
	return &limiter{
		ipLimiters: cache,
	}
}

type limiter struct {
	ipLimiters *lru.Cache[string, *rate.Limiter]
	mu         sync.Mutex
}

func (l *limiter) Allow(addr multiaddr.Multiaddr) bool {
	ip, err := manet.ToIP(addr)
	if err != nil {
		return false
	}

	l.mu.Lock()
	ipLimiter, found := l.ipLimiters.Get(ip.String())
	if !found {
		ipLimiter = rate.NewLimiter(rate.Every(1), 1)
		l.ipLimiters.Add(ip.String(), ipLimiter)
	}
	l.mu.Unlock()

	return ipLimiter.Allow()
}
