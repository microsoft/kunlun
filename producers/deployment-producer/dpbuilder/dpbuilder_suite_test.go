package dpbuilder_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDpbuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dpbuilder Suite")
}
