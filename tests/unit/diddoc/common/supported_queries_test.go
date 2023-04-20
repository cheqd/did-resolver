//go:build unit

package common

import (
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diff with values tests", func() {
	It("returns empty list if all values are supported", func() {
		supportedQueries := types.SupportedQueriesT{"a", "b", "c"}
		values := map[string][]string{
			"a": {"1"},
			"b": {"2"},
			"c": {"3"},
		}

		result := supportedQueries.DiffWithUrlValues(values)
		Expect(result).To(BeEmpty())
	})

	It("returns list with unsupported values", func() {
		supportedQueries := types.SupportedQueriesT{"a", "b", "c"}
		values := map[string][]string{
			"a": {"1"},
			"b": {"2"},
			"c": {"3"},
			"d": {"4"},
		}

		result := supportedQueries.DiffWithUrlValues(values)
		Expect(result).To(Equal([]string{"d"}))
	})
})
