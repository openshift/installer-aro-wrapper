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
// Root cause:
//
// WantedBy=network-online.target pulled aro-etchosts-resolver into the network-online.target
// activation chain; combined with Before=node-image-pull.service and node-image-pull's own
// After=network-online.target, that created the circular wait
//
// Fix: WantedBy=network-online.target to WantedBy=multi-user.target and adding After=network-online.target
//
//	node-image-pull        → After  → network-online.target
func TestGenerateEtcHostsAROUnit_OrderingCycle(t *testing.T) {
	unit, err := GenerateEtcHostsAROUnit()
	if err != nil {
		t.Fatalf("GenerateEtcHostsAROUnit() error: %v", err)
	}

	// ── must NOT be present ──────────────────────────────────────────────────
	//
	// If this line appears, dnsmasq will be skipped on RHCOS 252.
	if strings.Contains(unit, "WantedBy=network-online.target") {
		t.Errorf("unit must not contain 'WantedBy=network-online.target': "+
			"this was the directive that caused the ordering cycle\n\ngenerated unit:\n%s", unit)
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

	if !strings.Contains(unit, "StandardOutput=journal+console") {
		t.Errorf(
			"unit MUST contain 'StandardOutput=journal+console': "+
				"output must be routed to journal and console for bootstrap diagnostics\n\n"+
				"generated unit:\n%s", unit)
	}

	if !strings.Contains(unit, "StandardError=journal+console") {
		t.Errorf(
			"unit MUST contain 'StandardError=journal+console': "+
				"errors must be routed to journal and console for bootstrap diagnostics\n\n"+
				"generated unit:\n%s", unit)
	}

	if !strings.Contains(unit, "WantedBy=multi-user.target") {
		t.Errorf(
			"unit MUST contain 'WantedBy=multi-user.target' in [Install]\n\n"+
				"generated unit:\n%s", unit)
	}
}
