package templates

var providerTF = []byte(`
provider "azurerm" {
	client_id       = "${var.azure_client_id}"
	client_secret   = "${var.azure_client_secret}"
	subscription_id = "${var.azure_subscription_id}"
	tenant_id       = "${var.azure_tenant_id}"
	environment     = "${var.azure_environment}"
}

variable "azure_client_id" {}
variable "azure_client_secret" {}
variable "azure_subscription_id" {}
variable "azure_tenant_id" {}
variable "azure_environment" {}
`)

func NewProviderTemplate() (string, error) {
	return string(providerTF), nil
}
