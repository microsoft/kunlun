resource "azurerm_storage_account" "default_storage_account" {
  name                      = "${replace(var.env_name, "-", "")}-sa"
  resource_group_name       = "${azurerm_resource_group.kunlun_resource_group.name}"
  location                  = "${var.location}"
  account_kind              = "Storage"
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_blob_encryption    = true
  enable_file_encryption    = true
  enable_https_traffic_only = true
}
