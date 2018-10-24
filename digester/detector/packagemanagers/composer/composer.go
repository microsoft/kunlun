package composer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Microsoft/kunlun/digester/common"
	"github.com/Microsoft/kunlun/digester/detector/frameworks/laravel5"
)

type packageManager struct{}

type composerConfig struct {
	Name    string                 `json:"name"`
	Require map[string]interface{} `json:"require"`
}

func New() common.PackageManager {
	return &packageManager{}
}

func (p *packageManager) GetName() common.PackageManagerName {
	return "Composer"
}

func (p *packageManager) Identify(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.Name() == "composer.json" {
			return true
		}
	}
	return false
}

func (p *packageManager) DetectFramework(path string) []common.FrameworkName {
	composerFile, err := os.Open(path + "composer.json")
	if err != nil {
		log.Fatal(err)
	}
	defer composerFile.Close()
	composerByte, _ := ioutil.ReadAll(composerFile)
	var composerConfig composerConfig
	json.Unmarshal(composerByte, &composerConfig)

	possibleFrameworks := []common.FrameworkName{}

	version, ok := composerConfig.Require["laravel/framework"].(string)
	if ok {
		if strings.HasPrefix(version, "5") {
			possibleFrameworks = append(possibleFrameworks, laravel5.New().GetName())
		}
	}

	return possibleFrameworks
}
