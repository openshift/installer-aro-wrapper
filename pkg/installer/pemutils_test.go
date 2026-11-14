package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestPemToPrivateKey(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	rsaPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKey),
	})

	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ecDER, err := x509.MarshalECPrivateKey(ecKey)
	if err != nil {
		t.Fatal(err)
	}
	ecPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: ecDER,
	})

	ecParamsOID, err := asn1.Marshal(ecKey.Curve.Params().Name)
	if err != nil {
		ecParamsOID = []byte("bogus-params")
	}
	ecParamsBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PARAMETERS",
		Bytes: ecParamsOID,
	})

	ecPKCS8DER, err := x509.MarshalPKCS8PrivateKey(ecKey)
	if err != nil {
		t.Fatal(err)
	}
	ecPKCS8PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: ecPKCS8DER,
	})

	testCases := []struct {
		name    string
		data    []byte
		wantRSA bool
		wantEC  bool
		wantErr bool
	}{
		{
			name:    "RSA PKCS1 key",
			data:    rsaPEM,
			wantRSA: true,
		},
		{
			name:   "EC key",
			data:   ecPEM,
			wantEC: true,
		},
		{
			name:   "EC key with EC PARAMETERS prefix",
			data:   append(ecParamsBlock, ecPEM...),
			wantEC: true,
		},
		{
			name:   "EC PKCS8 key",
			data:   ecPKCS8PEM,
			wantEC: true,
		},
		{
			name:   "EC PKCS8 key with EC PARAMETERS prefix",
			data:   append(ecParamsBlock, ecPKCS8PEM...),
			wantEC: true,
		},
		{
			name:    "RSA key with EC PARAMETERS prefix",
			data:    append(ecParamsBlock, rsaPEM...),
			wantRSA: true,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "no key blocks",
			data:    ecParamsBlock,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := pemToPrivateKey(tc.data)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			switch {
			case tc.wantRSA:
				if _, ok := key.(*rsa.PrivateKey); !ok {
					t.Fatalf("expected *rsa.PrivateKey, got %T", key)
				}
			case tc.wantEC:
				if _, ok := key.(*ecdsa.PrivateKey); !ok {
					t.Fatalf("expected *ecdsa.PrivateKey, got %T", key)
				}
			}
		})
	}
}

func TestPemToCertificate(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	nonCertBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PARAMETERS",
		Bytes: []byte("bogus"),
	})

	testCases := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name: "valid certificate",
			data: certPEM,
		},
		{
			name: "certificate preceded by non-cert block",
			data: append(nonCertBlock, certPEM...),
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "no certificate blocks",
			data:    nonCertBlock,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cert, err := pemToCertificate(tc.data)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cert.Subject.CommonName != "test" {
				t.Fatalf("expected CN 'test', got %q", cert.Subject.CommonName)
			}
		})
	}
}
