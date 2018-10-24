package apis

// LoadBalancer contains needed information to create a load balancer on Azure.
type LoadBalancer struct {
	Name                string                           `yaml:"name"`
	SKU                 string                           `yaml:"sku"`
	BackendAddressPools []LoadBalancerBackendAddressPool `yaml:"backend_address_pools"`
	HealthProbes        []LoadBalancerHealthProbe        `yaml:"health_probes"`
	Rules               []LoadBalancerRule               `yaml:"rules"`
}

type LoadBalancerBackendAddressPool struct {
	Name string `yaml:"name"`
}

type LoadBalancerHealthProbe struct {
	Name        string `yaml:"name"`
	Protocol    string `yaml:"protocol"`
	Port        int    `yaml:"port"`
	RequestPath string `yaml:"request_path"`
}

type LoadBalancerRule struct {
	Name                   string `yaml:"name"`
	Protocol               string `yaml:"protocol"`
	FrontendPort           int    `yaml:"frontend_port"`
	BackendPort            int    `yaml:"backend_port"`
	BackendAddressPoolName string `yaml:"backend_address_pool_name"`
	HealthProbeName        string `yaml:"health_probe_name"`
}

func loadBalancerValidator(m Manifest) error {
	for _, lb := range m.LoadBalancers {
		occuredNames := make(map[string]bool)
		if lb.Name == "" {
			return validationError("load balancer name can't be empty")
		}
		if occuredNames[lb.Name] {
			return validationError("load balancer name must be unique, but %s occurs at least twice", lb.Name)
		}
		occuredNames[lb.Name] = true
		if lb.SKU != "Basic" && lb.SKU != "Standard" {
			return validationError("sku for load balancer %s must be \"Standard\" or \"Basic\"", lb.Name)
		}
	}
	return nil
}

const (
	LoadBalancerStandardSKU = "standard"
	LoadBalancerBasicSKU    = "basic"
)
