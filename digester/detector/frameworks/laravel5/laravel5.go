package laravel5

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/Microsoft/kunlun/digester/common"
	"github.com/Microsoft/kunlun/digester/detector/util"
)

type framework struct{}

func New() common.Framework {
	return &framework{}
}

func (f *framework) GetName() common.FrameworkName {
	return "Laravel5"
}

func (f *framework) GetProgrammingLanguage() common.ProgrammingLanguage {
	return common.PHP
}

func (f *framework) DetectConfig(path string) []common.Database {
	res := []common.Database{}
	envPath := fmt.Sprintf(
		"%s%s.env",
		path,
		string(os.PathSeparator),
	)
	envBytes, err := ioutil.ReadFile(envPath)
	if err != nil {
		fmt.Print(err)
	}
	envStr := string(envBytes)
	var rex = regexp.MustCompile("([a-zA-Z_][a-zA-Z0-9_]*)=([^\r]*)")
	envMatch := rex.FindAllStringSubmatch(envStr, -1)
	envMap := make(map[string]string)
	for _, kv := range envMatch {
		envMap[kv[1]] = kv[2]
	}

	configPath := fmt.Sprintf(
		"%s%sconfig%sdatabase.php",
		path,
		string(os.PathSeparator),
		string(os.PathSeparator),
	)
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Print(err)
	}
	configStr := string(configBytes)
	var l int
	l = -1
	for i, v := range configStr {
		if v == '[' {
			l = i
		}
		if v == ']' && l != -1 {
			block := configStr[l : i+1]
			l = -1
			driver := util.KeyValueParser1(block, "driver", false)
			if driver == "mysql" {
				host := util.KeyValueParser1(block, "host", true)
				database := util.KeyValueParser1(block, "database", true)
				username := util.KeyValueParser1(block, "username", true)
				password := util.KeyValueParser1(block, "password", true)

				res = append(res, common.Database{
					Driver:         driver,
					OriginHost:     envMap[host],
					OriginName:     envMap[database],
					OriginUsername: envMap[username],
					OriginPassword: envMap[password],
					/*
					   EnvVarHost: host,
					   EnvVarDatabase: database,
					   EnvVarUsername: username,
					   EnvVarPassword: password,
					*/
				})
			}
		}
	}
	return res
}
