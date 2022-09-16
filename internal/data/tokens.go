package data

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"hash"
	"time"
)


const (
	ScopeActivation = "activation"
)


type Token struct {
	Plaintext string
	Hash[] byte
	UserID int64
	Expiry time.Time
	Scope string
}


func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error){

	token := &Token{
		UserID:    userID,
		Expiry:    time.Time{},
		Scope:     ScopeActivation,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
