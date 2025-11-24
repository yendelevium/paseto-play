package pasetotokens

import (
	"aidanwoods.dev/go-paseto"
)

// For Dependency Injection
type PasetoPublicKeyPair struct {
	privateKey paseto.V4AsymmetricSecretKey
	publicKey  paseto.V4AsymmetricPublicKey
}
