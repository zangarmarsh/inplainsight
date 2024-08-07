package inplainsight

import (
	"golang.org/x/exp/rand"
)

type HostsPool struct {
	pool []*SecretsContainer
}

func NewHostsPool() *HostsPool {
	return &HostsPool{
		pool: make([]*SecretsContainer, 0),
	}
}

func (p *HostsPool) Add(host *SecretsContainer) {
	p.pool = append(p.pool, host)
}

func (p *HostsPool) Random(requiredSpace int) *SecretsContainer {
	var eligibles []*SecretsContainer

	for _, container := range p.pool {
		if int((*container).Host.Cap()-(*container).Host.Len()) >= requiredSpace {
			eligibles = append(eligibles, container)
		}
	}

	if len(eligibles) == 0 {
		return nil
	}

	return eligibles[rand.Intn(len(eligibles))]
}

func (p *HostsPool) Reset() {
	p.pool = nil
}
