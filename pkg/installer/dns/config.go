package dns

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

// DNSConfig holds the DNS-related network configuration used by both
// dnsmasq and CustomDNS (CoreDNS) modes in the ARO installer wrapper.
type DNSConfig struct {
	APIIntIP                 string
	IngressIP                string
	GatewayDomains           []string
	GatewayPrivateEndpointIP string
}
