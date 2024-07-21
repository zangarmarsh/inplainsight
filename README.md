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
- Complete refactoring
- Output image formats other than actual `png`
- `stealth mode` file header encryption
- Dockerization
- self-hostable version
  - optional `2FA`
    - Evaluate if it makes sense using it even locally taking care of using strict file permissions
- Support new data sources
  - HTTPS
  - S3
  - FTP
  - SSH
- Pool of data-sources wrapped in a file
- Steganography on the following media formats:
    - Audio files (MP3/WAV)
    - MP4
    - ?
