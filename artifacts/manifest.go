package apis

import (
	yaml "gopkg.in/yaml.v2"
)

// Manifest contains all needed information, all later on modules will
// use this manifest
type Manifest struct {
	Schema                string                 `yaml:"schema,omitempty"`
	EnvName               string                 `yaml:"env_name"`
	ResourceGroupName     string                 `yaml:"resource_group_name"`
	Location              string                 `yaml:"location,omitempty"`
	IaaS                  string                 `yaml:"iaas,omitempty"`
	Platform              *Platform              `yaml:"platform,omitempty"`
	VMGroups              []VMGroup              `yaml:"vm_groups,omitempty"`
	VNets                 []VirtualNetwork       `yaml:"vnets,omitempty"`
	LoadBalancers         []LoadBalancer         `yaml:"load_balancers,omitempty"`
	StorageAccounts       []StorageAccount       `yaml:"storage_accounts,omitempty"`
	NetworkSecurityGroups []NetworkSecurityGroup `yaml:"network_security_groups,omitempty"`
	MysqlDatabases        []MysqlDatabase        `yaml:"mysql_databases,omitempty"`
}

func (m *Manifest) validate() error {
	return newValidator().Validate(*m)
}

// ToYAML converts the object to YAML bytes
func (m *Manifest) ToYAML() (b []byte, err error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	return yaml.Marshal(m)
}

// NewManifestFromYAML convert yaml bytes to Manifest object
func NewManifestFromYAML(b []byte) (m *Manifest, err error) {
	var manifest Manifest
	if err := yaml.Unmarshal(b, &manifest); err != nil {
		return nil, err
	}
	if err := manifest.validate(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m *Manifest) GetSubnetByName(subnetName string) *Subnet {
	for _, vnet := range m.VNets {
		for _, snet := range vnet.Subnets {
			if snet.Name == subnetName {
				return &snet
			}
		}
	}
	return nil
}

func (m *Manifest) GetLoadBalancerByName(loadBalancerName string) *LoadBalancer {
	for _, lb := range m.LoadBalancers {
		if lb.Name == loadBalancerName {
			return &lb
		}
	}
	return nil
}

func (m *Manifest) GetNetworkSecurityGroupByName(nsgName string) *NetworkSecurityGroup {
	for _, nsg := range m.NetworkSecurityGroups {
		if nsg.Name == nsgName {
			return &nsg
		}
	}
	return nil
}
