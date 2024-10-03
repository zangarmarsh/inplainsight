package secrets

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"hash/crc32"
	"sort"
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
					secret.AssignRandomID()

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

func (c *Container) sort() {
	sort.Slice(c.secrets, func(i, j int) bool {
		return c.secrets[i].GetTitle() > c.secrets[j].GetTitle()
	})
}

func (c *Container) checksum() uint32 {
	c.sort()

	var buffer []byte

	for _, secret := range c.secrets {
		buffer = append(buffer, secret.Serialize()...)
	}

	return crc32.Checksum(buffer, crc32.MakeTable(crc32.IEEE))
}
