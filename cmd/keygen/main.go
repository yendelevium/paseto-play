package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
)

// Run this script ONCE to make
func main() {
	secret := os.Getenv("SECRET_PASTEO")
	hash := sha256.Sum256(([]byte(secret)))
	standardKey := ed25519.NewKeyFromSeed(hash[:])

	log.Printf("COPY TO .ENV\n")
	log.Printf("64-byte hashed-secret: %v \n", hex.EncodeToString(standardKey[:]))
}
