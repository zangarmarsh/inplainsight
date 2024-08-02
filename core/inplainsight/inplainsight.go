package inplainsight

import (
	"errors"
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/cryptography"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"github.com/zangarmarsh/inplainsight/ui/events"
	"strings"
	"time"

	_ "github.com/zangarmarsh/inplainsight/core/steganography/medium/image"
)

const (
	Version = "1.0.0"
)

const secretSeparator = "\x02"

var InPlainSight = &InPlainSightClient{
	Secrets: make(map[string]*Secret),
}

func Conceal(secret *Secret) error {
	secretMessage := []byte(secret.Serialize())

	// path := fmt.Sprintf("%s/%s", strings.TrimRight(InPlainSight.Path, "/\\"), fileName)
	path := secret.Host.GetPath()

	if len(secretMessage) == 0 {
		return errors.New("provided secret is empty")
	}

	var contentEncryptionKey []byte
	var err error

	if len(InPlainSight.MasterPassword) != 0 {
		contentEncryptionKey, _, err = cryptography.DeriveEncryptionKeysFromPassword([]byte(InPlainSight.MasterPassword))
		if err != nil {
			return err
		}

		secretMessage, err = cryptography.Encrypt(secretMessage, contentEncryptionKey)
		if err != nil {
			return err
		}
	}

	if secret.Host == nil {
		secret.Host = steganography.New(path)
	}

	if secret.Host != nil {
		err = secret.Host.Interweave(string(secretMessage))
		if err != nil {
			return err
		}

		// InPlainSight.Secrets[fileName] = secret

		InPlainSight.Trigger(events.Event{
			CreatedAt: time.Now(),
			EventType: events.AddedNewSecret,
			Data: map[string]interface{}{
				"secret": secret,
			},
		})

	} else {
		return errors.New("unable to interweave secret")
	}

	return nil
}

func Reveal(fileName string) error {
	var decrypted []byte
	path := fmt.Sprintf("%s/%s", strings.TrimRight(InPlainSight.Path, "/\\"), fileName)
	host := steganography.New(path)

	if host != nil && len(host.Data().Encrypted) > 0 {
		var contentEncryptionKey []byte
		var err error

		contentEncryptionKey, _, err = cryptography.DeriveEncryptionKeysFromPassword(
			[]byte(InPlainSight.MasterPassword),
		)

		if err != nil {
			return err
		}

		decrypted, err = cryptography.Decrypt([]byte(host.Data().Encrypted), contentEncryptionKey)

		if err != nil {
			return err
		}
	}

	secret := &Secret{}
	secret.Unserialize(string(decrypted))
	secret.Host = host

	InPlainSight.Secrets[fileName] = secret

	if decrypted != nil {
		InPlainSight.Trigger(events.Event{
			CreatedAt: time.Now(),
			EventType: events.DiscoveredNewSecret,
			Data: map[string]interface{}{
				"secret": *secret,
			},
		})
	}

	return nil
}
