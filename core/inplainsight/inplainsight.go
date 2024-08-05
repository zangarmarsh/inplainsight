package inplainsight

import (
	"errors"
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/cryptography"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"github.com/zangarmarsh/inplainsight/ui/events"
	"log"
	"strings"
	"time"

	_ "github.com/zangarmarsh/inplainsight/core/steganography/medium/image"
)

const (
	Version = "1.0.0"
)

const secretSeparator = "\x02"

var InPlainSight = &InPlainSightClient{
	Hosts: *NewHostsPool(),
}

func Conceal(secret *Secret) error {
	isCreating := false

	// ToDo maybe worth creating a isEmpty() method on Secret
	secretMessage := []byte(secret.Serialize())
	secretMessage = append(secretMessage, separator)

	if len(secretMessage) == 0 {
		return errors.New("provided secret is empty")
	}

	var contentEncryptionKey []byte
	var err error

	// If `secret.Container` is null it will likely mean that we're creating it
	if secret.Container == nil {
		isCreating = true

		// At this point there should be already a bunch of secret hosts
		secret.Container = InPlainSight.Hosts.Random(len(secretMessage))

		if secret.Container != nil {
			secret.Container.Add(secret)
		}
	}

	if secret.Container != nil {
		if len(InPlainSight.MasterPassword) != 0 {
			contentEncryptionKey, _, err = cryptography.DeriveEncryptionKeysFromPassword([]byte(InPlainSight.MasterPassword))
			if err != nil {
				return err
			}

			log.Printf("Secret: %+v (%+v)", secret, &secret)
			log.Printf("Updating container file with secrets: %+v", secret.Container.secrets)
			secretMessage = []byte(secret.Container.Serialize())

			secretMessage, err = cryptography.Encrypt(secretMessage, contentEncryptionKey)
			if err != nil {
				return err
			}
		}

		log.Printf("Media at %v (%d/%d) has been chosen to host secret %+v", secret.Container.Host.GetPath(), secret.Container.Host.Len(), secret.Container.Host.Cap(), secret)

		err = secret.Container.Host.Interweave(string(secretMessage))
		if err != nil {
			return err
		}

		if isCreating {
			InPlainSight.Secrets = append(InPlainSight.Secrets, secret)

			InPlainSight.Trigger(events.Event{
				CreatedAt: time.Now(),
				EventType: events.AddedNewSecret,
				Data: map[string]interface{}{
					"secret": secret,
				},
			})
		} else {
			InPlainSight.Trigger(events.Event{
				CreatedAt: time.Now(),
				EventType: events.UpdatedSecret,
				Data: map[string]interface{}{
					"secret": secret,
				},
			})
		}

		// Clean up removable secret once they have been physically removed from the medium
		for i, item := range InPlainSight.Secrets {
			if item.deleatable {
				log.Println("Secret", item, "is deleatable")
				InPlainSight.Secrets = append(InPlainSight.Secrets[:i], InPlainSight.Secrets[i+1:]...)

				break
			}
		}
	} else {
		return errors.New("unable to interweave secret")
	}

	return nil
}

func Reveal(fileName string) error {
	var decrypted []byte
	path := fmt.Sprintf("%s/%s", strings.TrimRight(InPlainSight.Path, "/\\"), fileName)
	host := steganography.New(path)

	if host != nil {
		container := SecretsContainer{Host: host}

		if len(*host.Data()) > 0 {
			var contentEncryptionKey []byte
			var err error

			contentEncryptionKey, _, err = cryptography.DeriveEncryptionKeysFromPassword(
				[]byte(InPlainSight.MasterPassword),
			)
			if err != nil {
				return err
			}

			decrypted, err = cryptography.Decrypt([]byte(*host.Data()), contentEncryptionKey)
			if err != nil {
				return err
			}

			container.Unserialize(string(decrypted))

			for _, secret := range container.secrets {
				secret.Container = &container
				InPlainSight.Secrets = append(InPlainSight.Secrets, secret)

				InPlainSight.Trigger(events.Event{
					CreatedAt: time.Now(),
					EventType: events.DiscoveredNewSecret,
					Data: map[string]interface{}{
						"secret": secret,
					},
				})
			}
		}

		InPlainSight.Hosts.Add(&container)
	}

	return nil
}
