package inplainsight

import (
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"golang.org/x/exp/rand"
	"strings"
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

func (p *HostsPool) SearchByContainerPath(query string) []*secrets.Container {
	var results []*secrets.Container

	for _, pool := range p.pool {
		if strings.Contains(pool.Host.GetPath(), query) {
			results = append(results, pool)
		}
	}

	return results
}

func (p *HostsPool) List() []*secrets.Container {
	return p.pool
}

func (p *HostsPool) Reset() {
	p.pool = nil
}
