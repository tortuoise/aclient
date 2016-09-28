package aclient

import (
	"crypto/rand"
	rsa "crypto/rsa"
	"crypto/sha256"
	"hash"
	"io"
)

type RSA struct {
	Hash    hash.Hash
	Entropy io.Reader
	Pubkey  *rsa.PublicKey
	Prikey  *rsa.PrivateKey
	Cipher  []byte
	Plain   []byte
	Label   []byte
}

func NewRSA() (*RSA, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &RSA{
		Hash:    sha256.New(),
		Entropy: rand.Reader,
		Prikey:  key,
		Pubkey:  &key.PublicKey,
	}, nil
}

func (this *RSA) Encrypt() ([]byte, error) {
	out, err := rsa.EncryptOAEP(this.Hash, this.Entropy, this.Pubkey, this.Plain, this.Label)
	return out, err
}

func (this *RSA) Decrypt() ([]byte, error) {
	msg, err := rsa.DecryptOAEP(this.Hash, this.Entropy, this.Prikey, this.Cipher, this.Label)
	return msg, err
}

func (this *RSA) Randomize() {
	this.Entropy = rand.Reader
}
