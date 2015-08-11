package application

import (
	"crypto/rand"
	"encoding/base64"
)

type KeyGenerator struct {
	DataStore dataStore
}

func (g KeyGenerator) New() (string, error) {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	newKey := base64.URLEncoding.EncodeToString(randomBytes)
	err = g.DataStore.Set(newKey, []byte("[]"))
	if err != nil {
		return "", err
	}
	return newKey, nil
}
