//go:build integration

package metadata

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMetadata(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[Integration Test]: DIDDoc Metadata query parameter")
}
