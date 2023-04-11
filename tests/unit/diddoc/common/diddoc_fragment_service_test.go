//go:build unit

package common

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
)

var _ = Describe("Test GetDIDFragment method", func() {
	It("can find a existent fragment in VerificationMethod", func() {
		fragmentId := testconstants.ValidDIDDocResolution.VerificationMethod[0].Id
		expectedFragment := &testconstants.ValidDIDDocResolution.VerificationMethod[0]

		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(fragmentId, testconstants.ValidDIDDocResolution)
		Expect(fragment).To(Equal(expectedFragment))
	})

	It("can find a existent fragment in Service", func() {
		fragmentId := testconstants.ValidDIDDocResolution.Service[0].Id
		expectedFragment := &testconstants.ValidDIDDocResolution.Service[0]

		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(fragmentId, testconstants.ValidDIDDocResolution)
		Expect(fragment).To(Equal(expectedFragment))
	})

	It("cannot find a not-existent fragment", func() {
		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(testconstants.NotExistentFragment, testconstants.ValidDIDDocResolution)
		Expect(fragment).To(BeNil())
	})
})
