package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type keyPair struct {
	public, private string
}

func genKey(bitSize int) (keyPair, error) {

	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return keyPair{}, err
	}

	pub := key.Public()

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	return keyPair{public: string(pubPEM), private: string(keyPEM)}, nil
}
