package common

type PackageManagerName string

type PackageManager interface {
	Identify(path string) bool
	DetectFramework(path string) []FrameworkName
	GetName() PackageManagerName
}
