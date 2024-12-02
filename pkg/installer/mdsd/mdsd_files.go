package mdsd

import "embed"

//go:embed etc/* systemd/*
var assets embed.FS
