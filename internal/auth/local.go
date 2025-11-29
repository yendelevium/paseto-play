package auth

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

func MakePasetoLocalKey() *PasetoLocalKey {
	localKey := paseto.NewV4SymmetricKey()
	// Write this into a file later on to persist the local key
	// Then everytime this function is read, check if the file exists, read it and return the same key otherwise new key + file-write
	return &PasetoLocalKey{
		localKey: localKey,
	}
}

// Create the token -> add the payload to token -> sign it with privatekey
func (maker *PasetoLocalKey) CreateTokenEncrypted(payload any) (string, error) {
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

	// ENCRYPT the token -> as it's symmetric key, it needs to be secret and encrypted
	return token.V4Encrypt(maker.localKey, nil), nil
}

// Decrypt the token by using the same symmetric key and extracting the information
func (maker *PasetoLocalKey) DecryptToken(tokenString string) ([]byte, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	token, err := parser.ParseV4Local(maker.localKey, tokenString, nil)
	if err != nil {
		return nil, err
	}

	parsedJSON := token.ClaimsJSON()
	return parsedJSON, nil
}
