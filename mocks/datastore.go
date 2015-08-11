package mocks

type DataStore struct {
	GetCall struct {
		Args struct {
			Key string
		}
		Return struct {
			Value []byte
			Error error
		}
	}

	SetCall struct {
		Args struct {
			Key   string
			Value []byte
		}
		Return struct {
			Error error
		}
	}
}

func (s *DataStore) Get(key string) ([]byte, error) {
	s.GetCall.Args.Key = key
	return s.GetCall.Return.Value, s.GetCall.Return.Error
}

func (s *DataStore) Set(key string, value []byte) error {
	s.SetCall.Args.Key = key
	s.SetCall.Args.Value = value
	return s.GetCall.Return.Error
}
