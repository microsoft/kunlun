package apis

type AzureFile struct {
	StorageAccount string `yaml:"storage_account"`
	Name           string `yaml:"name"`
	MountPoint     string `yaml:"mount_point"`
}
