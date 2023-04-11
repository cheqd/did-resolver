package collection_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCollection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collection of Resources Integration Tests")
}
