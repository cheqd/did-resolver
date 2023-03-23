package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

var _ = Describe("Test GetDIDFragment method", func() {
	DIDDoc := types.NewDidDoc(&validDIDDoc)

	It("can find a existent fragment in VerificationMethod", func() {
		fragmentId := DIDDoc.VerificationMethod[0].Id
		expectedFragment := &DIDDoc.VerificationMethod[0]

		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(fragmentId, DIDDoc)
		Expect(fragment).To(Equal(expectedFragment))
	})

	It("can find a existent fragment in Service", func() {
		fragmentId := DIDDoc.Service[0].Id
		expectedFragment := &DIDDoc.Service[0]

		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(fragmentId, DIDDoc)
		Expect(fragment).To(Equal(expectedFragment))
	})

	It("cannot find a not-existent fragment", func() {
		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(NotExistFragmentId, DIDDoc)
		Expect(fragment).To(BeNil())
	})
})
