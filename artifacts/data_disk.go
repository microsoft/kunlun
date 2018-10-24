package apis

type DataDisk struct {
	ManagedDiskType string `yaml:"managed_disk_type"`
	Caching         string `yaml:"caching"`
	CreateOption    string `yaml:"create_option"`
	DiskSizeGB      int    `yaml:"disk_size_gb"`
}
