package secrets

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"strings"
)

const separator uint8 = '\x03'

// Container is a convenient wrapper which contains several SimpleSecret objects within the same host
type Container struct {
	secrets []SecretInterface
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
			secretHeader := NewHeader(singleSecretContent[0])
			singleSecretContent = singleSecretContent[1:]

			if SecretsModelRegister[secretHeader.mn] != nil {
				if secret := SecretsModelRegister[secretHeader.mn](singleSecretContent); secret != nil {
					secret.SetHeader(secretHeader)
					c.secrets = append(c.secrets, secret)
				}
			}
		}
	}
}

func (c *Container) Add(secret SecretInterface) {
	c.secrets = append(c.secrets, secret)
}

func (c *Container) GetItems() []SecretInterface {
	return c.secrets
}
