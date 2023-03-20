package tests

import (
	"net/http/httptest"

	diddocServices "github.com/cheqd/did-resolver/services/diddoc"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Content/Accept encoding checks", func() {
	var context echo.Context
	var rec *httptest.ResponseRecorder

	BeforeEach(func() {
		validDIDDoc := ValidDIDDoc()
		validMetadata := ValidMetadata()
		validResource := ValidResource()
		ledgerService := NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource)
		context, rec = setupContext(
			"/1.0/identifiers/:did",
			[]string{"did"},
			[]string{ValidDid},
			types.DIDJSON,
			ledgerService)
	})
	Context("Gzip in Accept-Encoding", func() {
		It("should return gzip in Content-Encoding", func() {
			// Setup Accept header to gzip
			context.Request().Header.Set("Accept-Encoding", "gzip")

			err := diddocServices.DidDocEchoHandler(context)
			Expect(err).To(BeNil())

			// Check if Content-Encoding is gzip
			Expect(rec.Header().Get("Content-Encoding")).To(Equal("gzip"))
		})
	})
	Context("Gzip not in Accept-Encoding", func() {
		It("should not return gzip in Content-Encoding", func() {
			err := diddocServices.DidDocEchoHandler(context)
			Expect(err).To(BeNil())

			// Check if Content-Encoding is Empty
			Expect(rec.Header().Get("Content-Encoding")).To(BeEmpty())
		})
	})
	Context("Not supported compressing", func() {
		It("should not return gzip in Content-Encoding", func() {
			// Setup Accept header to gzip
			context.Request().Header.Set("Accept-Encoding", "br")

			err := diddocServices.DidDocEchoHandler(context)
			Expect(err).To(BeNil())

			// Check if Content-Encoding is empty
			Expect(rec.Header().Get("Content-Encoding")).To(BeEmpty())
		})
	})
})
