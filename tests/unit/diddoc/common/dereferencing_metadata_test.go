//go:build unit

package common

import (
	"time"

	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DereferencingMetadata", func() {
	var versionList types.DidDocMetadataList

	BeforeEach(func() {
		_tcreated := utils.MustParseDate("2021-08-23T09:00:00Z")
		_t1 := utils.MustParseDate("2021-08-23T09:30:00Z")
		_t2 := utils.MustParseDate("2021-08-23T09:40:00Z")
		versionList = types.DidDocMetadataList{
			{
				VersionId:   "1",
				Deactivated: false,
				Created:     &_tcreated,
				Updated:     nil,
			},
			{
				VersionId:   "2",
				Deactivated: false,
				Created:     &_tcreated,
				Updated:     &_t1,
			},
			{
				VersionId:   "3",
				Deactivated: false,
				Created:     &_tcreated,
				Updated:     &_t2,
			},
		}
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
})
