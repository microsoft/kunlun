package common

type FrameworkName string

type Framework interface {
	DetectConfig(path string) []Database
	GetName() FrameworkName
	GetProgrammingLanguage() ProgrammingLanguage
}
