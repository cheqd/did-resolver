package resource_id_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResourceId(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[Integration Test]:  ResourceId query parameter")
}
