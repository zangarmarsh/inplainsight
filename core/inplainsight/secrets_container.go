package inplainsight

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"strings"
)

const separator uint8 = '\x03'

// SecretsContainer is a convenient wrapper which contains several Secret objects within the same host
type SecretsContainer struct {
	secrets []Secret
	Host    steganography.HostInterface
}

func (c *SecretsContainer) Serialize() (serialized string) {
	for _, secret := range c.secrets {
		serialized += secret.Serialize() + string(separator)
	}

	return
}

func (c *SecretsContainer) Unserialize(content string) {
	for _, singleSecretContent := range strings.Split(content, string(separator)) {
		if len(singleSecretContent) > 0 {
			if secret := UnserializeSecret(singleSecretContent); secret != nil {
				c.secrets = append(c.secrets, *secret)
			}
		}
	}
}

func (c *SecretsContainer) Add(secret *Secret) {
	c.secrets = append(c.secrets, *secret)
}
