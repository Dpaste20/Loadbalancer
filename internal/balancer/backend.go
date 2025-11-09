package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Backend Representation
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
	Connections  int
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

func (b *Backend) GetConnections() int {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Connections
}

func (b *Backend) IncrementConnections() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Connections++
}

func (b *Backend) DecrementConnections() {
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.Connections > 0 {
		b.Connections--
	}
}

func NewBackend(urlStr string) (*Backend, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:          u,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
		Connections:  0,
	}, nil
}
