package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// condition functions should return an error only if it's not retryable
// if a condition function encounters a retryable error it should return false, nil.

func (m *manager) bootstrapConfigMapReady(ctx context.Context) (bool, error) {
	cm, err := m.kubernetescli.CoreV1().ConfigMaps("kube-system").Get(ctx, "bootstrap", metav1.GetOptions{})
	if err != nil {
		// During bootstrap the control plane nodes may go down while they
		// update various components. In this case, it is usually temporary and
		// we should let it poll until the timeout rather than fail out
		// instantly on such a blip.
		m.log.Printf("bootstrapConfigMapReady condition error %s, continuing to poll", err)
		return false, nil
	}
	return err == nil && cm.Data["status"] == "complete", nil
}
