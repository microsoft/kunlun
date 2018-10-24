package apis

type ValueGenerator interface {
	Generate(interface{}) (interface{}, error)
}
