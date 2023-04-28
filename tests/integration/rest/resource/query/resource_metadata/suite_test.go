package resource_metadata_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResourceMetadata(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[Integration Test]: Resource metadata query parameter")
}
