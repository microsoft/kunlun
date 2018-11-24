package handler

type executor interface {
	Version() (string, error)
	Setup(terraformTemplate string, inputs map[string]interface{}) error
	Init() error
	Apply(credentials map[string]string) error
	Validate(credentials map[string]string) error
	Destroy(credentials map[string]string) error
	Outputs() (map[string]interface{}, error)
	Output(string) (string, error)
	IsPaved() (bool, error)
}
