package virtualbox

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// VirtualMachine ...
type VirtualMachine struct {
	UUID       string
	Name       string
	GuestOS    string
	ConfigFile string
	BaseFolder string
}

// IVBoxManage ...
type IVBoxManage interface {
	CreateVM(vm VirtualMachine, shouldRegister bool) (string, error)
	StartVM(name string) error
	AddStorageCtl(vmName string, name string, ctlType string, controller string) error
	AttachStorage(vmName string, controllerName string, port int32, device int32, storageType string, medium string) error
	CreateMedium(mediumType string, filePath string, size int32, format string) error
	VMInfo(name string) (*VirtualMachine, error)
	UnRegisterVM(name string, deleteFiles bool) error
}

// VBoxManage ...
type vBoxManage struct {
}

func (m *vBoxManage) CreateVM(vm VirtualMachine, shouldRegister bool) (string, error) {
	cmd := exec.Command("VBoxManage")
	cmd.Args = append(cmd.Args, "createvm")
	cmd.Args = append(cmd.Args, "--name")
	cmd.Args = append(cmd.Args, vm.Name)
	if vm.GuestOS != "" {
		cmd.Args = append(cmd.Args, "--ostype")
		cmd.Args = append(cmd.Args, vm.GuestOS)
	}

	if vm.BaseFolder != "" {
		cmd.Args = append(cmd.Args, "--basefolder")
		cmd.Args = append(cmd.Args, vm.BaseFolder)
	}

	if shouldRegister {
		cmd.Args = append(cmd.Args, "--register")
	}
	out, _, err := m.execute(cmd)
	if err != nil {
		return "", err
	}
	info := m.parseOutput(out)
	return info["UUID"], nil
}

func (m *vBoxManage) StartVM(name string) error {
	cmd := exec.Command("VBoxManage")
	cmd.Args = append(cmd.Args, "startvm")
	cmd.Args = append(cmd.Args, name)
	_, _, err := m.execute(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *vBoxManage) AddStorageCtl(vmName string, name string, ctlType string, controller string) error {
	cmd := exec.Command("VBoxManage")
	cmd.Args = append(cmd.Args, "storagectl")
	cmd.Args = append(cmd.Args, vmName)
	cmd.Args = append(cmd.Args, "--name")
	cmd.Args = append(cmd.Args, name)
	cmd.Args = append(cmd.Args, "--add")
	cmd.Args = append(cmd.Args, ctlType)
	cmd.Args = append(cmd.Args, "--controller")
	cmd.Args = append(cmd.Args, controller)
	_, _, err := m.execute(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *vBoxManage) AttachStorage(vmName string, controllerName string, port int32, device int32, storageType string, medium string) error {
	cmd := exec.Command("VBoxManage")
	cmd.Args = append(cmd.Args, "storageattach")
	cmd.Args = append(cmd.Args, vmName)
	cmd.Args = append(cmd.Args, "--storagectl")
	cmd.Args = append(cmd.Args, controllerName)
	cmd.Args = append(cmd.Args, "--port")
	cmd.Args = append(cmd.Args, fmt.Sprintf("%d", port))
	cmd.Args = append(cmd.Args, "--device")
	cmd.Args = append(cmd.Args, fmt.Sprintf("%d", device))
	cmd.Args = append(cmd.Args, "--type")
	cmd.Args = append(cmd.Args, storageType)
	cmd.Args = append(cmd.Args, "--medium")
	cmd.Args = append(cmd.Args, medium)
	_, _, err := m.execute(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *vBoxManage) CreateMedium(mediumType string, filePath string, size int32, format string) error {
	cmd := exec.Command("VBoxManage")
	cmd.Args = append(cmd.Args, "createmedium")
	cmd.Args = append(cmd.Args, mediumType)
	cmd.Args = append(cmd.Args, "--filename")
	cmd.Args = append(cmd.Args, filePath)
	cmd.Args = append(cmd.Args, "--size")
	cmd.Args = append(cmd.Args, fmt.Sprintf("%d", size))
	cmd.Args = append(cmd.Args, "--format")
	cmd.Args = append(cmd.Args, format)
	_, _, err := m.execute(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *vBoxManage) VMInfo(name string) (*VirtualMachine, error) {
	cmd := exec.Command("VBoxManage")
	cmd.Args = append(cmd.Args, "showvminfo")
	cmd.Args = append(cmd.Args, name)
	out, _, err := m.execute(cmd)
	if err != nil {
		return nil, err
	}
	vmInfo := m.parseOutput(out)
	return &VirtualMachine{
		UUID:       vmInfo["UUID"],
		Name:       vmInfo["Name"],
		GuestOS:    vmInfo["Guest OS"],
		ConfigFile: vmInfo["Config File"],
		BaseFolder: filepath.Dir(vmInfo["Config File"]),
	}, nil
}

func (m *vBoxManage) UnRegisterVM(name string, deleteFiles bool) error {
	cmd := exec.Command("VBoxManage")
	cmd.Args = append(cmd.Args, "unregistervm")
	cmd.Args = append(cmd.Args, name)
	if deleteFiles {
		cmd.Args = append(cmd.Args, "--delete")

	}
	_, _, err := m.execute(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (m *vBoxManage) parseOutput(out string) map[string]string {
	s := bufio.NewScanner(strings.NewReader(out))
	line := regexp.MustCompile(`(.*): *(.*)`)
	vmInfo := map[string]string{}
	for s.Scan() {
		match := line.FindStringSubmatch(s.Text())
		if match == nil {
			continue
		}

		if len(match) < 3 {
			continue
		}
		vmInfo[match[1]] = match[2]
	}
	return vmInfo
}

func (m *vBoxManage) execute(cmd *exec.Cmd) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// NewVBoxManage ...
func NewVBoxManage() IVBoxManage {
	return &vBoxManage{}
}
