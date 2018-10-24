package apis

type OSDisk struct {
	ManagedDiskType string `yaml:"managed_disk_type"`
	Caching         string `yaml:"caching"`
	CreateOption    string `yaml:"create_option"`
}
