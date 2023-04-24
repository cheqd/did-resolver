//go:build unit

package request

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type QueriesDIDDocTestCase struct {
	didURL             string
	resolutionType     types.ContentType
	expectedResolution types.ResolutionResultI
	expectedError      error
}

type ResourceTestCase struct {
	didURL           string
	resolutionType   types.ContentType
	expectedResource types.ContentStreamI
	expectedError    error
}

var (
	Data                   = []byte("{\"attr\":[\"name\",\"age\"]}")
	Checksum               = sha256.New().Sum(Data)
	VersionId1             = uuid.New().String()
	VersionId2             = uuid.New().String()
	ResourceIdName1        = uuid.New().String()
	ResourceIdName12       = uuid.New().String()
	ResourceIdName2        = uuid.New().String()
	ResourceIdType1        = uuid.New().String()
	ResourceIdType12       = uuid.New().String()
	ResourceIdType2        = uuid.New().String()
	DidDocBeforeCreated, _ = time.Parse(time.RFC3339, "2021-08-23T08:59:00Z")
	DidDocCreated, _       = time.Parse(time.RFC3339, "2021-08-23T09:00:00Z")
	DidDocAfterCreated, _  = time.Parse(time.RFC3339, "2021-08-23T09:00:30Z")
	Resource1Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:01:00Z")
	Resource12Created, _   = time.Parse(time.RFC3339, "2021-08-23T09:01:30Z")
	Resource2Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:02:00Z")
	DidDocUpdated, _       = time.Parse(time.RFC3339, "2021-08-23T09:03:00Z")
	DidDocAfterUpdated, _  = time.Parse(time.RFC3339, "2021-08-23T09:03:30Z")
	Resource3Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:04:00Z")
	Resource34Created, _   = time.Parse(time.RFC3339, "2021-08-23T09:04:30Z")
	Resource4Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:05:00Z")
)

var (
	ResourceName1   = generateResource(ResourceIdName1, "Name1", "string", "1", timestamppb.New(Resource1Created))
	ResourceName12  = generateResource(ResourceIdName12, "Name1", "string", "12", timestamppb.New(Resource12Created))
	ResourceName2   = generateResource(ResourceIdName2, "Name2", "string2", "2", timestamppb.New(Resource2Created))
	ResourceType1   = generateResource(ResourceIdType1, "Name", "Type1", "3", timestamppb.New(Resource3Created))
	ResourceType12  = generateResource(ResourceIdType12, "Name", "Type1", "34", timestamppb.New(Resource34Created))
	ResourceType2   = generateResource(ResourceIdType2, "Name2", "Type2", "4", timestamppb.New(Resource4Created))
	DidDocMetadata1 = generateMetadata(
		VersionId1,
		timestamppb.New(DidDocCreated),
		nil,
	)
	DidDocMetadata2 = generateMetadata(
		VersionId2,
		timestamppb.New(DidDocCreated),
		timestamppb.New(DidDocUpdated),
	)
)

var MockLedger = utils.NewMockLedgerService(
	&testconstants.ValidDIDDoc,
	[]*didTypes.Metadata{
		&DidDocMetadata1,
		&DidDocMetadata2,
	},
	[]resourceTypes.ResourceWithMetadata{
		ResourceName1,
		ResourceName12,
		ResourceName2,
		ResourceType1,
		ResourceType12,
		ResourceType2,
	},
)

func generateResource(resourceId, name, rtype, version string, created *timestamppb.Timestamp) resourceTypes.ResourceWithMetadata {
	return resourceTypes.ResourceWithMetadata{
		Resource: &resourceTypes.Resource{
			Data: Data,
		},
		Metadata: &resourceTypes.Metadata{
			CollectionId: testconstants.ValidIdentifier,
			Id:           resourceId,
			Name:         name,
			ResourceType: rtype,
			MediaType:    "application/json",
			Checksum:     fmt.Sprintf("%x", Checksum),
			Created:      created,
			Version:      version,
		},
	}
}

func generateMetadata(versionId string, created, updated *timestamppb.Timestamp) didTypes.Metadata {
	return didTypes.Metadata{
		VersionId:   versionId,
		Deactivated: false,
		Created:     created,
		Updated:     updated,
	}
}