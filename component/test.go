package component

type TestSuite struct {
	Component *Component
	Tests []*Test
}

type Test struct {
	Params map[string]interface{}
	Expected []byte
}
