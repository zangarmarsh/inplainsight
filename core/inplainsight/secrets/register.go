package secrets

type SecretFactoryInitializer func(serialized string) SecretInterface

var SecretsModelRegister = make(map[MagicNumber]SecretFactoryInitializer)
