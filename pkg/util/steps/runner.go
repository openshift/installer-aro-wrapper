package steps

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/ARO-Installer/pkg/api"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// FriendlyName returns a "friendly" stringified name of the given func.
func FriendlyName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// Step is the interface for steps that Runner can execute.
type Step interface {
	run(ctx context.Context, log *logrus.Entry) error
	String() string
}

// Run executes the provided steps in order until one fails or all steps
// are completed. Errors from failed steps are returned directly.
func Run(ctx context.Context, log *logrus.Entry, pollInterval time.Duration, steps []Step) error {
	for _, step := range steps {
		log.Infof("running step %s", step)
		err := step.run(ctx, log)

		if err != nil {
			log.Errorf("step %s encountered error: %s", step, err.Error())
			if strings.Contains(err.Error(), "ThrottlingLimitExceeded") || strings.Contains(err.Error(), "TooManyRequests") {
				correlation_id := fmt.Sprintf("%v", log.Data["correlation_id"])
				message := fmt.Sprintf("Requests are being throttled due to Azure Storage limits being exceeded. Please visit https://learn.microsoft.com/en-us/azure/openshift/troubleshoot#exceeding-azure-storage-limits for more details. CorrelationId: " + correlation_id)
				err = api.NewCloudError(
					http.StatusTooManyRequests,
					api.CloudErrorCodeThrottlingLimitExceeded,
					"",
					message)
				return err
			}

			if err, ok := err.(stackTracer); ok {
				trace := ""
				for _, f := range err.StackTrace() {
					trace = trace + fmt.Sprintf("%+s:%d\n", f, f)
				}
				log.Error(trace)
			}

			return err
		}
	}
	return nil
}
