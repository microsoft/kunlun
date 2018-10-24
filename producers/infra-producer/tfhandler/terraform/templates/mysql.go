package templates

import (
	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/helpers"
)

var mysqlTF = []byte(`
resource "azurerm_mysql_server" "{{.mysqlServerName}}" {
	name                      = "${var.env_name}-{{.mysqlServerName}}"
	resource_group_name       = "${azurerm_resource_group.kunlun_resource_group.name}"
	location                  = "${var.location}"
  
	sku {
	  name = "${var.{{.mysqlServerName}}_ams_sku_tier_abbr}_${var.{{.mysqlServerName}}_ams_sku_family}_${var.{{.mysqlServerName}}_ams_sku_capacity}"
	  capacity = "${var.{{.mysqlServerName}}_ams_sku_capacity}"
	  tier = "${var.{{.mysqlServerName}}_ams_sku_tier}"
	  family = "${var.{{.mysqlServerName}}_ams_sku_family}"
	}
  
	storage_profile {
	  storage_mb = "${var.{{.mysqlServerName}}_ams_storage_profile_stroage_gb * 1024}"
	  backup_retention_days = "${var.{{.mysqlServerName}}_ams_storage_profile_backup_retention_days}"
	}
  
	administrator_login = "${var.{{.mysqlServerName}}_ams_administrator_login}"
	administrator_login_password = "${var.{{.mysqlServerName}}_ams_administrator_login_password}"
	version = "${var.{{.mysqlServerName}}_ams_version}"
	ssl_enforcement = "${var.{{.mysqlServerName}}_ams_ssl_enforcement}"
}
  
variable "{{.mysqlServerName}}_ams_sku_capacity" {}
variable "{{.mysqlServerName}}_ams_sku_tier" {}
variable "{{.mysqlServerName}}_ams_sku_tier_abbr" {}
variable "{{.mysqlServerName}}_ams_sku_family" {}
variable "{{.mysqlServerName}}_ams_storage_profile_stroage_gb" {}
variable "{{.mysqlServerName}}_ams_storage_profile_backup_retention_days" {}
variable "{{.mysqlServerName}}_ams_administrator_login" {}
variable "{{.mysqlServerName}}_ams_administrator_login_password" {}
variable "{{.mysqlServerName}}_ams_version" {}
variable "{{.mysqlServerName}}_ams_ssl_enforcement" {}

resource "azurerm_mysql_database" "{{.mysqlDatabaseName}}" {
	name = "{{.mysqlDatabaseName}}"
	resource_group_name = "${azurerm_resource_group.kunlun_resource_group.name}"
	server_name = "${azurerm_mysql_server.{{.mysqlServerName}}.name}"
	charset = "utf8"
	collation = "utf8_unicode_ci"
}

`)

var mysqlTFVars = []byte(`
{{.mysqlServerName}}_ams_sku_capacity = "{{.ams_sku_capacity}}" 
{{.mysqlServerName}}_ams_sku_tier = "{{.ams_sku_tier}}"
{{.mysqlServerName}}_ams_sku_tier_abbr = "{{.ams_sku_tier_abbr}}"
{{.mysqlServerName}}_ams_sku_family = "{{.ams_sku_family}}"
{{.mysqlServerName}}_ams_storage_profile_stroage_gb = "{{.ams_storage_profile_stroage_gb}}"
{{.mysqlServerName}}_ams_storage_profile_backup_retention_days = "{{.ams_storage_profile_backup_retention_days}}"
{{.mysqlServerName}}_ams_administrator_login = "{{.ams_administrator_login}}"
{{.mysqlServerName}}_ams_administrator_login_password = "{{.ams_administrator_login_password}}"
{{.mysqlServerName}}_ams_version = "{{.ams_version}}"
{{.mysqlServerName}}_ams_ssl_enforcement = "{{.ams_ssl_enforcement}}"
`)

func NewMysqlTemplate(mysql artifacts.MysqlDatabase) (string, error) {
	return helpers.Render(mysqlTF, getMysqlTFParams(mysql))
}

func NewMysqlInput(mysql artifacts.MysqlDatabase) (string, error) {
	return helpers.Render(mysqlTFVars, getMysqlTFVarsParams(mysql))
}

func getMysqlTFParams(mysql artifacts.MysqlDatabase) map[string]interface{} {
	return map[string]interface{}{
		"mysqlServerName":   mysql.Name,
		"mysqlDatabaseName": mysql.Name,
	}
}

func getMysqlTFVarsParams(mysql artifacts.MysqlDatabase) map[string]interface{} {
	var tierAbbr string
	switch mysql.Tier {
	case "Basic":
		tierAbbr = "B"
	case "GeneralPurpose":
		tierAbbr = "GP"
	case "MemoryOptimized":
		tierAbbr = "MO"
	}

	return map[string]interface{}{
		"mysqlServerName":                           mysql.Name,
		"ams_sku_capacity":                          mysql.Cores,
		"ams_sku_tier":                              mysql.Tier,
		"ams_sku_tier_abbr":                         tierAbbr,
		"ams_sku_family":                            mysql.Family,
		"ams_storage_profile_stroage_gb":            mysql.Storage,
		"ams_storage_profile_backup_retention_days": mysql.BackupRetentionDays,
		"ams_administrator_login":                   mysql.Username,
		"ams_administrator_login_password":          mysql.Password,
		"ams_version":                               mysql.Version,
		"ams_ssl_enforcement":                       mysql.SSLEnforcement,
	}
}
