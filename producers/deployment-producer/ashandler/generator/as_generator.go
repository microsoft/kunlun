package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/artifacts/builtinroles"
	"github.com/Microsoft/kunlun/artifacts/deployments"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	patching "github.com/Microsoft/kunlun/patching"
	yaml "gopkg.in/yaml.v2"
)

type ASGenerator struct {
	stateStore storage.Store
	ui         *ui.UI
	fs         fileio.Fs
}

func NewASGenerator(
	stateStore storage.Store,
	ui *ui.UI,
	fs fileio.Fs,
) ASGenerator {
	return ASGenerator{
		stateStore: stateStore,
		ui:         ui,
		fs:         fs,
	}
}

type SSHConfig struct {
	UserKnownHostsFile    string
	SSHPrivateKey         string
	JumpboxUser           string
	JumpboxHost           string
	StrictHostKeyChecking string
}

// https://docs.ansible.com/ansible/latest/user_guide/playbooks_reuse_roles.html?highlight=roles
func (a ASGenerator) Generate(hostGroups []deployments.HostGroup, deployments []deployments.Deployment) error {
	// generate the ansible config file.
	filesToCopy := []string{"ansible.cfg"}
	for _, f := range filesToCopy {
		builtInRolesFS, err := builtinroles.FSByte(false, fmt.Sprintf("/%s", f))
		if err != nil {
			return err
		}
		ansibleDir, err := a.stateStore.GetAnsibleDir()
		filePath := path.Join(ansibleDir, f)
		a.fs.WriteFile(filePath, builtInRolesFS, 0644)
	}

	// generate the hosts files.
	hostsFileContent, err := a.generateHostsFile(hostGroups)
	if err != nil {
		return err
	}
	ansibleInventoriesDir, _ := a.stateStore.GetAnsibleInventoriesDir()
	hostsFile := path.Join(ansibleInventoriesDir, "hosts.yml")
	a.ui.Printf("writing hosts file to %s\n", hostsFile)
	err = a.fs.WriteFile(hostsFile, hostsFileContent, 0644)
	if err != nil {
		a.ui.Printf("write file failed: %s\n", err.Error())
		return err
	}

	err = a.prepareBuiltInRoles(deployments)
	if err != nil {
		a.ui.Printf("prepare built in roles failed: %s\n", err.Error())
		return err
	}
	// generate the roles files.
	playbookContent := a.generatePlaybookFile(deployments)
	ansibleMainFile, err := a.stateStore.GetAnsibleMainFile()

	a.ui.Printf("writing playbook file to %s\n", ansibleMainFile)
	err = ioutil.WriteFile(ansibleMainFile, playbookContent, 0644)
	if err != nil {
		a.ui.Printf("write file failed: %s\n", err.Error())
		return err
	}

	// generate the private key.
	privateSshKey, err := a.getAdminSSHPrivateKey()
	if err != nil {
		a.ui.Printf("get admin ssh private key failed: %s\n", err.Error())
		return err
	}
	sshPrivateKeyPath, err := a.stateStore.GetAdminSSHPrivateKeyPath()
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
		a.ui.Printf("write file failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (a ASGenerator) getSSHKnownHostsFilePath() (string, error) {
	deploymentFolder, err := a.stateStore.GetDeploymentsDir()
	if err != nil {
		return "", err
	}
	sshKnownHostsFilePath := path.Join(deploymentFolder, "known_hosts")
	return sshKnownHostsFilePath, nil
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

	err := a.provisionJumpboxParameters(hostGroups)
	if err != nil {
		return nil, err
	}
	for _, hostGroup := range hostGroups {
		hosts := yaml.MapSlice{}

		for _, host := range hostGroup.Hosts {

			hostSlice := yaml.MapItem{
				Key: host.Alias,
				Value: AnsibleHost{
					Host:          host.Host,
					SSHUser:       host.User,
					SSHCommonArgs: host.SSHCommonArgs,
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

func (a ASGenerator) provisionJumpboxParameters(hostGroups []deployments.HostGroup) error {
	// find the ip of the jumpbox
	var (
		jumpboxUser string
		jumpboxHost string
	)
	for _, hostGroup := range hostGroups {
		if hostGroup.GroupType == artifacts.JumpboxHostGroupType {
			jumpboxUser = hostGroup.Hosts[0].User
			jumpboxHost = hostGroup.Hosts[0].Host
		}
	}
	privateKeyPath, err := a.stateStore.GetAdminSSHPrivateKeyPath()
	if err != nil {
		return err
	}
	knownHostsFilePath, err := a.getSSHKnownHostsFilePath()
	if err != nil {
		return err
	}
	jumpBoxSSHCommonArgs := "-o UserKnownHostsFile={{.UserKnownHostsFile}} -o StrictHostKeyChecking={{.StrictHostKeyChecking}}"
	nonJumpboxSSHCommonArgs := "-o UserKnownHostsFile={{.UserKnownHostsFile}} -o StrictHostKeyChecking={{.StrictHostKeyChecking}} " +
		"-o ProxyCommand=\"ssh -o StrictHostKeyChecking={{.StrictHostKeyChecking}} -W %h:%p -q {{.JumpboxUser}}@{{.JumpboxHost}} -i {{.SSHPrivateKey}}\""
	sshConfig := SSHConfig{
		UserKnownHostsFile:    knownHostsFilePath,
		SSHPrivateKey:         privateKeyPath,
		JumpboxUser:           jumpboxUser,
		JumpboxHost:           jumpboxHost,
		StrictHostKeyChecking: "{{ansible_host_key_checking}}",
	}
	jumpBoxSSHCommonArgs, err = a.provisionSSHCommonArgs(jumpBoxSSHCommonArgs, sshConfig)
	if err != nil {
		return err
	}
	nonJumpboxSSHCommonArgs, err = a.provisionSSHCommonArgs(nonJumpboxSSHCommonArgs, sshConfig)
	if err != nil {
		return err
	}
	for _, hostGroup := range hostGroups {
		for index := range hostGroup.Hosts {
			if hostGroup.GroupType != artifacts.JumpboxHostGroupType {
				hostGroup.Hosts[index].SSHCommonArgs = nonJumpboxSSHCommonArgs
			} else {
				hostGroup.Hosts[index].SSHCommonArgs = jumpBoxSSHCommonArgs
			}
		}
	}
	return nil
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

func (a ASGenerator) provisionSSHCommonArgs(commonArgs string, sshConfig SSHConfig) (string, error) {
	var tpl bytes.Buffer
	tmpl, err := template.New("ssh_args").Parse(commonArgs)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&tpl, sshConfig)
	if err != nil {
		return "", err
	}
	sshCommonArgs := tpl.String()
	return sshCommonArgs, nil
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
		mainVarsFilePath, _ := a.stateStore.GetMainArtifactVarsFilePath()
		varsDir, _ := a.stateStore.GetAnsibleDir()
		varsFile := path.Join(varsDir, dep.HostGroupName+".yml")
		varsContent, _ := yaml.Marshal(dep.Vars)

		a.ui.Printf("writing vars file to %s\n", varsFile)
		err := ioutil.WriteFile(varsFile, varsContent, 0644)
		if err != nil {
			a.ui.Printf("write vars file failed: %s\n", err.Error())
		}
		depItem := depItem{
			Hosts:    dep.HostGroupName,
			VarsFile: []string{mainVarsFilePath, varsFile},
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
				// a.ui.Printf("base dir is %s\n", path.Dir(targetPath))
				err = a.fs.MkdirAll(path.Dir(targetPath), 0744)
				if err != nil {
					return err
				}
				// a.ui.Printf("writing to %s\n", targetPath)
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
	sshPrivateKeyPath, err := a.stateStore.GetAdminSSHPrivateKeyPath()
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
