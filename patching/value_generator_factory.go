package apis

type ValueGeneratorFactory interface {
	GetGenerator(valueType string) (ValueGenerator, error)
}
