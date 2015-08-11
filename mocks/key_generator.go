package mocks

type KeyGenerator struct {
	NewCall struct {
		Return string
		Error  error
	}
}

func (g KeyGenerator) New() (string, error) {
	return g.NewCall.Return, g.NewCall.Error
}
