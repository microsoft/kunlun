package templates

import (
	"strings"

	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/helpers"
)

var vmTF = []byte(`
resource "azurerm_availability_set" "{{.vmGroupName}}_as" {
	name                        = "${var.env_name}-{{.vmGroupName}}-as"
	location                    = "${var.location}"
	resource_group_name         = "${azurerm_resource_group.kunlun_resource_group.name}"
	managed                     = "true"
	platform_fault_domain_count = 2
}
  
{{if .publicIPAddressAllocation -}}
resource "azurerm_public_ip" "{{.vmGroupName}}_public_ip" {
	name                         = "${var.env_name}-{{.vmGroupName}}-public-ip-${count.index}"
	count 						 = "${var.{{.vmGroupName}}_avm_count}"
	location                     = "${var.location}"
	resource_group_name          = "${azurerm_resource_group.kunlun_resource_group.name}"
	public_ip_address_allocation = "{{.publicIPAddressAllocation}}"
}
{{- end}}

resource "azurerm_network_interface" "{{.vmGroupName}}_nic" {
	name                      = "${var.env_name}-{{.vmGroupName}}-nic-${count.index}"
	count                     = "${var.{{.vmGroupName}}_avm_count}"
	location                  = "${var.location}"
	resource_group_name       = "${azurerm_resource_group.kunlun_resource_group.name}"
	network_security_group_id = "${azurerm_network_security_group.{{.nsgName}}.id}"
	ip_configuration {
	  name                          = "${var.env_name}-{{.vmGroupName}}-nicip-${count.index}"
	  subnet_id                     = "${azurerm_subnet.{{.subnetName}}.id}"
	  private_ip_address_allocation = "dynamic"
	  {{if .publicIPAddressAllocation -}}
	  public_ip_address_id          = "${azurerm_public_ip.{{.vmGroupName}}_public_ip.*.id[count.index]}"
	  {{- end}}
	}
}
  
resource "azurerm_virtual_machine" "{{.vmGroupName}}" {
	name                  = "${var.env_name}-{{.vmGroupName}}-vm-${count.index}"
	location              = "${var.location}"
	resource_group_name   = "${azurerm_resource_group.kunlun_resource_group.name}"
	network_interface_ids = ["${azurerm_network_interface.{{.vmGroupName}}_nic.*.id[count.index]}"]
	availability_set_id   = "${azurerm_availability_set.{{.vmGroupName}}_as.id}"
	vm_size               = "${var.{{.vmGroupName}}_avm_vm_size}"
	count                 = "${var.{{.vmGroupName}}_avm_count}"
  
	# Comment this line to keep the OS disk after deleting the VM
	delete_os_disk_on_termination = true
  
	# Comment this line to keep the data disks after deleting the VM
	delete_data_disks_on_termination = true
  
	storage_image_reference {
	  publisher = "${var.{{.vmGroupName}}_avm_storage_image_reference_publisher}"
	  offer     = "${var.{{.vmGroupName}}_avm_storage_image_reference_offer}"
	  sku       = "${var.{{.vmGroupName}}_avm_storage_image_reference_sku}"
	  version   = "${var.{{.vmGroupName}}_avm_storage_image_reference_version}"
	}
  
	storage_os_disk {
	  name              = "{{.vmGroupName}}-osdisk-${count.index}"
	  caching           = "${var.{{.vmGroupName}}_avm_storage_os_disk_caching}"
	  create_option     = "${var.{{.vmGroupName}}_avm_storage_os_disk_create_option}"
	  managed_disk_type = "${var.{{.vmGroupName}}_avm_storage_os_disk_managed_disk_type}"
	}

	storage_data_disk {
		name			  = "{{.vmGroupName}}-datadisk-${count.index}"
		lun               = 0
		caching           = "${var.{{.vmGroupName}}_avm_storage_data_disk_caching}"
		create_option     = "${var.{{.vmGroupName}}_avm_storage_data_disk_create_option}"
		managed_disk_type = "${var.{{.vmGroupName}}_avm_storage_data_disk_managed_disk_type}"
		disk_size_gb      = "${var.{{.vmGroupName}}_avm_storage_data_disk_disk_size_gb}"

	}
  
	os_profile {
	  computer_name  = "{{.vmGroupName}}-vm-${count.index}"
	  admin_username = "${var.{{.vmGroupName}}_avm_os_profile_admin_username}"
	}
  
	os_profile_linux_config {
	  disable_password_authentication = true
	  ssh_keys {
		  path = "/home/${var.{{.vmGroupName}}_avm_os_profile_admin_username}/.ssh/authorized_keys"
		  key_data = "${var.{{.vmGroupName}}_avm_os_profile_linux_config_ssh_keys_key_data}"
	  }
	}
}
  
{{if .loadBalancerBackendAddressPoolName}}
resource "azurerm_network_interface_backend_address_pool_association" "{{.vmGroupName}}_backend_address_pool_association" {
	count                   = "${var.{{.vmGroupName}}_avm_count}"
	network_interface_id    = "${azurerm_network_interface.{{.vmGroupName}}_nic.*.id[count.index]}"
	ip_configuration_name   = "${var.env_name}-{{.vmGroupName}}-nicip-${count.index}"
	backend_address_pool_id = "${azurerm_lb_backend_address_pool.{{.loadBalancerBackendAddressPoolName}}.id}"
}
{{end}}

variable "{{.vmGroupName}}_avm_vm_size" {}
variable "{{.vmGroupName}}_avm_count" {}
variable "{{.vmGroupName}}_avm_storage_image_reference_publisher" {}
variable "{{.vmGroupName}}_avm_storage_image_reference_offer" {}
variable "{{.vmGroupName}}_avm_storage_image_reference_sku" {}
variable "{{.vmGroupName}}_avm_storage_image_reference_version" {}
variable "{{.vmGroupName}}_avm_storage_os_disk_caching" {}
variable "{{.vmGroupName}}_avm_storage_os_disk_create_option" {}
variable "{{.vmGroupName}}_avm_storage_os_disk_managed_disk_type" {}
variable "{{.vmGroupName}}_avm_storage_data_disk_caching" {}
variable "{{.vmGroupName}}_avm_storage_data_disk_create_option" {}
variable "{{.vmGroupName}}_avm_storage_data_disk_managed_disk_type" {}
variable "{{.vmGroupName}}_avm_storage_data_disk_disk_size_gb" {}
variable "{{.vmGroupName}}_avm_os_profile_admin_username" {}
variable "{{.vmGroupName}}_avm_os_profile_linux_config_ssh_keys_key_data" {}
`)

