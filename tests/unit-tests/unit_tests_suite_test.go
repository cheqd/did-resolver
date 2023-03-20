package tests_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUnitTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UnitTests Suite")
}
