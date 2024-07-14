package sealedsecret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
)

func (s SealedSecretService) hybridEncrypt(value, label string) (string, error) {
	// Generate a random AES key
	aesKey := make([]byte, 32) // Using AES-256
	if _, err := rand.Read(aesKey); err != nil {
		return "", err
	}

	// Prepare AES cipher
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	// Prepare GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Because we use a random AES key for each encryption, we can use a zero nonce.
	nonce := make([]byte, gcm.NonceSize())

	// Encrypt the data using AES-GCM
	cipherText := gcm.Seal(nil, nonce, []byte(value), nil)

	// Encrypt the AES key using RSA-OAEP
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, s.pubKey, aesKey, []byte(label))
	if err != nil {
		return "", err
	}

	// Prepend the length of the RSA encrypted key to the ciphertext
	lenRSA := uint16(len(encryptedKey))
	fullCipherText := make([]byte, 2+len(encryptedKey)+len(cipherText))
	binary.BigEndian.PutUint16(fullCipherText, lenRSA)
	copy(fullCipherText[2:], encryptedKey)
	copy(fullCipherText[2+len(encryptedKey):], cipherText)

	// Encode the result in base64
	encodedResult := base64.StdEncoding.EncodeToString(fullCipherText)

	return encodedResult, nil
}
