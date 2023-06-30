package bootstraplogging

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File generation", Ordered, func() {

	It("generates successfully", func() {

		config := &Config{}

		files, units, err := Files(config)
		Expect(err).To(BeNil())

		Expect(files).To(HaveLen(7))
		Expect(units).To(HaveLen(3))

	})

})

func TestBootstrapLogging(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BootstrapLogging Suite")
}
