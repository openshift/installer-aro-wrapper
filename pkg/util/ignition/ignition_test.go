package ignition

import (
	"embed"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/*
var testData embed.FS

func TestIgnition(t *testing.T) {
	files, units, err := GetFiles(
		testData, map[string]string{"name": "world"}, map[string]bool{
			"testenabledunit.service":  true,
			"testdisabledunit.service": false,
		}, map[string]int{
			"/etc/NetworkManager/config.d/90-something-else": 0755,
		})
	assert.NoError(t, err, "GetFiles")
	assert.Len(t, units, 2)

	for _, u := range units {
		switch u.Name {
		case "testenabledunit.service":
			assert.Equal(t, true, *u.Enabled)
			assert.Equal(t, *u.Contents, "not entirely empty\n")
		case "testdisabledunit.service":
			assert.Equal(t, false, *u.Enabled)
			assert.Equal(t, *u.Contents, "hello world\n")
		default:
			assert.Failf(t, "unknown unit", "not known: %s", u.Name)
		}
	}

	assert.Len(t, files, 3)
	for _, f := range files {
		switch f.Path {
		case "/etc/NetworkManager/config.d/90-something-else":
			assert.Equal(t, 0755, *f.Mode)
			assert.Equal(
				t, *f.Contents.Source,
				"data:text/plain;charset=utf-8;base64,"+base64.StdEncoding.EncodeToString([]byte("1234\n")))
		case "/etc/NetworkManager/something.conf":
			assert.Equal(t, 0555, *f.Mode)
			assert.Equal(
				t, *f.Contents.Source,
				"data:text/plain;charset=utf-8;base64,"+base64.StdEncoding.EncodeToString([]byte("test\n")))
		case "/etc/motd":
			assert.Equal(t, 0555, *f.Mode)
			assert.Equal(
				t, *f.Contents.Source,
				"data:text/plain;charset=utf-8;base64,"+base64.StdEncoding.EncodeToString([]byte("hello :)\n")))
		default:
			assert.Failf(t, "unknown path", "not known: %s", f.Path)
		}
	}
}