var vmTFVars = []byte(`
{{.vmGroupName}}_avm_vm_size = "{{.avm_vm_size}}"
{{.vmGroupName}}_avm_count = "{{.avm_count}}"
{{.vmGroupName}}_avm_storage_image_reference_publisher = "{{.avm_storage_image_reference_publisher}}"
{{.vmGroupName}}_avm_storage_image_reference_offer = "{{.avm_storage_image_reference_offer}}"
{{.vmGroupName}}_avm_storage_image_reference_sku = "{{.avm_storage_image_reference_sku}}"
{{.vmGroupName}}_avm_storage_image_reference_version = "{{.avm_storage_image_reference_version}}"
{{.vmGroupName}}_avm_storage_os_disk_caching = "{{.avm_storage_os_disk_caching}}"
{{.vmGroupName}}_avm_storage_os_disk_create_option = "{{.avm_storage_os_disk_create_option}}"
{{.vmGroupName}}_avm_storage_os_disk_managed_disk_type = "{{.avm_storage_os_disk_managed_disk_type}}"
{{.vmGroupName}}_avm_storage_data_disk_caching = "{{.avm_storage_data_disk_caching}}"
{{.vmGroupName}}_avm_storage_data_disk_create_option = "{{.avm_storage_data_disk_create_option}}"
{{.vmGroupName}}_avm_storage_data_disk_managed_disk_type = "{{.avm_storage_data_disk_managed_disk_type}}"
{{.vmGroupName}}_avm_storage_data_disk_disk_size_gb = "{{.avm_storage_data_disk_disk_size_gb}}"
{{.vmGroupName}}_avm_os_profile_admin_username = "{{.avm_os_profile_admin_username}}"
{{.vmGroupName}}_avm_os_profile_linux_config_ssh_keys_key_data = "{{.avm_os_profile_linux_config_ssh_keys_key_data}}"
`)

