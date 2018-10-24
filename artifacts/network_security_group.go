package apis

type NetworkSecurityGroup struct {
	Name                 string                `yaml:"name"`
	NetworkSecurityRules []NetworkSecurityRule `yaml:"network_security_rules"`
}

type NetworkSecurityRule struct {
	Name                     string `yanl:"name"`
	Priority                 int    `yaml:"priority"`
	Direction                string `yaml:"direction"`
	Access                   string `yaml:"access"`
	Protocol                 string `yaml:"protocol"`
	SourcePortRange          string `yaml:"source_port_range"`
	DestinationPortRange     string `yaml:"destination_port_range"`
	SourceAddressPrefix      string `yaml:"source_address_prefix"`
	DestinationAddressPrefix string `yaml:"destination_address_prefix"`
}
