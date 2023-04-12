//go:build unit

package common

import (
	"time"

	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	utils "github.com/cheqd/did-resolver/tests/unit"
)

var _ = Describe("DereferencingContentStream Find before time", func() {
	var versionList types.DereferencedResourceList

	BeforeEach(func() {
		_tcreated :=utils.MustParseDate("2021-08-23T09:00:00Z")
		_t1 := utils.MustParseDate("2021-08-23T09:30:00Z")
		_t2 := utils.MustParseDate("2021-08-23T09:40:00Z")
		versionList = types.DereferencedResourceList{
			{
				Created:    &_tcreated,
				ResourceId: "1",
			},
			{
				Created:    &_t1,
				ResourceId: "2",
			},
			{
				Created:    &_t2,
				ResourceId: "3",
			},
		}
	})

	Context("FindBeforeTime", func() {
		// Should return the first resource
		It("should return resourceId of the first resource", func() {
			Expect(versionList.FindBeforeTime(utils.MustParseDate("2021-08-23T09:00:01Z").Format(time.RFC3339))).To(Equal("1"))
		})
		// Should return the second resource
		It("should return resourceId of the second resource", func() {
			Expect(versionList.FindBeforeTime(utils.MustParseDate("2021-08-23T09:30:01Z").Format(time.RFC3339))).To(Equal("2"))
		})
		// Should return the latest resource
		It("should return resourceId of the latest resource", func() {
			Expect(versionList.FindBeforeTime(utils.MustParseDate("2021-08-23T09:40:01Z").Format(time.RFC3339))).To(Equal("3"))
		})
		// Time before the creation
		It("should return empty string if no metadata found", func() {
			Expect(versionList.FindBeforeTime(utils.MustParseDate("2021-08-23T08:59:59Z").Format(time.RFC3339))).To(Equal(""))
		})
	})
})
