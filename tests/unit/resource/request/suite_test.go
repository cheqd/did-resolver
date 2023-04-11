//go:build unit

package request_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRequest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[Unit Test]: Resource Request Service")
}
