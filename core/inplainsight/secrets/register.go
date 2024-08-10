package secrets

type SecretFactoryInitializer func(serialized string) SecretInterface

var RegisteredSecrets = make(map[MagicNumber]SecretFactoryInitializer)
