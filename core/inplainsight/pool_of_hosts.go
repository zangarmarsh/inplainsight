package inplainsight

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"golang.org/x/exp/rand"
)

type PoolOfHosts struct {
	pool []*steganography.SecretHostInterface
}

func (p *PoolOfHosts) Add(host *steganography.SecretHostInterface) {
	p.pool = append(p.pool, host)
}

func (p *PoolOfHosts) Random(requiredSpace int) *steganography.SecretHostInterface {
	var eligibles []*steganography.SecretHostInterface

	for _, host := range p.pool {
		if int((*host).Cap()-(*host).Len()) >= requiredSpace {
			eligibles = append(eligibles, host)
		}
	}

	if len(eligibles) == 0 {
		return nil
	}

	return eligibles[rand.Intn(len(eligibles))]
}
