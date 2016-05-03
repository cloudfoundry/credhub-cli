package cmclitest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCmclitest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmclitest Suite")
}
