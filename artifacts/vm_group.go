package apis

import (
	"gopkg.in/yaml.v2"
)

// VMGroup contains needed information to create a set of VMs on Azure. VMs in the group
// will have the same SKU, using the same subnet.
type VMGroup struct {
	Name         string          `yaml:"name"`
	Meta         yaml.MapSlice   `yaml:"meta,omitempty"`
	Count        int             `yaml:"count"`
	SKU          string          `yaml:"sku"`
	Type         string          `yaml:"type"`
	OSProfile    VMOSProfile     `yaml:"os_profile"`
	Storage      *VMStorage      `yaml:"storage"`
	NetworkInfos []VMNetworkInfo `yaml:"networks"`
	Roles        []Role          `yaml:"roles"`
}

type VMOSProfile struct {
	AdminName          string             `yaml:"admin_name"`
	LinuxConfiguration LinuxConfiguration `yaml:"linux_configuration"`
}

type LinuxConfiguration struct {
	SSH SSH `yaml:"ssh"`
}

type SSH struct {
	PublicKeys []string `yaml:"public_keys,omitempty"`
}

type VMStorage struct {
	Image      *Image      `yaml:"image"`
	OSDisk     *OSDisk     `yaml:"os_disk"`
	DataDisks  []DataDisk  `yaml:"data_disks"`
	AzureFiles []AzureFile `yaml:"azure_files"`
}

type VMNetworkInfo struct {
	SubnetName                         string            `yaml:"subnet_name"`
	LoadBalancerName                   string            `yaml:"load_balancer_name"`
	LoadBalancerBackendAddressPoolName string            `yaml:"load_balancer_backend_address_pool_name"`
	NetworkSecurityGroupName           string            `yaml:"network_security_group_name"`
	PublicIP                           string            `yaml:"public_ip"`
	Outputs                            []VMNetworkOutput `yaml:"outputs,omitempty"`
}

type VMNetworkOutput struct {
	IP       string `yaml:"ip"`
	PublicIP string `yaml:"public_ip"`
	Host     string `yaml:"host"`
}

func vmGroupValidator(m Manifest) error {
	occuredNames := make(map[string]bool)
	for _, vmGroup := range m.VMGroups {
		if vmGroup.Name == "" {
			return validationError("vm group name can't be empty")
		}
		if occuredNames[vmGroup.Name] {
			return validationError("vm group name must be unique, but %s occurs at least twice", vmGroup.Name)
		}
		occuredNames[vmGroup.Name] = true
		if vmGroup.Count <= 0 {
			return validationError("count in vm group %s must be greater than 0", vmGroup.Name)
		}
		if vmGroup.Type != "vm" {
			return validationError("type in vm group %s is %s, but only vm is supported for now", vmGroup.Type, vmGroup.Name)
		}
	}
	return nil
}

const (
	VMStandardB1s  = "Standard_B1s"
	VMStandardB1ms = "Standard_B1ms"
	VMStandardB2s  = "Standard_B2s"
	VMStandardB2ms = "Standard_B2ms"
	VMStandardB4ms = "Standard_B4ms"
	VMStandardB8ms = "Standard_B8ms"

	VMStandardD2sV3  = "Standard_D2s_v3"
	VMStandardD4sV3  = "Standard_D4s_v3"
	VMStandardD8sV3  = "Standard_D8s_v3"
	VMStandardD16sV3 = "Standard_D16s_v3"
	VMStandardD32sV3 = "Standard_D32s_v3"
	VMStandardD64sV3 = "Standard_D64s_v3"

	VMStandardD2V3  = "Standard_D2_v3"
	VMStandardD4V3  = "Standard_D4_v3"
	VMStandardD8V3  = "Standard_D8_v3"
	VMStandardD16V3 = "Standard_D16_v3"
	VMStandardD32V3 = "Standard_D32_v3"
	VMStandardD64V3 = "Standard_D64_v3"

	VMStandardDS1V2 = "Standard_DS1_v2"
	VMStandardDS2V2 = "Standard_DS2_v2"
	VMStandardDS3V2 = "Standard_DS3_v2"
	VMStandardDS4V2 = "Standard_DS4_v2"
	VMStandardDS5V2 = "Standard_DS5_v2"

	VMStandardD1V2 = "Standard_D1_v2"
	VMStandardD2V2 = "Standard_D2_v2"
	VMStandardD3V2 = "Standard_D3_v2"
	VMStandardD4V2 = "Standard_D4_v2"
	VMStandardD5V2 = "Standard_D5_v2"

	VMStandardA1V2  = "Standard_A1_v2"
	VMStandardA2V2  = "Standard_A2_v2"
	VMStandardA4V2  = "Standard_A4_v2"
	VMStandardA8V2  = "Standard_A8_v2"
	VMStandardA2mV2 = "Standard_A2m_v2"
	VMStandardA4mV2 = "Standard_A4m_v2"
	VMStandardA8mV2 = "Standard_A8m_v2"

	VMStandardDC2s = "Standard_DC2s"
	VMStandardDC4s = "Standard_DC4s"

	VMStandardF2sV2  = "Standard_F2s_v2"
	VMStandardF4sV2  = "Standard_F4s_v2"
	VMStandardF8sV2  = "Standard_F8s_v2"
	VMStandardF16sV2 = "Standard_F16s_v2"
	VMStandardF32sV2 = "Standard_F32s_v2"
	VMStandardF64sV2 = "Standard_F64s_v2"
	VMStandardF72sV2 = "Standard_F72s_v2"

	VMStandardF1s  = "Standard_F1s"
	VMStandardF2s  = "Standard_F2s"
	VMStandardF4s  = "Standard_F4s"
	VMStandardF8s  = "Standard_F8s"
	VMStandardF16s = "Standard_F16s"

	VMStandardF1  = "Standard_F1"
	VMStandardF2  = "Standard_F2"
	VMStandardF4  = "Standard_F4"
	VMStandardF8  = "Standard_F8"
	VMStandardF16 = "Standard_F16"
)
