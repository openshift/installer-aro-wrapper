package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func pemToPrivateKey(data []byte) (crypto.PrivateKey, error) {
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}

		switch block.Type {
		case "RSA PRIVATE KEY":
			if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
				return key, nil
			}
		case "EC PRIVATE KEY":
			if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
				return key, nil
			}
		case "PRIVATE KEY":
			if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
				switch key.(type) {
				case *rsa.PrivateKey, *ecdsa.PrivateKey:
					return key, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("data does not contain a valid RSA or ECDSA private key")
}

func pemToCertificate(data []byte) (*x509.Certificate, error) {
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			return x509.ParseCertificate(block.Bytes)
		}
	}
	return nil, fmt.Errorf("data does not contain a valid certificate")
}

func publicKeyToPem(pub crypto.PublicKey) ([]byte, error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: keyBytes}), nil
}
