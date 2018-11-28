package artifacts_test

import (
	"encoding/json"
	"reflect"

	"gopkg.in/yaml.v2"

	. "github.com/Microsoft/kunlun/artifacts"
	"github.com/go-test/deep"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest", func() {

	var (
		m *Manifest
	)
	BeforeEach(func() {

		platform := Platform{
			Type: "php",
		}

		networks := []VirtualNetwork{
			{
				Name:         "vnet-1",
				AddressSpace: "10.0.0.0/16",
				Subnets: []Subnet{
					{
						Range:   "10.10.0.0/24",
						Gateway: "10.10.0.1",
						Name:    "snet-1",
					},
				},
			}}

		loadBalancers := []LoadBalancer{
			{
				Name: "kunlun-wenserver-lb",
				SKU:  "Standard",
				BackendAddressPools: []LoadBalancerBackendAddressPool{
					{
						Name: "backend-address-pool-1",
					},
				},
				HealthProbes: []LoadBalancerHealthProbe{
					{
						Name:        "http-probe",
						Protocol:    "Http",
						Port:        80,
						RequestPath: "/",
					},
				},
				Rules: []LoadBalancerRule{
					{
						Name:                   "http_rule",
						Protocol:               "Tcp",
						FrontendPort:           80,
						BackendPort:            80,
						BackendAddressPoolName: "backend-address-pool-1",
						HealthProbeName:        "http-probe",
					},
				},
			},
		}

		networkSecurityGroups := []NetworkSecurityGroup{
			{
				Name: "nsg_1",
				NetworkSecurityRules: []NetworkSecurityRule{
					{
						Name:                     "allow-ssh",
						Priority:                 100,
						Direction:                "Inbound",
						Access:                   "Allow",
						Protocol:                 "Tcp",
						SourcePortRange:          "*",
						DestinationPortRange:     "22",
						SourceAddressPrefix:      "*",
						DestinationAddressPrefix: "*",
					},
				},
			},
		}

		vmGroups := []VMGroup{
			{
				Name: "jumpbox",
				Meta: yaml.MapSlice{
					{
						Key:   "group_type",
						Value: "jumpbox",
					},
				},
				SKU:   VMStandardDS1V2,
				Count: 1,
				Type:  "vm",
				Storage: &VMStorage{
					Image: &Image{
						Offer:     "offer1",
						Publisher: "ubuntu",
						SKU:       "sku1",
						Version:   "latest",
					},
					OSDisk: &OSDisk{
						ManagedDiskType: "Standard_LRS",
						Caching:         "ReadWrite",
						CreateOption:    "FromImage",
					},
					DataDisks: []DataDisk{
						{
							ManagedDiskType: "Standard_LRS",
							Caching:         "ReadWrite",
							CreateOption:    "FromImage",
							DiskSizeGB:      100,
						},
					},
					AzureFiles: []AzureFile{
						{
							StorageAccount: "storage_account_1",
							Name:           "azure_file_1",
							MountPoint:     "/mnt/azurefile_1",
						},
					},
				},
				NetworkInfos: []VMNetworkInfo{
					{
						SubnetName:                         networks[0].Subnets[0].Name,
						LoadBalancerName:                   loadBalancers[0].Name,
						LoadBalancerBackendAddressPoolName: loadBalancers[0].BackendAddressPools[0].Name,
						NetworkSecurityGroupName:           networkSecurityGroups[0].Name,
						PublicIP:                           "dynamic",
					},
				},
				Roles: []Role{},
			},
			{
				Name:  "d2v3_group",
				SKU:   VMStandardDS1V2,
				Count: 2,
				Type:  "vm",
				Storage: &VMStorage{
					OSDisk: &OSDisk{
						ManagedDiskType: "Standard_LRS",
						Caching:         "ReadWrite",
						CreateOption:    "FromImage",
					},
					DataDisks: []DataDisk{
						{
							ManagedDiskType: "Standard_LRS",
							Caching:         "ReadWrite",
							CreateOption:    "FromImage",
							DiskSizeGB:      100,
						},
					},
					AzureFiles: []AzureFile{},
				},
				NetworkInfos: []VMNetworkInfo{
					{
						SubnetName:                         networks[0].Subnets[0].Name,
						LoadBalancerName:                   loadBalancers[0].Name,
						LoadBalancerBackendAddressPoolName: loadBalancers[0].BackendAddressPools[0].Name,
					},
				},
				Roles: []Role{},
			},
		}

		storageAccounts := []StorageAccount{
			{
				Name:     "storage_account_1",
				SKU:      "standard",
				Location: "eastus",
			},
		}

		mysqlDatabases := []MysqlDatabase{
			{
				Name: "kunlundb",
				MigrationInformation: &MigrationInformation{
					OriginHost:     "asd",
					OriginDatabase: "asd",
					OriginUsername: "asd",
					OriginPassword: "asd",
				},
				Version:             "5.7",
				Cores:               2,
				Tier:                "GeneralPurpose",
				Family:              "Gen5",
				Storage:             5,
				BackupRetentionDays: 35,
				SSLEnforcement:      "Enabled",
				Username:            "dbuser",
				Password:            "abcd1234!",
			},
		}

		// The checker add needed resource to manifest
		m = &Manifest{
			Schema:                "v0.1",
			EnvName:               "kunlun",
			IaaS:                  "azure",
			Location:              "eastus",
			ResourceGroupName:     "mykunlun",
			Platform:              &platform,
			VMGroups:              vmGroups,
			VNets:                 networks,
			LoadBalancers:         loadBalancers,
			StorageAccounts:       storageAccounts,
			NetworkSecurityGroups: networkSecurityGroups,
			MysqlDatabases:        mysqlDatabases,
		}

	})
	Describe("ToYAML", func() {
		Context("Everything OK", func() {
			It("should can be deserialize correctly", func() {
				b, err := m.ToYAML()
				Expect(err).To(BeNil())
				mCopy, err := NewManifestFromYAML(b)
				Expect(err).To(BeNil())
				deep_equal := reflect.DeepEqual(m, mCopy)
				if !deep_equal {
					if diff := deep.Equal(m, mCopy); diff != nil {
						diff_bytes, _ := json.Marshal(diff)
						println(string(diff_bytes))
					}
				}
				Expect(deep_equal).To(BeTrue())
				Expect(mCopy.VMGroups[0].Meta[0].Key == "jumpbox")
				Expect(mCopy.VMGroups[0].Meta[0].Value == "true")
			})
		})
	})
})
