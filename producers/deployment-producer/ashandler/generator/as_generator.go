package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Microsoft/kunlun/artifacts/deployments"
	builtinroles "github.com/Microsoft/kunlun/built-in-roles"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
	patching "github.com/Microsoft/kunlun/patching"
	yaml "gopkg.in/yaml.v2"
)

type logger interface {
	Step(string, ...interface{})
	Printf(string, ...interface{})
	Println(string)
	Prompt(string) bool
}

type ASGenerator struct {
	stateStore storage.Store
	logger     logger
	fs         fileio.Fs
}

func NewASGenerator(
	stateStore storage.Store,
	logger logger,
	fs fileio.Fs,
) ASGenerator {
	return ASGenerator{
		stateStore: stateStore,
		logger:     logger,
		fs:         fs,
	}
}

// https://docs.ansible.com/ansible/latest/user_guide/playbooks_reuse_roles.html?highlight=roles
func (a ASGenerator) Generate(hostGroups []deployments.HostGroup, deployments []deployments.Deployment) error {
	// generate the ansible config file.
	builtInRolesFS, err := builtinroles.FSByte(false, "/ansible.cfg")
	if err != nil {
		return err
	}
	ansibleDir, err := a.stateStore.GetAnsibleDir()
	ansibleConfigFile := path.Join(ansibleDir, "ansible.cfg")
	a.fs.WriteFile(ansibleConfigFile, builtInRolesFS, 0644)

	// generate the hosts files.
	hostsFileContent, err := a.generateHostsFile(hostGroups)
	if err != nil {
		return err
	}
	ansibleInventoriesDir, _ := a.stateStore.GetAnsibleInventoriesDir()
	hostsFile := path.Join(ansibleInventoriesDir, "hosts.yml")
	a.logger.Printf("writing hosts file to %s\n", hostsFile)
	err = a.fs.WriteFile(hostsFile, hostsFileContent, 0644)
	if err != nil {
		a.logger.Printf("write file failed: %s\n", err.Error())
		return err
	}

	err = a.prepareBuiltInRoles(deployments)
	if err != nil {
		a.logger.Printf("prepare built in roles failed: %s\n", err.Error())
		return err
	}
	// generate the roles files.
	playbookContent := a.generatePlaybookFile(deployments)
	ansibleMainFile, err := a.stateStore.GetAnsibleMainFile()

	a.logger.Printf("writing playbook file to %s\n", ansibleMainFile)
	err = ioutil.WriteFile(ansibleMainFile, playbookContent, 0644)
	if err != nil {
		a.logger.Printf("write file failed: %s\n", err.Error())
		return err
	}

	// generate the private key.
	privateSshKey, err := a.getAdminSSHPrivateKey()
	if err != nil {
		a.logger.Printf("get admin ssh private key failed: %s\n", err.Error())
		return err
	}
	sshPrivateKeyPath, err := a.getSSHPrivateKeyPath()
	if err != nil {
		return err
	}
	err = a.fs.WriteFile(sshPrivateKeyPath, ([]byte)(privateSshKey), 0600)
	if err != nil {
		return err
	}

	// generate the deployment script file.
	deploymentScriptFilePath, err := a.stateStore.GetDeploymentScriptFile()
	deploymentScriptContent, err := a.generateDeploymentScript()
	if err != nil {
		return err
	}
	err = a.fs.WriteFile(deploymentScriptFilePath, deploymentScriptContent, 0744)
	if err != nil {
		a.logger.Printf("write file failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (a ASGenerator) getSSHPrivateKeyPath() (string, error) {
	varsFolder, err := a.stateStore.GetVarsDir()
	if err != nil {
		return "", err
	}
	sshPrivateKeyPath := path.Join(varsFolder, "admin_ssh_private_key")
	return sshPrivateKeyPath, nil
}

func (a ASGenerator) getAdminSSHPrivateKey() (string, error) {
	varsStore := patching.NewVarsFSStore(
		a.fs,
	)

	varsStoreFilePath, err := a.stateStore.GetMainArtifactVarsStoreFilePath()
	if err != nil {
		return "", err
	}

	err = varsStore.UnmarshalFlag(varsStoreFilePath)
	if err != nil {
		return "", err
	}
	privateKey := patching.VariableDefinition{
		Name: "admin_ssh", // TODO(andy) HARD CODE THIS, think about a better way to get the private key.
	}

	adminSSH, _, err := varsStore.Get(privateKey)

	if err != nil {
		return "", err
	}
	for k, v := range adminSSH.(map[interface{}]interface{}) {
		if k.(string) == "private_key" {
			return v.(string), nil
		}
	}
	return "", errors.New("no privat key found")
}

// TODO error handling.
func (a ASGenerator) generateHostsFile(hostGroups []deployments.HostGroup) ([]byte, error) {
	// ---
	// sample_server:
	// 	 hosts:
	// 	   172.16.8.4:
	// 	     ansible_ssh_user: andy
	// 	     ansible_ssh_common_args: '-o ProxyCommand="ssh -W %h:%p -q andy@65.52.176.243" -i private key'
	hostGroupsSlices := yaml.MapSlice{}
	privateKeyPath, err := a.getSSHPrivateKeyPath()
	if err != nil {
		return nil, err
	}
	for _, hostGroup := range hostGroups {
		hosts := yaml.MapSlice{}

		for _, host := range hostGroup.Hosts {
			sshCommonArgs := ""
			if host.SSHCommonArgs != "" {
				sshCommonArgs = host.SSHCommonArgs[0:len(host.SSHCommonArgs)-1] + " -i " + privateKeyPath + "\""
			}
			hostSlice := yaml.MapItem{
				Key: host.Alias,
				Value: AnsibleHost{
					Host:          host.Host,
					SSHUser:       host.User,
					SSHCommonArgs: sshCommonArgs,
				},
			}
			hosts = append(hosts, hostSlice)
		}

		hostGroupSlice := yaml.MapItem{
			Key: hostGroup.Name,
			Value: yaml.MapSlice{
				{
					Key:   "hosts",
					Value: hosts,
				},
			},
		}
		hostGroupsSlices = append(hostGroupsSlices, hostGroupSlice)
	}
	content, _ := yaml.Marshal(hostGroupsSlices)
	return content, nil
}

type AnsibleHost struct {
	Host          string `yaml:"ansible_host"`
	SSHUser       string `yaml:"ansible_ssh_user"`
	SSHCommonArgs string `yaml:"ansible_ssh_common_args,omitempty"`
}

type role struct {
	Role       string `yaml:"role"`
	Become     string `yaml:"become,omitempty"`
	BecomeUser string `yaml:"become_user,omitempty"`
}
type depItem struct {
	Hosts    string   `yaml:"hosts"`
	VarsFile []string `yaml:"vars_files"`
	Roles    []role   `yaml:"roles"`
}

// TODO error handling.
func (a ASGenerator) generatePlaybookFile(deployments []deployments.Deployment) []byte {
	// ---
	// - hosts: sample_server
	//   vars_files:
	// 	   - vars/sample.yml
	//   roles:
	// 	   - role: 'geerlingguy.composer'
	// 	     become: true
	// 	     become_user: root
	// - hosts: sample_server2
	//   vars_files:
	// 	   - vars/sample.yml
	//   roles:
	// 	   - role: 'geerlingguy.php'
	// 	     become: true
	// 	     become_user: root

	depItems := []depItem{}
	// write the vars files
	for _, dep := range deployments {
		// write the files
		varsDir, _ := a.stateStore.GetAnsibleDir()
		varsFile := path.Join(varsDir, dep.HostGroupName+".yml")
		varsContent, _ := yaml.Marshal(dep.Vars)

		a.logger.Printf("writing vars file to %s\n", varsFile)
		err := ioutil.WriteFile(varsFile, varsContent, 0644)
		if err != nil {
			a.logger.Printf("write vars file failed: %s\n", err.Error())
		}
		depItem := depItem{
			Hosts:    dep.HostGroupName,
			VarsFile: []string{varsFile},
		}

		roles := []role{}
		for _, r := range dep.Roles {
			if r.BecomeUser != "" {
				roles = append(roles, role{
					Role:       r.Name,
					Become:     "true",
					BecomeUser: r.BecomeUser,
				})
			} else {
				roles = append(roles, role{
					Role: r.Name,
				})
			}
		}

		depItem.Roles = roles
		depItems = append(depItems, depItem)
	}
	content, _ := yaml.Marshal(depItems)
	return content
}

func (a ASGenerator) prepareBuiltInRoles(deployments []deployments.Deployment) error {
	//
	ansibleDir, err := a.stateStore.GetAnsibleDir()
	if err != nil {
		return err
	}
	stack := NewStack()
	stack.Push("/built.in")
	stack.Push("/roles.galaxy")

	fs := builtinroles.FS(false)

	for stack.Len() > 0 {
		currentItem := stack.Pop()
		currentPath := currentItem.(string)
		file, err := fs.Open(currentPath)
		if err != nil {
			return err
		}
		files, err := file.Readdir(0)
		for _, file := range files {
			filePath := path.Join(currentPath, file.Name())
			if file.IsDir() {
				stack.Push(filePath)
			} else {
				content, err := builtinroles.FSByte(false, filePath)
				if err != nil {
					return err
				}
				targetPath := path.Join(ansibleDir, "./"+filePath)
				// a.logger.Printf("base dir is %s\n", path.Dir(targetPath))
				err = a.fs.MkdirAll(path.Dir(targetPath), 0744)
				if err != nil {
					return err
				}
				// a.logger.Printf("writing to %s\n", targetPath)
				err = a.fs.WriteFile(targetPath, content, 0644)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a ASGenerator) generateDeploymentScript() ([]byte, error) {
	// varsDir, _ := a.stateStore.GetAnsibleDir()
	sshPrivateKeyPath, err := a.getSSHPrivateKeyPath()
	if err != nil {
		return nil, err
	}
	deploymentScript := fmt.Sprintf(`#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export ANSIBLE_CONFIG=$DIR/ansible/ansible.cfg
ansible-playbook -i $DIR/ansible/inventories $DIR/ansible/main.yml -vv --private-key=%s
`, sshPrivateKeyPath)
	return []byte(deploymentScript), nil
}
