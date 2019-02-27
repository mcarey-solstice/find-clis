package main_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mcarey-solstice/find-clis"
)

var _ = Describe("FindClis", func() {

	const (
		filename = "jumpbox-cli"
		root = "fixtures"
	)

	It("Should find all clis", func() {
		expected := []string{"/var/vcap/packages/bosh/bin/bosh", "/var/vcap/packages/jq/bin/jq", "/var/vcap/packages/cf/bin/cf", "/var/vcap/packages/cf-mgmt/bin/cf-mgmt", "/var/vcap/packages/cf-mgmt/bin/cf-mgmt-config"}

		actual, err := FindClis(filename, root)
		Expect(err).To(BeNil())

		sort.Strings(actual)
		sort.Strings(expected)

		Expect(actual).To(Equal(expected))
	})

	It("Should return the correct paths", func() {
		expected := []string{"/var/vcap/packages/bosh/bin", "/var/vcap/packages/jq/bin", "/var/vcap/packages/cf/bin", "/var/vcap/packages/cf-mgmt/bin"}

		actual, err := FindPaths(filename, root)
		Expect(err).To(BeNil())

		sort.Strings(actual)
		sort.Strings(expected)

		Expect(actual).To(Equal(expected))
	})
})
