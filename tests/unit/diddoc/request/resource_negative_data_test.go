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
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Test resource negative cases. Data returning case", func(testCase ResourceTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	request.Header.Set("Content-Type", string(testCase.resolutionType))
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, MockLedger)
	expectedContentType := types.ContentType(testconstants.ValidResource[0].Metadata.MediaType)

	err := didDocService.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedResource.GetBytes(), rec.Body.Bytes())
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},
	Entry(
		"Negative. ResourceId not found",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s",
				testconstants.ValidDid,
				uuid.New().String(),
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceId wrong format",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s",
				testconstants.ValidDid,
				"SomeNotUUID",
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. resourceVersionTime wrong format",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersionTime=%s",
				testconstants.ValidDid,
				"SomeNotUUID",
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Only ResourceCollectionId is ambiguous query",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceCollectionId=%s",
				testconstants.ValidDid,
				uuid.New().String(),
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceType is not found",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceType=%s",
				testconstants.ValidDid,
				"NotExistentType",
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceName is not found",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceName=%s",
				testconstants.ValidDid,
				"NotExistentName",
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Only ResourceVersion is Ambiguous query",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersion=%s",
				testconstants.ValidDid,
				"NotExistentVersion",
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. checksum wrong",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&checksum=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				"wrongChecksum",
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. checksum query returns multiple resources",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?checksum=%s",
				testconstants.ValidDid,
				fmt.Sprintf("%x", Checksum),
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Only ResourceVersionTime is Ambiguous query",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersionTime=%s",
				testconstants.ValidDid,
				DidDocBeforeCreated.Format(time.RFC3339),
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. There are several resources with the same type but different names",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceType=%s",
				testconstants.ValidDid,
				ResourceType1.Metadata.ResourceType,
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. There are several names with the same name but different types",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceName=%s",
				testconstants.ValidDid,
				ResourceType2.Metadata.Name,
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
