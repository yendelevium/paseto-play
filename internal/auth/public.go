package auth

import (
	"fmt"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

// The function to generate the key pair
// Always generates the same key-pair as we have the same secret
func MakePasetoKeyPair() (*PasetoPublicKeyPair, error) {
	secret := os.Getenv("SECRET_PASETO")
	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey := privateKey.Public()
	return &PasetoPublicKeyPair{
		privateKey,
		publicKey,
	}, nil
}

// Create the token -> add the payload to token -> sign it with privatekey
func (maker *PasetoPublicKeyPair) CreateToken(payload any) (string, error) {
	// Create the bare token with no payload
	token := paseto.NewToken()
	token.SetExpiration(time.Now().Add(5 * time.Minute))
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	// Set your payload under the "data" key
	err := token.Set("data", payload)
	if err != nil {
		return "", fmt.Errorf("failed to set payload: %w", err)
	}

	// Sign the token
	return token.V4Sign(maker.privateKey, nil), nil
}

// verify the token by using the public key and extracting the information
func (maker *PasetoPublicKeyPair) VerifyToken(tokenString string) ([]byte, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	token, err := parser.ParseV4Public(maker.publicKey, tokenString, nil)
	if err != nil {
		return nil, err
	}

	parsedJSON := token.ClaimsJSON()
	return parsedJSON, nil
}
