package inplainsight

import (
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"golang.org/x/exp/rand"
)

type HostsPool struct {
	pool []*secrets.Container
}

func NewHostsPool() *HostsPool {
	return &HostsPool{
		pool: make([]*secrets.Container, 0),
	}
}

func (p *HostsPool) Add(host *secrets.Container) {
	p.pool = append(p.pool, host)
}

func (p *HostsPool) Random(requiredSpace int) *secrets.Container {
	var eligibles []*secrets.Container

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
