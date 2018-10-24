package generator_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	artifacts "github.com/Microsoft/kunlun/artifacts"
	dep "github.com/Microsoft/kunlun/artifacts/deployments"
	clogger "github.com/Microsoft/kunlun/common/logger"
	"github.com/Microsoft/kunlun/common/storage"
	. "github.com/Microsoft/kunlun/producers/deployment-producer/ashandler/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("AsGenerator", func() {

	var (
		generator   ASGenerator
		hostGroups  []dep.HostGroup
		deployments []dep.Deployment
	)

	BeforeEach(func() {
		log.SetFlags(0)
		fs := afero.NewOsFs()
		afs := &afero.Afero{Fs: fs}

		// Configuration
		logger := clogger.NewLogger(os.Stdout, os.Stdin)
		tempDir, err := ioutil.TempDir("", "")
		fmt.Printf("root folder is %s\n", tempDir)
		Expect(err).To(BeNil())
		stateStore := storage.NewStore(tempDir, afs)

		generator = NewASGenerator(stateStore, logger, afs)
		hostGroupName := "fake_host_group"
		hostGroups = []dep.HostGroup{
			dep.HostGroup{
				Name: hostGroupName,
				Hosts: []dep.Host{
					{
						Alias: "FakeAlias",
					},
				},
			},
		}
		deployments = []dep.Deployment{
			{
				HostGroupName: hostGroupName,
				Vars:          yaml.MapSlice{},
				Roles: []artifacts.Role{
					{
						Name: "kunlun.php",
					},
				},
			},
		}
	})

	Describe("Generate", func() {
		Context("Everything OK", func() {
			It("should succeed", func() {
				Expect(generator.Generate(hostGroups, deployments)).To(BeNil())
			})
		})
	})
})
