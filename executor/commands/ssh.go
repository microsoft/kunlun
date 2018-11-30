package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	patching "github.com/Microsoft/kunlun/patching"
	"github.com/Microsoft/kunlun/ssh"

	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/flags"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
)

type SSH struct {
	stateStore storage.Store
	fs         fileio.Fs
	ui         *ui.UI
	cli        sshCLI
	randomPort randomPort
}

type sshCLI interface {
	Run([]string) error
	Start([]string) (*exec.Cmd, error)
}

type randomPort interface {
	GetPort() (string, error)
}

func NewSSH(
	stateStore storage.Store,
	fs fileio.Fs,
	ui *ui.UI,
) SSH {
	sshCLI := ssh.NewCLI(os.Stdin, os.Stdout, os.Stderr)
	randomPort := ssh.RandomPort{}
	return SSH{
		stateStore: stateStore,
		fs:         fs,
		ui:         ui,
		cli:        sshCLI,
		randomPort: randomPort,
	}
}

func (s SSH) CheckFastFails(args []string, state storage.State) error {
	return nil
}

func (s SSH) Execute(args []string, state storage.State) error {
	var (
		vmGroupName string
		indexStr    string
		indexValue  int
	)
	sshFlags := flags.New("ssh")
	sshFlags.String(&vmGroupName, "group", "")
	sshFlags.String(&indexStr, "index", "0")
	err := sshFlags.Parse(args)
	if err != nil {
		return err
	}

	indexValue, err = strconv.Atoi(indexStr)
	if err != nil {
		return err
	}
	patching := patching.NewPatching(s.stateStore, s.fs)
	manifest, err := patching.ProvisionManifest()
	if err != nil {
		return err
	}

	jumpBoxUserName := ""
	jumpBoxIP := ""
	targetVMUserName := ""
	targetVMIP := ""

	// TODO this will not work for the plain_vm scenario, should handle that.
	for _, vmGroup := range manifest.VMGroups {
		if vmGroup.Jumpbox() {
			jumpBoxIP = vmGroup.NetworkInfos[0].Outputs[0].PublicIP
			jumpBoxUserName = vmGroup.OSProfile.AdminName
		}
		if vmGroup.Name == vmGroupName {
			targetVMIP = vmGroup.NetworkInfos[0].Outputs[indexValue].IP
			targetVMUserName = vmGroup.OSProfile.AdminName
		}
	}
	adminSSHPrivateKey, err := s.stateStore.GetAdminSSHPrivateKeyPath()
	if err != nil {
		return err
	}

	port, err := s.randomPort.GetPort()
	if err != nil {
		return fmt.Errorf("Open proxy port: %s", err)
	}
	s.ui.Println("checking host key")
	err = s.cli.Run([]string{
		"-T",
		fmt.Sprintf("%s@%s", jumpBoxUserName, jumpBoxIP),
		"-i", adminSSHPrivateKey,
		"echo", "host key confirmed",
	})
	if err != nil {
		return fmt.Errorf("unable to verify host key fingerprint: %s", err)
	}
	backgroundTunnel, err := s.cli.Start([]string{
		"-4",
		"-D", port,
		"-nNC",
		fmt.Sprintf("%s@%s", jumpBoxUserName, jumpBoxIP),
		"-i", adminSSHPrivateKey,
	})
	if err != nil {
		return fmt.Errorf("Open tunnel to jumpbox: %s", err)
	}
	defer func() {
		if backgroundTunnel != nil {
			backgroundTunnel.Process.Signal(syscall.SIGINT)
		} // removing this will break the acceptance test
	}()

	proxyCommandPrefix := "nc -x"
	time.Sleep(5 * time.Second) // make sure we give that tunnel a moment to open
	return s.cli.Run([]string{
		"-tt",
		"-o", "StrictHostKeyChecking=no",
		"-o", "ServerAliveInterval=300",
		"-o", fmt.Sprintf("ProxyCommand=%s localhost:%s %%h %%p", proxyCommandPrefix, port),
		"-i", adminSSHPrivateKey,
		fmt.Sprintf("%s@%s", targetVMUserName, targetVMIP),
	})
}
