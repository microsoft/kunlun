package config

type GlobalFlags struct {
	Help      bool   `short:"h" long:"help"`
	Debug     bool   `short:"d" long:"debug"        env:"KL_DEBUG"`
	Version   bool   `short:"v" long:"version"`
	NoConfirm bool   `short:"n" long:"no-confirm"`
	StateDir  string `short:"s" long:"state-dir"    env:"KL_STATE_DIRECTORY"`
	EnvID     string `          long:"name"`
	IAAS      string `          long:"iaas"         env:"KL_IAAS"`

	AzureEnvironment    string `long:"azure-environment"      env:"KL_AZURE_ENVIRONMENT"`
	AzureRegion         string `long:"azure-region"           env:"KL_AZURE_REGION"`
	AzureSubscriptionID string `long:"azure-subscription-id"  env:"KL_AZURE_SUBSCRIPTION_ID"`
	AzureTenantID       string `long:"azure-tenant-id"        env:"KL_AZURE_TENANT_ID"`
	AzureClientID       string `long:"azure-client-id"        env:"KL_AZURE_CLIENT_ID"`
	AzureClientSecret   string `long:"azure-client-secret"    env:"KL_AZURE_CLIENT_SECRET"`
}
