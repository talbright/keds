package system_test

import (
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/talbright/keds/utils/system"

	"testing"
)

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Suite")
}

var _ = Describe("system", func() {
	Describe("AbsPathify", func() {
		It("should expand the $HOME variable", func() {
			p1 := "$HOME/foo"
			expected := UserHomeDir() + "/foo"
			Expect(AbsPathify(p1)).Should(Equal(expected))
		})
		It("should expand the path", func() {
			_, file, _, _ := runtime.Caller(0)
			p1 := "foo"
			expected := filepath.Dir(file) + "/foo"
			Expect(AbsPathify(p1)).Should(Equal(expected))
		})
		It("shouldn't expand a root based path", func() {
			p1 := "/foo"
			Expect(AbsPathify(p1)).Should(Equal(p1))
		})
	})
})
