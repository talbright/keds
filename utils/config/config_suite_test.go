package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//TODO add tests
	// . "github.com/talbright/keds/utils/config"

	"testing"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}
