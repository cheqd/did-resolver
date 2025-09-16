package unit

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "testing"
)

func TestEndpointManager(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "EndpointManager Unit Suite")
}


