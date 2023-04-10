//go:build unit

package common

import (
	"net/http"
	"net/http/httptest"

	diddocServices "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Content/Accept encoding checks", func() {
	var context echo.Context
	var rec *httptest.ResponseRecorder

	BeforeEach(func() {
		request := httptest.NewRequest(http.MethodGet, "/1.0/identifiers/"+testconstants.ExistentDid, nil)
		context, rec = utils.SetupEmptyContext(request, types.DIDJSON, utils.MockLedger)
	})
	Context("Gzip in Accept-Encoding", func() {
		It("should return gzip in Content-Encoding", func() {
			// Setup Accept header to gzip
			context.Request().Header.Set(echo.HeaderAcceptEncoding, "gzip")

			err := diddocServices.DidDocEchoHandler(context)
			Expect(err).To(BeNil())

			// Check if Content-Encoding is gzip
			Expect(rec.Header().Get(echo.HeaderContentEncoding)).To(Equal("gzip"))
		})
	})
	Context("Gzip not in Accept-Encoding", func() {
		It("should not return gzip in Content-Encoding", func() {
			err := diddocServices.DidDocEchoHandler(context)
			Expect(err).To(BeNil())

			// Check if Content-Encoding is Empty
			Expect(rec.Header().Get(echo.HeaderContentEncoding)).To(BeEmpty())
		})
	})
	Context("Not supported compressing", func() {
		It("should not return gzip in Content-Encoding", func() {
			// Setup Accept header to gzip
			context.Request().Header.Set(echo.HeaderAcceptEncoding, "br")

			err := diddocServices.DidDocEchoHandler(context)
			Expect(err).To(BeNil())

			// Check if Content-Encoding is empty
			Expect(rec.Header().Get(echo.HeaderContentEncoding)).To(BeEmpty())
		})
	})
	Context("* in Accept-Encoding", func() {
		It("should return gzip in Content-Encoding", func() {
			// Setup Accept header to all possible variants
			context.Request().Header.Set(echo.HeaderAcceptEncoding, "*")

			err := diddocServices.DidDocEchoHandler(context)
			Expect(err).To(BeNil())

			// Check if Content-Encoding is Empty
			Expect(rec.Header().Get(echo.HeaderContentEncoding)).To(Equal("gzip"))
		})
	})
})
