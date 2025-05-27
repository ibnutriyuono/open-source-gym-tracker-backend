package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)


func EncryptWithPublicKey(data []byte, pubPath string) ([]byte, error) {
	pubPem, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubPem)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key PEM")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func DecryptWithPrivateKey(cipherText []byte, privPath string) ([]byte, error) {
	privPem, err := os.ReadFile(privPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privPem)
	if block == nil {
		return nil, fmt.Errorf("invalid private key PEM: failed to decode")
	}

	var priv *rsa.PrivateKey

	switch block.Type {
	case "RSA PRIVATE KEY":
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#1 private key: %v", err)
		}
	case "PRIVATE KEY":
		privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#8 private key: %v", err)
		}
		var ok bool
		priv, ok = privInterface.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not RSA")
		}
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}