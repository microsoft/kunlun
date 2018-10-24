package common

import (
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
)

type NonInfra struct {
	ProjectSourceCodePath string
	ProgrammingLanguage   ProgrammingLanguage
	Databases             []Database
}

type Infra struct {
	Size InfraSize
}

type Misc struct {
	ResourceGroupName    string `question:"What's the Azure resource group name for this Kunlun deployment?" default:"kl-test"`
	Location             string `question:"What's the Azure location for this Kunlun deployment?" default:"eastus"`
	AdminName            string `question:"What's the admin name for the jumpbox?" default:"kluser"`
	ConcurrentUserNumber int    `question:"What's your expected number of concurrent users?" default:"1000"`
}

type Blueprint struct {
	NonInfra NonInfra
	Infra    Infra
	Misc     Misc
}

type blueprintForYaml struct {
	ProjectSourceCodePath  string `yaml:"project_source_code_path,omitempty"`
	ProjectGitRevision     string `yaml:"project_git_revision,omitempty"`
	ProgrammingLanguage    string `yaml:"programming_language,omitempty"`
	DatabaseDriver         string `yaml:"database_driver,omitempty"`
	DatabaseVersion        string `yaml:"database_version,omitempty"`
	DatabaseStorage        int    `yaml:"database_storage,omitempty"`
	DatabaseOriginHost     string `yaml:"database_origin_host,omitempty"`
	DatabaseOriginName     string `yaml:"database_origin_name,omitempty"`
	DatabaseOriginUsername string `yaml:"database_origin_username,omitempty"`
	DatabaseOriginPassword string `yaml:"database_origin_password,omitempty"`
	VMGroupSize            string `yaml:"vmgroup_size,omitempty"`
	ResourceGroupName      string `yaml:"resource_group_name,omitempty"`
	Location               string `yaml:"location,omitempty"`
	AdminName              string `yaml:"admin_name,omitempty"`
	ConcurrentUserNumber   int    `yaml:"concurrent_user_number,omitempty"`
}

// TODO check if it fits into one of the artifacts templates
func (b Blueprint) finalValidate() error {
	return nil
}

func (b Blueprint) ExposeYaml(filePath string) error {
	if err := b.finalValidate(); err != nil {
		return err
	}
	bpfy := blueprintForYaml{
		ProjectSourceCodePath: b.NonInfra.ProjectSourceCodePath,
		ProjectGitRevision:    "master", // hard code the master, TODO @zhongyi, please help add it into the question later.
		ProgrammingLanguage:   string(b.NonInfra.ProgrammingLanguage),
		VMGroupSize:           string(b.Infra.Size),
		ResourceGroupName:     b.Misc.ResourceGroupName,
		Location:              b.Misc.Location,
		AdminName:             b.Misc.AdminName,
		ConcurrentUserNumber:  b.Misc.ConcurrentUserNumber,
	}
	// TODO support more. Assume at most one database for now.
	if len(b.NonInfra.Databases) > 0 {
		bpfy.DatabaseDriver = b.NonInfra.Databases[0].Driver
		bpfy.DatabaseVersion = b.NonInfra.Databases[0].Version
		bpfy.DatabaseStorage = b.NonInfra.Databases[0].Storage
		bpfy.DatabaseOriginHost = b.NonInfra.Databases[0].OriginHost
		bpfy.DatabaseOriginName = b.NonInfra.Databases[0].OriginName
		bpfy.DatabaseOriginUsername = b.NonInfra.Databases[0].OriginUsername
		bpfy.DatabaseOriginPassword = b.NonInfra.Databases[0].OriginPassword
	}
	bpBytes, _ := yaml.Marshal(bpfy)
	return ioutil.WriteFile(filePath, bpBytes, 0644)
}

func ImportBlueprintYaml(filePath string) (Blueprint, error) {
	bp := Blueprint{}
	bpfy := blueprintForYaml{}
	bpBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return bp, err
	}
	if err = yaml.Unmarshal(bpBytes, &bpfy); err != nil {
		return bp, err
	}

	bp.Infra = Infra{}
	bp.Infra.Size, err = ParseInfraSize(bpfy.VMGroupSize)
	if err != nil {
		return bp, err
	}

	bp.NonInfra = NonInfra{
		ProjectSourceCodePath: bpfy.ProjectSourceCodePath,
	}

	bp.NonInfra.ProgrammingLanguage, err =
		ParseProgrammingLanguage(bpfy.ProgrammingLanguage)
	if err != nil {
		return bp, err
	}

	// TODO support more. Assume at most one database for now.
	db := Database{
		Driver:         bpfy.DatabaseDriver,
		Version:        bpfy.DatabaseVersion,
		Storage:        bpfy.DatabaseStorage,
		OriginHost:     bpfy.DatabaseOriginHost,
		OriginName:     bpfy.DatabaseOriginName,
		OriginUsername: bpfy.DatabaseOriginUsername,
		OriginPassword: bpfy.DatabaseOriginPassword,
	}
	allEmpty := true
	s := reflect.ValueOf(&db).Elem()
	for i := 0; i < s.NumField(); i++ {
		valField := s.Field(i)
		val := valField.Interface()
		if val != reflect.Zero(valField.Type()).Interface() {
			allEmpty = false
		}
	}
	if !allEmpty {
		bp.NonInfra.Databases = append(bp.NonInfra.Databases, db)
	}

	bp.Misc = Misc{
		ResourceGroupName:    bpfy.ResourceGroupName,
		Location:             bpfy.Location,
		AdminName:            bpfy.AdminName,
		ConcurrentUserNumber: bpfy.ConcurrentUserNumber,
	}

	if err = bp.finalValidate(); err != nil {
		return bp, err
	}

	return bp, nil
}
