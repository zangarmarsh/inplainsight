package secrets

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"strings"
)

const separator uint8 = '\x03'

// Container is a convenient wrapper which contains several SimpleSecret objects within the same host
type Container struct {
	secrets []*SimpleSecret
	Host    steganography.HostInterface
}

func (c *Container) Serialize() (serialized string) {
	for _, secret := range c.secrets {
		serialized += secret.Serialize() + string(separator)
	}

	return
}

func (c *Container) Unserialize(content string) {
	for _, singleSecretContent := range strings.Split(content, string(separator)) {
		if len(singleSecretContent) > 0 {
			if secret := (&SimpleSecret{}).UnserializeSecret(singleSecretContent); secret != nil {
				c.secrets = append(c.secrets, secret)
			}
		}
	}
}

func (c *Container) Add(secret *SimpleSecret) {
	c.secrets = append(c.secrets, secret)
}

func (c *Container) GetItems() []*SimpleSecret {
	return c.secrets
}
