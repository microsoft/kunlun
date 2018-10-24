package apis

type StorageAccount struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	SKU      string `yaml:"sku"`
}
