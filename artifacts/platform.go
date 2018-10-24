package apis

// Platform contains information about the project type. Several possible values for type:
// `tomcat/golang/python/java/php/nodejs`
type Platform struct {
	Type string `yaml:"type,omitempty"`
}
