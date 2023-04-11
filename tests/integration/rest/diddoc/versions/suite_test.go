package versions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVersions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DIDDoc Versions Integration Tests")
}
