package inplainsight

import (
	"errors"
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/cryptography"
	"github.com/zangarmarsh/inplainsight/core/events"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"log"
	"sort"
	"strings"
	"time"

	_ "github.com/zangarmarsh/inplainsight/core/steganography/medium/image"
)

const (
	Version = "1.0.0"
)

var InPlainSight = &InPlainSightClient{
	Hosts: *NewHostsPool(),
}

func Conceal(secret secrets.SecretInterface) error {
	isCreating := true

	// ToDo maybe it is worth creating a isEmpty() method on SimpleSecret
	secretMessage := []byte(secret.Serialize())
	secretMessage = append(secretMessage, secrets.SecretSeparator)

	if len(secretMessage) == 0 {
		return errors.New("provided secret is empty")
	}

	var contentEncryptionKey []byte
	var err error

	if secret.GetID() != "" {
		for _, s := range InPlainSight.Secrets {
			if s.GetID() == secret.GetID() {
				isCreating = false
				break
			}
		}
	} else {
		secret.AssignRandomID()
	}

	log.Println("Is creating?", isCreating, "container", secret.GetContainer())

	if secret.GetContainer() == nil {
		// At this point there should be already a bunch of secret hosts
		container := InPlainSight.Hosts.Random(len(secretMessage))

		if container != nil {
			secret.SetContainer(container)
			secret.GetContainer().Add(secret)
		} else {
			return errors.New("there's no eligible medium within the `Pool path` which can contain this secret")
		}
	}

	if secret.GetContainer() != nil {
		if len(InPlainSight.MasterPassword) != 0 {
			contentEncryptionKey, _, err = cryptography.DeriveEncryptionKeysFromPassword([]byte(InPlainSight.MasterPassword))
			if err != nil {
				return err
			}

			log.Printf("SimpleSecret: %+v (%+v)", secret, &secret)
			log.Printf("Updating container file with secrets: %+v", secret.GetContainer().GetItems())
			secretMessage = []byte(secret.GetContainer().Serialize())

			secretMessage, err = cryptography.Encrypt(secretMessage, contentEncryptionKey)
			if err != nil {
				return err
			}
		}

		log.Printf("Media at %v (%d/%d) has been chosen to host secret %+v", secret.GetContainer().Host.GetPath(), secret.GetContainer().Host.Len(), secret.GetContainer().Host.Cap(), secret)

		err = secret.GetContainer().Host.Interweave(string(secretMessage))
		if err != nil {
			return err
		}

		if isCreating {
			InPlainSight.Secrets = append(InPlainSight.Secrets, secret)

			InPlainSight.Trigger(events.Event{
				CreatedAt: time.Now(),
				EventType: events.SecretAdded,
				Data: map[string]interface{}{
					"secret": secret,
				},
			})
		} else {
			InPlainSight.Trigger(events.Event{
				CreatedAt: time.Now(),
				EventType: events.SecretUpdated,
				Data: map[string]interface{}{
					"secret": secret,
				},
			})
		}

		// Clean up removable secret once they have been physically removed from the medium
		for i, item := range InPlainSight.Secrets {
			if item.IsDeleatable() {
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
		container := secrets.Container{Host: host}
		InPlainSight.Hosts.Add(&container)

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

			for _, secret := range container.GetItems() {
				secret.SetContainer(&container)
				InPlainSight.Secrets = append(InPlainSight.Secrets, secret)

				sort.Slice(InPlainSight.Secrets, func(i, j int) bool {
					return InPlainSight.Secrets[i].GetTitle() < InPlainSight.Secrets[j].GetTitle()
				})

				InPlainSight.Trigger(events.Event{
					CreatedAt: time.Now(),
					EventType: events.SecretDiscovered,
					Data: map[string]interface{}{
						"secret": secret,
					},
				})
			}
		}
	}

	return nil
}
