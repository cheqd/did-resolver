//go:build unit

package common

import (
	"time"

	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DereferencingContentStream Find before time", func() {
	var (
		versionList types.DereferencedResourceList
		t0          = utils.MustParseDate("2021-08-23T08:59:59Z")
		t1          = utils.MustParseDate("2021-08-23T09:00:00Z")
		t1_2        = utils.MustParseDate("2021-08-23T09:30:01Z")
		t2          = utils.MustParseDate("2021-08-23T09:40:00Z")
		t2_3        = utils.MustParseDate("2021-08-23T09:40:01Z")
		t3          = utils.MustParseDate("2021-08-23T09:50:00Z")
		t3_         = utils.MustParseDate("2021-08-23T09:50:01Z")
	)

	BeforeEach(func() {
		versionList = types.DereferencedResourceList{
			{
				Created:      &t1,
				ResourceId:   "r1",
				CollectionId: "c1",
				Name:         "name1",
				ResourceType: "type1",
				Version:      "v1",
				Checksum:     "checksum1",
			},
			{
				Created:      &t2,
				ResourceId:   "r2",
				CollectionId: "c2",
				Name:         "name2",
				ResourceType: "type2",
				Version:      "v2",
				Checksum:     "checksum2",
			},
			{
				Created:      &t3,
				ResourceId:   "r3",
				CollectionId: "c3",
				Name:         "name2",
				ResourceType: "type2",
				Version:      "v2",
				Checksum:     "checksum3",
			},
		}
	})

	Context("FindBeforeTime", func() {
		// Should return the first resource
		It("should return resourceId of the first resource", func() {
			Expect(versionList.FindBeforeTime(t1_2.Format(time.RFC3339))).To(Equal("r1"))
		})
		// Should return the second resource
		It("should return resourceId of the second resource", func() {
			Expect(versionList.FindBeforeTime(t2_3.Format(time.RFC3339))).To(Equal("r2"))
		})
		// Should return the latest resource
		It("should return resourceId of the latest resource", func() {
			Expect(versionList.FindBeforeTime(t3_.Format(time.RFC3339))).To(Equal("r3"))
		})
		// Time before the creation
		It("should return empty string if no metadata found", func() {
			Expect(versionList.FindBeforeTime(t0.Format(time.RFC3339))).To(Equal(""))
		})
	})

	Context("FindAllBeforeTime", func() {
		It("should return 1 resource before the given time", func() {
			resources, err := versionList.FindAllBeforeTime(t1_2.Format(time.RFC3339))
			Expect(err).To(BeNil())
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].ResourceId).To(Equal("r1"))
		})
		It("should return 2 resources before the given time", func() {
			resources, err := versionList.FindAllBeforeTime(t2_3.Format(time.RFC3339))
			Expect(err).To(BeNil())
			Expect(len(resources)).To(Equal(2))
			Expect(resources[1].ResourceId).To(Equal("r1"))
			Expect(resources[0].ResourceId).To(Equal("r2"))
		})
		It("should return 3 resources before the given time", func() {
			resources, err := versionList.FindAllBeforeTime(t3_.Format(time.RFC3339))
			Expect(err).To(BeNil())
			Expect(len(resources)).To(Equal(3))
			Expect(resources[2].ResourceId).To(Equal("r1"))
			Expect(resources[1].ResourceId).To(Equal("r2"))
			Expect(resources[0].ResourceId).To(Equal("r3"))
		})
		It("should return empty list if no metadata found", func() {
			resources, err := versionList.FindAllBeforeTime(t0.Format(time.RFC3339))
			Expect(err).To(BeNil())
			Expect(len(resources)).To(Equal(0))
		})
	})
	Context("FilterByCollectionId", func() {
		It("should return list with one resource and the given collectionId", func() {
			resources := versionList.FilterByCollectionId("c1")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].CollectionId).To(Equal("c1"))
		})
		It("should return list with one resource and the given collectionId", func() {
			resources := versionList.FilterByCollectionId("c2")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].CollectionId).To(Equal("c2"))
		})
		It("should return empty list cause collectionId is not placed", func() {
			resources := versionList.FilterByCollectionId("c0")
			Expect(len(resources)).To(Equal(0))
		})
	})
	Context("GetByResourceId", func() {
		It("should return list with one resource and the given resourceId", func() {
			resources := versionList.GetByResourceId("r1")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].ResourceId).To(Equal("r1"))
		})
		It("should return list with one resource and the given resourceId", func() {
			resources := versionList.GetByResourceId("r2")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].ResourceId).To(Equal("r2"))
		})
		It("should return empty list cause resourceId is not placed", func() {
			resources := versionList.GetByResourceId("r0")
			Expect(len(resources)).To(Equal(0))
		})
	})
	Context("FilterByResourceName", func() {
		It("should return list with one resource and the given resourceName", func() {
			resources := versionList.FilterByResourceName("name1")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].Name).To(Equal("name1"))
		})
		It("should return list with 2 resources and the given resourceName", func() {
			resources := versionList.FilterByResourceName("name2")
			Expect(len(resources)).To(Equal(2))
			Expect(resources[0].Name).To(Equal("name2"))
			Expect(resources[1].Name).To(Equal("name2"))
		})
		It("should return empty list cause resourceName is not placed", func() {
			resources := versionList.FilterByResourceName("name0")
			Expect(len(resources)).To(Equal(0))
		})
	})
	Context("FilterByResourceType", func() {
		It("should return list with one resource and the given resourceType", func() {
			resources := versionList.FilterByResourceType("type1")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].ResourceType).To(Equal("type1"))
		})
		It("should return list with 2 resources and the given resourceType", func() {
			resources := versionList.FilterByResourceType("type2")
			Expect(len(resources)).To(Equal(2))
			Expect(resources[0].ResourceType).To(Equal("type2"))
			Expect(resources[1].ResourceType).To(Equal("type2"))
		})
		It("should return empty list cause resourceType is not placed", func() {
			resources := versionList.FilterByResourceType("type0")
			Expect(len(resources)).To(Equal(0))
		})
	})
	Context("FilterByResourceVersion", func() {
		It("should return list with one resource and the given resourceVersion", func() {
			resources := versionList.FilterByResourceVersion("v1")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].Version).To(Equal("v1"))
		})
		It("should return list with 2 resources and the given resourceVersion", func() {
			resources := versionList.FilterByResourceVersion("v2")
			Expect(len(resources)).To(Equal(2))
			Expect(resources[0].Version).To(Equal("v2"))
			Expect(resources[1].Version).To(Equal("v2"))
		})
		It("should return empty list cause resourceVersion is not placed", func() {
			resources := versionList.FilterByResourceVersion("v0")
			Expect(len(resources)).To(Equal(0))
		})
	})
	Context("FilterByChecksum", func() {
		It("should return list with one resource and the given checksum", func() {
			resources := versionList.FilterByChecksum("checksum1")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].Checksum).To(Equal("checksum1"))
		})
		It("should return list with 2 resources and the given checksum", func() {
			resources := versionList.FilterByChecksum("checksum2")
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].Checksum).To(Equal("checksum2"))
		})
		It("should return empty list cause checksum is not placed", func() {
			resources := versionList.FilterByChecksum("checksum0")
			Expect(len(resources)).To(Equal(0))
		})
	})
	Context("AreResourceNamesTheSame", func() {
		It("should return true", func() {
			versionList[0].Name = "name2"
			Expect(versionList.AreResourceNamesTheSame()).To(BeTrue())
		})
		It("should return false", func() {
			Expect(versionList.AreResourceNamesTheSame()).To(BeFalse())
		})
	})
	Context("AreResourceTypesTheSame", func() {
		It("should return true", func() {
			versionList[0].ResourceType = "type2"
			Expect(versionList.AreResourceTypesTheSame()).To(BeTrue())
		})
		It("should return false", func() {
			Expect(versionList.AreResourceTypesTheSame()).To(BeFalse())
		})
	})
})
