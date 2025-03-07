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

var _ = DescribeTable("Test Query handlers with resource params. Returns Resource", func(testCase ResourceTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	request.Header.Set("Accept", string(testCase.resolutionType))
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

	// Positive cases
	Entry(
		"Positive. ResourceId case. The first item",
		ResourceTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s?resourceId=%s", testconstants.ValidDid, ResourceIdName1),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId case. The first item",
		ResourceTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s", testconstants.ValidDid, ResourceIdName1, testconstants.ValidIdentifier),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				testconstants.ValidIdentifier,
				ResourceName1.Metadata.Name),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				testconstants.ValidIdentifier,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType + ResourceVersion case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s",
				testconstants.ValidDid,
				ResourceIdName2,
				testconstants.ValidIdentifier,
				ResourceName2.Metadata.Name,
				ResourceName2.Metadata.ResourceType,
				ResourceName2.Metadata.Version),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName2.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceVersionTime=%s",
				testconstants.ValidDid,
				ResourceIdName2,
				testconstants.ValidIdentifier,
				ResourceName2.Metadata.Name,
				ResourceName2.Metadata.ResourceType,
				ResourceName2.Metadata.Version,
				DidDocUpdated.Format(time.RFC3339)),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName2.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime + Checksum case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceVersionTime=%s&checksum=%s",
				testconstants.ValidDid,
				ResourceIdName2,
				testconstants.ValidIdentifier,
				ResourceName2.Metadata.Name,
				ResourceName2.Metadata.ResourceType,
				ResourceName2.Metadata.Version,
				DidDocUpdated.Format(time.RFC3339),
				ResourceName2.Metadata.Checksum),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName2.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. Unique Checksum case.",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?checksum=%s",
				testconstants.ValidDid,
				ResourceChecksum.Metadata.Checksum),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceChecksum.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceType case. Get the latest with the same resource type",
		ResourceTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s?resourceType=%s", testconstants.ValidDid, ResourceName1.Metadata.ResourceType),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceName case. Get the latest with the same resource name",
		ResourceTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s?resourceName=%s", testconstants.ValidDid, ResourceType1.Metadata.Name),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceType12.Resource.Data),
			expectedError:    nil,
		},
	),
)
