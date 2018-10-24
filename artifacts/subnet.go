package apis

type Subnet struct {
	Name    string `yaml:"name"`
	Range   string `yaml:"range"`
	Gateway string `yaml:"gateway"`
}
