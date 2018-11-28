package artifacts

import (
	"regexp"
)

// Database contains information to deploy a database on VM(s)
type MysqlDatabase struct {
	Name                 string                `yaml:"name"`
	Version              string                `yaml:"version"`
	Cores                int                   `yaml:"cores"`
	Tier                 string                `yaml:"tier"`
	Family               string                `yaml:"family"`
	Storage              int                   `yaml:"storage"`
	BackupRetentionDays  int                   `yaml:"backup_retention_days"`
	SSLEnforcement       string                `yaml:"ssl_enforcement"`
	Username             string                `yaml:"username"`
	Password             string                `yaml:"password"`
	MigrationInformation *MigrationInformation `yaml:"migrate_from,omitempty"`
}

type MigrationInformation struct {
	OriginHost     string `yaml:"origin_host"`
	OriginDatabase string `yaml:"origin_database"`
	OriginUsername string `yaml:"origin_username"`
	OriginPassword string `yaml:"origin_password"`
}

func mysqlDatabaseValidator(m Manifest) error {
	for _, db := range m.MysqlDatabases {
		occuredNames := make(map[string]bool)
		if match, _ := regexp.MatchString("^[0-9a-z]+$", db.Name); !match {
			return validationError("database name must be a non-empty string consisting of lowercase letters and numbers, but found \"%s\"")
		}
		if occuredNames[db.Name] {
			return validationError("database name must be unique, but %s occurs at least twice", db.Name)
		}
		occuredNames[db.Name] = true

		if db.Version != "5.6" && db.Version != "5.7" {
			return validationError("the version of database %s must be \"5.6\" or \"5.7\"", db.Name)
		}
		if db.Tier != "Basic" && db.Tier != "GeneralPurpose" && db.Tier != "MemoryOptimized" {
			return validationError("the tier of database %s must be \"Basic\" or \"GeneralPurpost\" or \"MemoryOptimized\"", db.Name)
		}
		if db.Family != "Gen4" && db.Family != "Gen5" {
			return validationError("the family of database %s must be \"Gen4\" or \"Gen5\"", db.Name)
		}
		if db.Storage < 5 || db.Storage > 4096 {
			return validationError("the storage of database %s must between 5 and 4096", db.Name)
		}
		if db.BackupRetentionDays < 7 || db.BackupRetentionDays > 35 {
			return validationError("the backup_retention_days for database %s must between 7 and 35", db.Name)
		}
		if db.SSLEnforcement != "Enabled" && db.SSLEnforcement != "Disabled" {
			return validationError("the ssl_enforcement for database %s must be \"Enabled\" or \"Disabled\"", db.Name)
		}
	}
	return nil
}
