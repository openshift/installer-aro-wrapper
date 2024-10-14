package pem

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	utiltls "github.com/openshift/installer-aro-wrapper/pkg/util/tls"
)

func TestEncode(t *testing.T) {
	validCaKey, validCaCerts, err := utiltls.GenerateKeyAndCertificate("validca", nil, nil, true, false)
	assert.NoError(t, err)

	t.Run("encoding key", func(t *testing.T) {
		var b bytes.Buffer
		err := Encode(&b, validCaKey)
		if assert.NoError(t, err) {
			assert.Regexp(t, regexp.MustCompile("-----BEGIN RSA PRIVATE KEY-----\n(?:[a-zA-Z0-9+/=]+\n)*-----END RSA PRIVATE KEY-----\n"), b.String())
		}
	})

	t.Run("encoding single certificate", func(t *testing.T) {
		var b bytes.Buffer
		err := Encode(&b, validCaCerts[0])
		if assert.NoError(t, err) {
			assert.Regexp(t, regexp.MustCompile("-----BEGIN CERTIFICATE-----\n(?:[a-zA-Z0-9+/=]+\n)*-----END CERTIFICATE-----\n"), b.String())
		}
	})

	t.Run("encoding multiple certificates", func(t *testing.T) {
		var b bytes.Buffer
		err := Encode(&b, validCaCerts[0], validCaCerts[0])
		if assert.NoError(t, err) {
			assert.Regexp(t, regexp.MustCompile("(?:-----BEGIN CERTIFICATE-----\n(?:[a-zA-Z0-9+/=]+\n)*-----END CERTIFICATE-----\n){2}"), b.String())
		}
	})

	t.Run("encoding multiple certificates and private key", func(t *testing.T) {
		var b bytes.Buffer
		err := Encode(&b, validCaCerts[0], validCaCerts[0])
		err = errors.Join(err, Encode(&b, validCaKey))
		if assert.NoError(t, err) {
			assert.Regexp(t, regexp.MustCompile("(?:-----BEGIN CERTIFICATE-----\n(?:[a-zA-Z0-9+/=]+\n)*-----END CERTIFICATE-----\n){2}-----BEGIN RSA PRIVATE KEY-----\n(?:[a-zA-Z0-9+/=]+\n)*-----END RSA PRIVATE KEY-----\n"), b.String())
		}
	})
}
