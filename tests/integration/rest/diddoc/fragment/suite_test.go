package fragment_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFragment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DIDDoc Fragment Integration Tests")
}
