package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"os"
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

type PublicKey struct {
	key ed25519.PublicKey
}

func Generate() (*PrivateKey, *PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return &PrivateKey{key: priv}, &PublicKey{key: pub}, nil
}

func NewPrivateKey(raw []byte) (*PrivateKey, error) {
	if len(raw) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key size")
	}
	return &PrivateKey{key: ed25519.PrivateKey(raw)}, nil
}

func NewPublicKey(raw []byte) (*PublicKey, error) {
	if len(raw) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key size")
	}
	return &PublicKey{key: ed25519.PublicKey(raw)}, nil
}

func (k *PrivateKey) Bytes() []byte {
	return k.key
}

func (k *PrivateKey) Hex() string {
	return hex.EncodeToString(k.key)
}

func (k *PrivateKey) Public() *PublicKey {
	return &PublicKey{key: k.key.Public().(ed25519.PublicKey)}
}

func (k *PrivateKey) Sign(data []byte) []byte {
	return ed25519.Sign(k.key, data)
}

func (k *PrivateKey) SaveToFile(path string) error {
	der, err := x509.MarshalPKCS8PrivateKey(k.key)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}
	return os.WriteFile(path, pem.EncodeToMemory(block), 0600)
}

func LoadPrivateKey(pemData []byte) (*PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("invalid pem data")
	}

	if len(block.Bytes) == ed25519.PrivateKeySize {
		return NewPrivateKey(block.Bytes)
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := parsed.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("not an ed25519 private key")
	}
	return &PrivateKey{key: key}, nil
}

func LoadPrivateKeyFromFile(path string) (*PrivateKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadPrivateKey(raw)
}

func (k *PublicKey) Bytes() []byte {
	return k.key
}

func (k *PublicKey) Hex() string {
	return hex.EncodeToString(k.key)
}

func (k *PublicKey) Verify(data, signature []byte) bool {
	return ed25519.Verify(k.key, data, signature)
}

func (k *PublicKey) SaveToFile(path string) error {
	der, err := x509.MarshalPKIXPublicKey(k.key)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return os.WriteFile(path, pem.EncodeToMemory(block), 0644)
}

func LoadPublicKey(pemData []byte) (*PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("invalid pem data")
	}

	if len(block.Bytes) == ed25519.PublicKeySize {
		return NewPublicKey(block.Bytes)
	}

	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := parsed.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("not an ed25519 public key")
	}
	return &PublicKey{key: key}, nil
}

func LoadPublicKeyFromFile(path string) (*PublicKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadPublicKey(raw)
}
