package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

const ReadOnly = 0o600

var (
	ErrIdentityEmpty   = errors.New("identity is empty")
	ErrBadPEMFile      = errors.New("bad PEM file")
	ErrUnknownPEMBlock = errors.New("unknown PEM block")
	ErrNoPrivateKey    = errors.New("no private key")
)

const RSAKeySize = 2048

// Identity is just a small struct that clearly differentiates between the
// private and public key of an RSA keypair.
type Identity struct {
	public  *rsa.PublicKey
	private *rsa.PrivateKey
}

func GenerateIdentity() *Identity {
	priv, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		panic(err)
	}

	return &Identity{
		private: priv,
		public:  &priv.PublicKey,
	}
}

func EmptyIdentity() *Identity {
	return &Identity{
		private: nil,
		public:  nil,
	}
}

func LoadIdentity(rsaFile string) (*Identity, error) {
	rsaText, err := os.ReadFile(rsaFile)
	if err != nil {
		return nil, err
	}

	id := EmptyIdentity()
	if err := id.UnmarshalText(rsaText); err != nil {
		return nil, err
	}

	return id, nil
}

func (r *Identity) SaveIdentity(rsaFile string) error {
	data, err := r.MarshalText()
	if err != nil {
		return err
	}

	return os.WriteFile(rsaFile, data, ReadOnly)
}

func (r *Identity) MarshalText() ([]byte, error) {
	if r.private != nil {
		return pem.EncodeToMemory(&pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(r.private),
		}), nil
	} else if r.public != nil {
		return pem.EncodeToMemory(&pem.Block{
			Type:    "RSA PUBLIC KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PublicKey(r.public),
		}), nil
	}

	return nil, ErrIdentityEmpty
}

func (r *Identity) UnmarshalText(text []byte) error {
	b, _ := pem.Decode(text)
	if b == nil {
		return ErrBadPEMFile
	}

	if b.Type == "RSA PRIVATE KEY" {
		p, e := x509.ParsePKCS1PrivateKey(b.Bytes)
		if e != nil {
			return e
		}

		r.private = p
		r.public = &p.PublicKey

		return nil
	} else if b.Type == "RSA PUBLIC KEY" {
		p, e := x509.ParsePKCS1PublicKey(b.Bytes)
		if e != nil {
			return e
		}

		r.public = p

		return nil
	}

	return ErrUnknownPEMBlock
}

func (r *Identity) String() string {
	b, e := r.MarshalText()
	if e != nil {
		return ""
	}

	return string(b)
}

// NewIdentity returns a new identity with spefied keys.
func NewIdentity(pri *rsa.PrivateKey) *Identity {
	return &Identity{
		private: pri,
		public:  &pri.PublicKey,
	}
}

// PubIdentity returns identity with public key. This identity object can
// only be used to verify messages.
func PubIdentity(pub *rsa.PublicKey) *Identity {
	return &Identity{
		private: nil,
		public:  pub,
	}
}

func (r *Identity) PublicKey() *rsa.PublicKey {
	return r.public
}

func (r *Identity) Public() *Identity {
	return &Identity{
		public:  r.public,
		private: nil,
	}
}

// Sign returns a signature made by combining the message and the signers private key
// With the r.Verify function, the signature can be checked.
func (r *Identity) Sign(msg []byte) ([]byte, error) {
	hs := r.getHashSum(msg)

	if r.private == nil {
		return nil, ErrNoPrivateKey
	}

	return rsa.SignPKCS1v15(rand.Reader, r.private, crypto.SHA256, hs)
}

// Verify checks if a message is signed by a given Public Key.
func (r *Identity) Verify(msg []byte, sig []byte, pubKey *rsa.PublicKey) error {
	hs := r.getHashSum(msg)

	if pubKey == nil {
		pubKey = r.PublicKey()
	}

	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hs, sig)
}

// Encrypt's the message using EncryptOAEP which encrypts the given message with RSA-OAEP.
// https://en.wikipedia.org/wiki/Optimal_asymmetric_encryption_padding
// Returns the encrypted message and an error.
func (r *Identity) Encrypt(msg []byte, key *rsa.PublicKey) ([]byte, error) {
	label := []byte("")
	hash := sha256.New()

	if key == nil {
		key = r.PublicKey()
	}

	return rsa.EncryptOAEP(hash, rand.Reader, key, msg, label)
}

// Decrypt a message using your private key.
// A received message should be encrypted using the receivers public key.
func (r *Identity) Decrypt(msg []byte) ([]byte, error) {
	label := []byte("")
	hash := sha256.New()

	return rsa.DecryptOAEP(hash, rand.Reader, r.private, msg, label)
}

func (r *Identity) getHashSum(msg []byte) []byte {
	h := sha256.New()
	h.Write(msg)

	return h.Sum(nil)
}
