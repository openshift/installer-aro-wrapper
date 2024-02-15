package steps

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/sirupsen/logrus"

	testlog "github.com/openshift/ARO-Installer/test/util/log"
)

func successfulFunc(context.Context) error { return nil }
func failingFunc(context.Context) error    { return errors.New("oh no!") }

func alwaysFalseCondition(context.Context) (bool, error) { return false, nil }
func alwaysTrueCondition(context.Context) (bool, error)  { return true, nil }
func timingOutCondition(ctx context.Context) (bool, error) {
	time.Sleep(60 * time.Millisecond)
	return false, nil
}
func internalTimeoutCondition(ctx context.Context) (bool, error) {
	return false, context.DeadlineExceeded
}

func TestStepRunner(t *testing.T) {
	for _, tt := range []struct {
		name        string
		steps       func(*gomock.Controller) []Step
		wantEntries []map[string]types.GomegaMatcher
		wantErr     string
	}{
		{
			name: "All successful Actions will have a successful run",
			steps: func(controller *gomock.Controller) []Step {
				return []Step{
					Action(successfulFunc),
					Action(successfulFunc),
					Action(successfulFunc),
				}
			},
			wantEntries: []map[string]types.GomegaMatcher{
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
			},
		},
		{
			name: "A failing Action will fail the run",
			steps: func(controller *gomock.Controller) []Step {
				return []Step{
					Action(successfulFunc),
					Action(failingFunc),
					Action(successfulFunc),
				}
			},
			wantEntries: []map[string]types.GomegaMatcher{
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.failingFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal(`step [Action github.com/openshift/ARO-Installer/pkg/util/steps.failingFunc] encountered error: oh no!`),
					"level": gomega.Equal(logrus.ErrorLevel),
				},
			},
			wantErr: `oh no!`,
		},
		{
			name: "A successful condition will allow steps to continue",
			steps: func(controller *gomock.Controller) []Step {
				return []Step{
					Action(successfulFunc),
					Condition(alwaysTrueCondition, 50*time.Millisecond, true),
					Action(successfulFunc),
				}
			},
			wantEntries: []map[string]types.GomegaMatcher{
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.alwaysTrueCondition, timeout 50ms]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
			},
		},
		{
			name: "A failed condition with fail=false will allow steps to continue",
			steps: func(controller *gomock.Controller) []Step {
				return []Step{
					Action(successfulFunc),
					Condition(alwaysFalseCondition, 50*time.Millisecond, false),
					Action(successfulFunc),
				}
			},
			wantEntries: []map[string]types.GomegaMatcher{
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.alwaysFalseCondition, timeout 50ms]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.alwaysFalseCondition, timeout 50ms] failed but has configured 'fail=false'. Continuing. Error: context deadline exceeded"),
					"level": gomega.Equal(logrus.WarnLevel),
				},
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
			},
		},
		{
			name: "A timed out Condition causes a failure",
			steps: func(controller *gomock.Controller) []Step {
				return []Step{
					Action(successfulFunc),
					&conditionStep{
						f:            timingOutCondition,
						fail:         true,
						pollInterval: 20 * time.Millisecond,
						timeout:      50 * time.Millisecond,
					},
					Action(successfulFunc),
				}
			},
			wantEntries: []map[string]types.GomegaMatcher{
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.timingOutCondition, timeout 50ms]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.timingOutCondition, timeout 50ms] encountered error: context deadline exceeded"),
					"level": gomega.Equal(logrus.ErrorLevel),
				},
			},
			wantErr: "context deadline exceeded",
		},
		{
			name: "A Condition that returns a timeout error causes a different failure from a timed out Condition",
			steps: func(controller *gomock.Controller) []Step {
				return []Step{
					Action(successfulFunc),
					&conditionStep{
						f:            internalTimeoutCondition,
						fail:         true,
						pollInterval: 20 * time.Millisecond,
						timeout:      50 * time.Millisecond,
					},
					Action(successfulFunc),
				}
			},
			wantEntries: []map[string]types.GomegaMatcher{
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.internalTimeoutCondition, timeout 50ms]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.internalTimeoutCondition, timeout 50ms] encountered error: condition encountered internal timeout: context deadline exceeded"),
					"level": gomega.Equal(logrus.ErrorLevel),
				},
			},
			wantErr: "condition encountered internal timeout: context deadline exceeded",
		},
		{
			name: "A Condition that does not return true in the timeout time causes a failure",
			steps: func(controller *gomock.Controller) []Step {
				return []Step{
					Action(successfulFunc),
					Condition(alwaysFalseCondition, 50*time.Millisecond, true),
					Action(successfulFunc),
				}
			},
			wantEntries: []map[string]types.GomegaMatcher{
				{
					"msg":   gomega.Equal("running step [Action github.com/openshift/ARO-Installer/pkg/util/steps.successfulFunc]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("running step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.alwaysFalseCondition, timeout 50ms]"),
					"level": gomega.Equal(logrus.InfoLevel),
				},
				{
					"msg":   gomega.Equal("step [Condition github.com/openshift/ARO-Installer/pkg/util/steps.alwaysFalseCondition, timeout 50ms] encountered error: context deadline exceeded"),
					"level": gomega.Equal(logrus.ErrorLevel),
				},
			},
			wantErr: "context deadline exceeded",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			controller := gomock.NewController(t)
			defer controller.Finish()

			h, log := testlog.New()
			steps := tt.steps(controller)

			err := Run(ctx, log, 25*time.Millisecond, steps)
			if err != nil && err.Error() != tt.wantErr ||
				err == nil && tt.wantErr != "" {
				t.Error(err)
			}

			err = testlog.AssertLoggingOutput(h, tt.wantEntries)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
