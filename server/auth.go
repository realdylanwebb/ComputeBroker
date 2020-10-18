package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/square/go-jose"

	"github.com/square/go-jose/jwt"
)

//KeyChain contains keys for encrypting and signing tokens
type KeyChain struct {
	SigningKey    *[]byte
	EncryptionKey *[]byte
	TTL           time.Duration
}

//PrivateClaims represent an id to be passed to the client
type PrivateClaims struct {
	ID string `json:"id"`
	jwt.Claims
}

func openKey(path string) (*[]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fSize := stat.Size()
	fBuf := make([]byte, int(fSize))

	bytesRead, err := file.Read(fBuf)
	if err != nil {
		return nil, err
	}
	if bytesRead == 0 {
		return nil, errors.New("0 bytes read from key")
	}
	return &fBuf, nil
}

//Init initializes the keychain with signing and encryption keys and a TTL for tokens
func (keyChain *KeyChain) Init(encKeyPath string, sigKeyPath string, TTL time.Duration) error {
	encKey, err := openKey(encKeyPath)
	if err != nil {
		return err
	}
	sigKey, err := openKey(sigKeyPath)
	if err != nil {
		return err
	}

	keyChain.SigningKey = sigKey
	keyChain.EncryptionKey = encKey
	keyChain.TTL = TTL

	log.Printf("LENGTH OF SIG KEY: %d", len(*keyChain.SigningKey))
	return nil
}

//Validate retrieves claims from a token
func (keyChain *KeyChain) Validate(tokenString string) (*PrivateClaims, error) {
	token, err := jwt.ParseSignedAndEncrypted(tokenString)
	if err != nil {
		return nil, err
	}

	decrypted, err := token.Decrypt(*keyChain.EncryptionKey)
	if err != nil {
		return nil, err
	}

	out := PrivateClaims{}
	err = decrypted.Claims(*keyChain.SigningKey, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

//Sign creates a token from a set of claims
func (keyChain *KeyChain) Sign(claims *PrivateClaims) (string, error) {

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: *keyChain.SigningKey}, nil)
	if err != nil {
		return "", err
	}

	enc, err := jose.NewEncrypter(
		jose.A128GCM,
		jose.Recipient{
			Algorithm: jose.A128GCMKW,
			Key:       *keyChain.EncryptionKey,
		},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)

	if err != nil {
		return "", err
	}

	pubClaims := jwt.Claims{
		Expiry:   jwt.NewNumericDate(time.Now().Add(keyChain.TTL)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	raw, err := jwt.SignedAndEncrypted(signer, enc).Claims(pubClaims).Claims(*claims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return raw, err
}
