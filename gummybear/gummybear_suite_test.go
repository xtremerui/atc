package gummybear_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGummybear(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gummybear Suite")
}
