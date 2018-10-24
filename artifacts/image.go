package apis

type Image struct {
	Offer     string `yaml:"offer"`
	Publisher string `yaml:"publisher"`
	SKU       string `yaml:"sku"`
	Version   string `yaml:"version"`
}
