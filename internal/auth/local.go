package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

// If there's and existing local key, we read it from the file and use it
func CheckLocalKeyExists() (*PasetoLocalKey, error) {
	readKey, err := os.ReadFile("local.key")
	if err != nil {
		return nil, fmt.Errorf("Failed to read symmetricKey from system, does local.key exist? %w", err)
	}

	localKey, err := paseto.V4SymmetricKeyFromBytes(readKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate local key from local.key: %w", err)
	}
	return &PasetoLocalKey{
		localKey,
	}, nil
}

// The function to generate the key
// Always gets the same key as long as we don't delete the .key file
func MakePasetoLocalKey() (*PasetoLocalKey, error) {
	// If the key exists, use it otherwise generate new (all tokens with previous key will be invalidated in that case)
	secretKey, err := CheckLocalKeyExists()
	if err == nil {
		log.Println("Using existing key for local paseto :)")
		return secretKey, nil
	}

	// Writing in a file to persist the key
	// Read-Write access of file ONLY TO OWNER as it's the secret local key
	localKey := paseto.NewV4SymmetricKey()
	err = os.WriteFile("local.key", localKey.ExportBytes(), 0600)
	if err != nil {
		return nil, fmt.Errorf("Failed to write localKey to file: %w", err)
	}

	log.Println("Generated new key for local paseto :)")
	return &PasetoLocalKey{
		localKey,
	}, nil
}

// Create the token -> add the payload to token -> encrypt it with the key
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
