package unit

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEndpointManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EndpointManager Unit Suite")
}
