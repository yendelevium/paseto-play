package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

// If there's an existing asymmetric key-pair, we read the key-pair from the files and use them
func CheckPublicKeyPairExists() (*PasetoPublicKeyPair, error) {
	readPrivateKey, err := os.ReadFile("public.rsa")
	if err != nil {
		return nil, fmt.Errorf("Failed to read privateKey from system, does public.rsa exist? %w", err)
	}
	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(readPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate private key from public.rsa: %w", err)
	}

	readPublicKey, err := os.ReadFile("public.pub.rsa")
	if err != nil {
		return nil, fmt.Errorf("Failed to read publicKey from system, does public.rsa exist? %w", err)
	}
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromBytes(readPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate public key from public.pub.rsa: %w", err)
	}

	return &PasetoPublicKeyPair{
		privateKey,
		publicKey,
	}, nil
}

// The function to generate the key pair
// Always generates the same key-pair as long as we don't delete the .rsa files
func MakePasetoKeyPair() (*PasetoPublicKeyPair, error) {
	// If the keypair exists, use those otherwise generate new (all tokens with previous keypair will be invalidated in that case)
	keyPair, err := CheckPublicKeyPairExists()
	if err == nil {
		log.Println("Using existing key-pair for public paseto :)")
		return keyPair, nil
	}

	privateKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := privateKey.Public()

	// Writing them in files to persist the keys
	// Read-Write access of file ONLY TO OWNER as it's the private key
	err = os.WriteFile("public.rsa", privateKey.ExportBytes(), 0600)
	if err != nil {
		return nil, fmt.Errorf("Failed to write privateKey to file: %w", err)
	}
	// Public key others can read but only owner can write
	err = os.WriteFile("public.pub.rsa", publicKey.ExportBytes(), 0644)
	if err != nil {
		return nil, fmt.Errorf("Failed to write publicKey to file: %w", err)
	}

	log.Println("Generated key-pair for public paseto :)")
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
