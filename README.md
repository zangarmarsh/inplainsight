# inplainsight

<img src="https://zangarmarsh.semaphoreci.com/badges/inplainsight/branches/main.svg">

This is a platform-independent **password (secret?)** manager which hides your secrets _in plain sight_. It takes extreme care about
the reliability and safety of your data.

## How does it work
### Encryption
Given an arbitrary master password _inplainsight_ derives two 32 bytes length keys through a slow hashing algorithm. They will encrypt the header and the data through AES-256 CTR. An HMAC is appended to the ciphertexts, to ensure the integrity of the encrypted secrets while decrypting.

### Steganography
The ciphertext is then interwoven within any supported media file through a process of adaptive steganography.

### Storage
Media file(s) might be stored locally or remotely, depending on the source of data used while logging in.
It would be advisable to keep a couple of remote backups, just to ensure a good level of data redundancy.

## Secrets
### How to implement your own secret structure
Secret models are defined in `core/inplainsight/secrets/`. To create your own model you just need to create a struct
which extends `secrets.AbstractSecret` and implements `secrets.SecretInterface`.

Look into [SimpleSecret](core/inplainsight/secrets/simple/)
or [WebsiteCredential](core/inplainsight/secrets/website/) if you need an example of implementation.

When defining a new model you'll also need to specify a custom `secrets.MagicNumber`.
Please keep in mind that value 0x00 and 0x03 are reserved.

In the end, [register](core/inplainsight/secrets/website/website_credentials.go#L20) it in an init function to make it globally available:
```golang
func init() {
	secrets.SecretsModelRegister[magicNumber] = func(serialized string)secrets.SecretInterface {
	// do something here if you need to...
        return (&YourSecretModel{}).Unserialize(serialized)
	}
}
```

### Supported secret structures
| ID   | Type     |Fields|
|------|----------|------|
 | 0x01 | `Secret` |Title, Description, Secret|
 | 0x02 | `WebsiteCredential` | URL, Note, Account, Password| 

## Media formats
### How to implement a new media format
Media formats live under the folder `core/steganography/medium/` and each one must have its own folder and dedicated tests.

Structs implementing `steganography.HostInterface` and registered can be used as media format.
This is how you register a new media format:

```go
// core/steganography/medium/yourmediaformat/register.go

package image

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
)

func init() {

  // Extend `Media` collection with a callback that returns an  instance
  // of `steganography.HostInterface` if the given `filePath` can be
  // handled by this `Media`. The check is typically based on specific conditions,
  // such as mimetype or content extension. Otherwise, return `nil`.
	
  steganography.Media = append(
    steganography.Media,
    func(filePath string) steganography.HostInterface {
      // ...
      
      return nil
    },
  )
}
```

### Supported media formats
- `images/*` - will eventually output `image/png` binary data 
- [ ... ]

## Supported sources & protocols
- file://
- [ ... ]

## Roadmap
- ~~Complete refactoring~~
- ~~Output image formats other than actual `png`~~ ( ngl, that was faked atm - needs lots of effort )
- Improve `Secret`
  - ~~Multiple secrets in one medium~~
  - ~~Make `Secret` more abstract and implementable in order to be easily extended~~
  - Give the user the ability to choose which file will be used (default will be `random`)
  - Exclusive host for one secret
  - `stealth mode` file header encryption
  - ~~Add secret icons~~
  - ~~Add custom `Action`~~
- Feature that allows to move a specific `Secret` into another choosable medium
- New `Secret` models:
  - ~~Website~~
  - Note
  - File
- Optional user preferences persistance
  - ~~Pool path at login~~
  - ~~Logout on screen lock~~
  - ~~Session timeout while inactive~~
  - `haveibeenpwned` optional periodical checks
- Blank image generation
- Support new data sources
  - `HTTPS`
  - `S3`
  - `FTP`
  - `SSH`
- Support `single-file` initialization mode
- Dockerization
- Self-hostable version
  - optional `2FA` through TOTP
- Support `hardware keys`
- Steganography the following media formats:
    - Audio files `MP3/WAV` (?)
    - `MP4` (?)
- Browser extension
- Pool of data-sources from text file (?)