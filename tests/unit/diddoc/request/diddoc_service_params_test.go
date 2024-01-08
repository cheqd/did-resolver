//go:build unit

package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	didDocService "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Test Query handlers with service and relativeRef params", func(testCase QueriesDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	err := didDocService.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(rec.Code).To(Equal(http.StatusSeeOther))
		Expect(string(testCase.expectedResolution.GetBytes())).To(Equal(context.Response().Header().Get("Location")))
		Expect(err).To(BeNil())
	}
},

	// Positive cases
	Entry(
		"Positive. Service case",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s", testconstants.ValidDid, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. relativeRef case",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionId + Service case",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionId + Service case + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionTime + Service case",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionTime + Service case + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),

	// Negative Cases
	Entry(
		"Negative. Service not found",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s", testconstants.ValidDid, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionId",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionId + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo&versionId=%s", testconstants.ValidDid, testconstants.InvalidServiceId, testconstants.ValidVersionId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionTime",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionTime + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo&versionTime=%s", testconstants.ValidDid, testconstants.InvalidServiceId, testconstants.CreatedAfter.Format(time.RFC3339)),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. RelativeRef without Service",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?relativeRef=%s", testconstants.ValidDid, "blabla"),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewRepresentationNotSupportedError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
)
