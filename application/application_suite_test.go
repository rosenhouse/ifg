package application_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIFG(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ifg Application Suite")
}
