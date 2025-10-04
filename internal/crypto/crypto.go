package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// CryptoManager handles cryptographic operations
type CryptoManager struct{}

// NewCryptoManager creates a new crypto manager instance
func NewCryptoManager() *CryptoManager {
	return &CryptoManager{}
}

// GenerateKeyPair generates an RSA public/private key pair
func (cm *CryptoManager) GenerateKeyPair() (privateKeyPEM, publicKeyPEM []byte, err error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Encode private key to PEM
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key to PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}

// SignData signs data with a private key
func (cm *CryptoManager) SignData(data []byte, privateKeyPEM []byte) (string, error) {
	// Decode private key
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", fmt.Errorf("failed to decode private key PEM")
	}

	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	privateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("not an RSA private key")
	}

	// Hash the data
	hash := sha256.Sum256(data)

	// Sign the hash
	signature, err := rsa.SignPSS(
		rand.Reader,
		privateKey,
		crypto.SHA256,
		hash[:],
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	// Return base64 encoded signature
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies a signature with a public key
func (cm *CryptoManager) VerifySignature(data []byte, signature string, publicKeyPEM []byte) (bool, error) {
	// Decode signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Decode public key
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return false, fmt.Errorf("failed to decode public key PEM")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("not an RSA public key")
	}

	// Hash the data
	hash := sha256.Sum256(data)

	// Verify the signature
	err = rsa.VerifyPSS(
		publicKey,
		crypto.SHA256,
		hash[:],
		signatureBytes,
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
		},
	)

	return err == nil, nil
}

// HashData generates SHA-256 hash of data
func (cm *CryptoManager) HashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// GenerateVoterID generates an anonymous voter ID from email
func (cm *CryptoManager) GenerateVoterID(email string, salt string) string {
	if salt == "" {
		salt = "voting-system"
	}
	combined := email + salt
	hash := sha256.Sum256([]byte(combined))
	return fmt.Sprintf("%x", hash[:8]) // Return first 16 hex characters
}

// SecureVoter represents a voter with cryptographic capabilities
type SecureVoter struct {
	VoterID    string
	Name       string
	Email      string
	PrivateKey []byte
	PublicKey  []byte
}

// NewSecureVoter creates a new secure voter
func NewSecureVoter(email, name string) (*SecureVoter, error) {
	cm := NewCryptoManager()

	voterID := cm.GenerateVoterID(email, "")
	privateKey, publicKey, err := cm.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	return &SecureVoter{
		VoterID:    voterID,
		Name:       name,
		Email:      email,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// SignVote signs vote data with the voter's private key
func (v *SecureVoter) SignVote(voteData []byte) (string, error) {
	cm := NewCryptoManager()
	return cm.SignData(voteData, v.PrivateKey)
}

// GetPublicCredentials returns public credentials for registration
func (v *SecureVoter) GetPublicCredentials() map[string]string {
	return map[string]string{
		"voter_id":   v.VoterID,
		"name":       v.Name,
		"public_key": string(v.PublicKey),
	}
}