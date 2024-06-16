package config

type Config interface {
	GetJWTSecretKey() []byte
}
