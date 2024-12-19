package mdsd

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import "embed"

//go:embed etc/* systemd/*
var Assets embed.FS
