package apis_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPatching(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Patching Suite")
}