var vmOutputTF = []byte(`
output "vm_groups_{{.vmGroupName}}_networks_0_outputs_{{.index}}" {
	value = {
		"ip" = "${azurerm_network_interface.{{.vmGroupName}}_nic.*.private_ip_address[{{.index}}]}"
		{{if .publicIPAddressAllocation -}}
		"public_ip" = "${azurerm_public_ip.{{.vmGroupName}}_public_ip.*.ip_address[{{.index}}]}"
		{{- end}}
	}
}`)

func NewVMTemplate(vm artifacts.VMGroup) (string, error) {
	template, err := helpers.Render(vmTF, getVMTFParams(vm))
	if err != nil {
		return "", err
	}

	for i := 0; i < vm.Count; i++ {
		outputTemplate, err := helpers.Render(vmOutputTF, map[string]interface{}{
			"vmGroupName":               vm.Name,
			"index":                     i,
			"publicIPAddressAllocation": vm.NetworkInfos[0].PublicIP == "static" || vm.NetworkInfos[0].PublicIP == "dynamic",
		})
		if err != nil {
			return "", err
		}
		template += outputTemplate
	}

	return template, nil
}

func NewVMInput(vm artifacts.VMGroup) (string, error) {
	return helpers.Render(vmTFVars, getVMTFVarsParams(vm))
}

func getVMTFParams(vm artifacts.VMGroup) map[string]interface{} {
	m := map[string]interface{}{
		"vmGroupName":                        vm.Name,
		"subnetName":                         vm.NetworkInfos[0].SubnetName,
		"nsgName":                            vm.NetworkInfos[0].NetworkSecurityGroupName,
		"loadBalancerBackendAddressPoolName": vm.NetworkInfos[0].LoadBalancerBackendAddressPoolName,
	}
	if vm.NetworkInfos[0].PublicIP == "static" || vm.NetworkInfos[0].PublicIP == "dynamic" {
		m["publicIPAddressAllocation"] = vm.NetworkInfos[0].PublicIP
	}
	return m
}

func getVMTFVarsParams(vm artifacts.VMGroup) map[string]interface{} {
	return map[string]interface{}{
		"vmGroupName":                                   vm.Name,
		"avm_vm_size":                                   vm.SKU,
		"avm_count":                                     vm.Count,
		"avm_storage_image_reference_publisher":         vm.Storage.Image.Publisher,
		"avm_storage_image_reference_offer":             vm.Storage.Image.Offer,
		"avm_storage_image_reference_sku":               vm.Storage.Image.SKU,
		"avm_storage_image_reference_version":           vm.Storage.Image.Version,
		"avm_storage_os_disk_caching":                   vm.Storage.OSDisk.Caching,
		"avm_storage_os_disk_create_option":             vm.Storage.OSDisk.CreateOption,
		"avm_storage_os_disk_managed_disk_type":         vm.Storage.OSDisk.ManagedDiskType,
		"avm_storage_data_disk_caching":                 vm.Storage.DataDisks[0].Caching,
		"avm_storage_data_disk_create_option":           vm.Storage.DataDisks[0].CreateOption,
		"avm_storage_data_disk_managed_disk_type":       vm.Storage.DataDisks[0].ManagedDiskType,
		"avm_storage_data_disk_disk_size_gb":            vm.Storage.DataDisks[0].DiskSizeGB,
		"avm_os_profile_admin_username":                 vm.OSProfile.AdminName,
		"avm_os_profile_linux_config_ssh_keys_key_data": strings.TrimSpace(vm.OSProfile.LinuxConfiguration.SSH.PublicKeys[0]),
	}
}
