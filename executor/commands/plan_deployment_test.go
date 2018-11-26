package commands_test

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	// testInfra "github.com/kun-lun/test-infra/pkg/apis"
)

var _ = Describe("PlanDeployment", func() {

	// var (
	// 	testInf testInfra.TestInfra
	// )
	// BeforeEach(func() {
	// 	testInf = testInfra.TestInfra{}
	// })
	// AfterEach(func() {
	// 	// os.RemoveAll(stateDir)
	// })
	// Describe("Execute", func() {
	// 	Context("when everything ok", func() {
	// 		It("should prepare deployment successfully", func() {
	// 			store := testInf.PrepareForDeploymentCmd()
	// 			fs := afero.NewOsFs()
	// 			afs := &afero.Afero{Fs: fs}

	// 			patch := patching.NewPatching(store, afs)
	// 			manifest, _ := patch.ProvisionManifest()
	// 			content, err := manifest.ToYAML()
	// 			fmt.Printf("%s", string(content))
	// 			Expect(err).To(BeNil())
	// 		})
	// 	})
	// })
})
