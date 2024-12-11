package mdsd

import "embed"

//go:embed etc/* systemd/*
var Assets embed.FS
