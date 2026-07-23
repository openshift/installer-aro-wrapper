package etchost

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"strings"
	"testing"
)

// TestGenerateEtcHostsAROUnit_OrderingCycle is the regression test for the
// ordering-cycle bug that caused OSProvisioningTimedOut on UDR clusters.
//
// Root cause: aro-etchosts-resolver.service had BOTH
//
//	Before=network-online.target   (in [Unit])
//	WantedBy=network-online.target (in [Install])
//
// On RHCOS (systemd 252) this creates a dependency cycle:
//
//	network-online.target wants  aro-etchosts-resolver  → activates it
//	network-online.target waits  aro-etchosts-resolver  → must finish first (Before=)
//
// systemd 252 breaks the cycle by skipping dnsmasq.service entirely, which
// logs "[ SKIP ] Ordering cycle found, skipping DNS caching server."  Without
// dnsmasq, /etc/resolv.conf keeps the raw DHCP nameserver (172.16.0.0), so
// arostgsvc.azurecr.io resolves to the public IP (40.64.135.171) which a UDR
// blackholes → node-image-pull loops → bootstrap never becomes healthy →
// all masters: OSProvisioningTimedOut.
//
// Fix: remove Before=network-online.target.  The ordering is preserved
// transitively:  aro-etchosts-resolver → Before → node-image-pull
//
//	node-image-pull        → After  → network-online.target
func TestGenerateEtcHostsAROUnit_OrderingCycle(t *testing.T) {
	unit, err := GenerateEtcHostsAROUnit()
	if err != nil {
		t.Fatalf("GenerateEtcHostsAROUnit() error: %v", err)
	}

	// ── must NOT be present ──────────────────────────────────────────────────
	//
	// Before=network-online.target + WantedBy=network-online.target is the
	// cycle.  If this line appears, dnsmasq will be skipped on RHCOS 252.
	if strings.Contains(unit, "Before=network-online.target") {
		t.Errorf(
			"unit MUST NOT contain 'Before=network-online.target': "+
				"combined with WantedBy=network-online.target this creates an "+
				"ordering cycle that causes systemd 252 to skip dnsmasq.service\n\n"+
				"generated unit:\n%s", unit)
	}

	// ── must be present ──────────────────────────────────────────────────────

	if !strings.Contains(unit, "Type=oneshot") {
		t.Errorf(
			"unit MUST contain 'Type=oneshot': the service runs a script and exits\n\n"+
				"generated unit:\n%s", unit)
	}

	if !strings.Contains(unit, "RemainAfterExit=yes") {
		t.Errorf(
			"unit MUST contain 'RemainAfterExit=yes': keeps service active after "+
				"execution so 'systemctl status' shows it ran successfully\n\n"+
				"generated unit:\n%s", unit)
	}

	if !strings.Contains(unit, "Before=node-image-pull.service") {
		t.Errorf(
			"unit MUST contain 'Before=node-image-pull.service': "+
				"ordering to node-image-pull is preserved after removing Before=network-online.target\n\n"+
				"generated unit:\n%s", unit)
	}

	if !strings.Contains(unit, "WantedBy=multi-user.target") {
		t.Errorf(
			"unit MUST contain 'WantedBy=multi-user.target' in [Install]\n\n"+
				"generated unit:\n%s", unit)
	}
}
