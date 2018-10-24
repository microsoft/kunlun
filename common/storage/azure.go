package storage

type Azure struct {
	Region         string `json:"region,omitempty"`
	Environment    string `json:"environment,omitempty"`
	SubscriptionID string `json:"-"`
	TenantID       string `json:"-"`
	ClientID       string `json:"-"`
	ClientSecret   string `json:"-"`
}
