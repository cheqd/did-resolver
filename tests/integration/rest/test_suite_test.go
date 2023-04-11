//go:build integration

package rest

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestsGeneral(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[Integration Test]: REST API")
}
