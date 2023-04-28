package resource_version_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResourceVersion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[Integration Test]: Resource version query parameter")
}
