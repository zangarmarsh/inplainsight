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

### Supported secret structures
| ID   | Type     |Fields|
|------|----------|------|
 | 0x01 | `Secret` |Title, Description, Secret|

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
  - Make `Secret` more abstract and implementable in order to be easily extended
  - Support `single-file` mode
  - Give the user the ability to choose which file will be used (default will be `random`)
  - Exclusive host for one secret
  - `stealth mode` file header encryption
- Optional user preferences persistance
  - ~~Pool path at login~~
  - ~~Logout on screen lock~~
  - Session timeout while inactive
  - `haveibeenpwned` optional periodical checks
- Blank image generation
- Support new data sources
  - `HTTPS`
  - `S3`
  - `FTP`
  - `SSH`
- Dockerization
- self-hostable version
  - optional `2FA`
    - Evaluate if it makes sense using it even locally
- Pool of data-sources wrapped in a file
- Support `hardware keys`
- Steganography the following media formats:
    - Audio files `MP3/WAV`
    - `MP4`
- Browser extension