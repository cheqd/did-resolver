//go:build unit

package common

import (
	"sort"
	"time"

	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DidDocMetadataList", func() {
	var (
		versionList types.DidDocMetadataList
		// t0 = utils.MustParseDate("2021-08-23T08:59:59Z")
		tcreated = utils.MustParseDate("2021-08-23T09:00:00Z")
		t1       = utils.MustParseDate("2021-08-23T09:30:00Z")
		t1_2     = utils.MustParseDate("2021-08-23T09:30:01Z")
		t2       = utils.MustParseDate("2021-08-23T09:40:00Z")
		t2_      = utils.MustParseDate("2021-08-23T09:40:01Z")
	)

	BeforeEach(func() {
		// By default we are adding all resources to each version and after that grouping them
		versionList = types.DidDocMetadataList{
			{
				VersionId:   "1",
				Deactivated: false,
				Created:     &tcreated,
				Updated:     nil,
				Resources: []types.DereferencedResource{
					{
						Created: &t1_2,
					},
					{
						Created: &t2_,
					},
				},
			},
			{
				VersionId:   "2",
				Deactivated: false,
				Created:     &tcreated,
				Updated:     &t1,
				Resources: []types.DereferencedResource{
					{
						Created: &t1_2,
					},
					{
						Created: &t2_,
					},
				},
			},
			{
				VersionId:   "3",
				Deactivated: false,
				Created:     &tcreated,
				Updated:     &t2,
				Resources: []types.DereferencedResource{
					{
						Created: &t1_2,
					},
					{
						Created: &t2_,
					},
				},
			},
		}
		sort.Sort(versionList)
	},
	)

	Context("FindBeforeTime", func() {
		// Time right after creation but before first update
		It("should return versionId of metadata with created", func() {
			Expect(versionList.FindActiveForTime(utils.MustParseDate("2021-08-23T09:00:01Z").Format(time.RFC3339))).To(Equal("1"))
		})
		// Time after first update but the the latest
		It("should return versionId of metadata with the first updated", func() {
			Expect(versionList.FindActiveForTime(utils.MustParseDate("2021-08-23T09:30:01Z").Format(time.RFC3339))).To(Equal("2"))
		})

		It("should return versionId of metadata with the latest updated", func() {
			Expect(versionList.FindActiveForTime(utils.MustParseDate("2021-08-23T09:40:01Z").Format(time.RFC3339))).To(Equal("3"))
		})
		// Time before the creation
		It("should return empty string if no metadata found", func() {
			Expect(versionList.FindActiveForTime(utils.MustParseDate("2021-08-23T08:59:59Z").Format(time.RFC3339))).To(Equal(""))
		})
	})

	Context("GetByVersionId", func() {
		It("should return metadata with the given versionId", func() {
			Expect(versionList.GetByVersionId("1")).To(Equal(types.DidDocMetadataList{versionList[2]}))
		})

		It("should return empty list if no metadata found", func() {
			Expect(len(versionList.GetByVersionId("4"))).To(Equal(0))
		})
	})

	Context("SortInDescendingOrder", func() {
		It("should sort metadata in descending order", func() {
			sort.Sort(versionList)
			Expect(versionList[0].VersionId).To(Equal("3"))
			Expect(versionList[1].VersionId).To(Equal("2"))
			Expect(versionList[2].VersionId).To(Equal("1"))
		})
	})

	Context("GetResourcesBeforeNextVersion", func() {
		It("should return empty list of resources for the first version", func() {
			Expect(versionList.GetResourcesBeforeNextVersion("1")).To(Equal(types.DereferencedResourceList{}))
		})
		It("should return resource created before next version", func() {
			Expect(versionList.GetResourcesBeforeNextVersion("2")).To(Equal(types.DereferencedResourceList{
				versionList[1].Resources[0],
			}))
		})

		It("should return list with all resources for the latest versionId", func() {
			Expect(versionList.GetResourcesBeforeNextVersion("3")).To(Equal(types.DereferencedResourceList{
				versionList[0].Resources[0],
				versionList[0].Resources[1],
			}))
		})
	})
})
