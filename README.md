# inplainsight

<img src="https://zangarmarsh.semaphoreci.com/badges/inplainsight/branches/main.svg">

This is a platform-independent **password (secret?)** manager which hides your secrets _in plain sight_. It takes extreme care about
the reliability and safety of your data.

## How does it work
Given an arbitrary master password _inplainsight_ derives two 32 bytes length keys through a slow hashing algorithm. They will encrypt the header and the data through AES-256 CTR. An HMAC is appended to the ciphertexts, to ensure the integrity of the encrypted secrets while decrypting.

The ciphertext is then interwoven within the pixels of an image through a process of adaptive steganography.


The media file might be stored locally, but we advise to keep a couple of online backups, just to ensure a good level of data redundancy.

## How does it work from a user perspective
You only need to remember a master password to get access to your secrets.

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
- Implement `haveibeenpwned` optional periodical checks
- Browser extension