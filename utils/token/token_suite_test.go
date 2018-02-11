package token_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/talbright/keds/utils/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"testing"
)

func TestToken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Suite")
}

var _ = Describe("token", func() {
	Describe("AddTokenToContext", func() {
		It("should add the token to the context metadata", func() {
			val := "abc"
			ctx := AddTokenToContext(context.Background(), val)
			md, _ := metadata.FromOutgoingContext(ctx)
			toke, _ := md[TOKEN_KEY]
			Expect(toke[0]).Should(Equal(val))
		})
	})
	Describe("GetTokenFromContext", func() {
		It("should get the token from the context metadata", func() {
			val := "abc"
			ctx := AddTokenToContext(context.Background(), val)
			Expect(GetTokenFromContext(ctx)).Should(Equal(val))
		})
	})
})
