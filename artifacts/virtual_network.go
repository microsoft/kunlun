package apis

type VirtualNetwork struct {
	Name         string   `yaml:"name"`
	AddressSpace string   `yaml:"address_space"`
	Subnets      []Subnet `yaml:"subnets"`
}

func virtualNetworkValidator(m Manifest) error {
	for _, vn := range m.VNets {
		occuredNames := make(map[string]bool)
		if vn.Name == "" {
			return validationError("vnet name can't be empty")
		}
		if occuredNames[vn.Name] {
			return validationError("vnet name must be unique, but %s occurs at least twice", vn.Name)
		}
		occuredNames[vn.Name] = true
	}
	return nil
}
