package apis_test

import (
	patching "github.com/Microsoft/kunlun/patching"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PasswordGenerator", func() {

	var (
		generator patching.ValueGenerator
	)

	BeforeEach(func() {
		generator = patching.NewPasswordGenerator()
	})
	Describe("Generate", func() {
		Context("Everything OK", func() {
			It("should succeed", func() {
				params := map[string]int{
					"length": 0,
				}
				password, err := generator.Generate(params)
				Expect(err).To(BeNil())
				Expect(password).NotTo(BeNil())
				Expect(len(password.(string))).To(Equal(20))
			})
		})
	})
})
