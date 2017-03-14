package token

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const TOKEN_KEY = "token"

func GetTokenFromContext(ctx context.Context) (token string) {
	if md, ok := metadata.FromContext(ctx); ok {
		token = GetTokenFromMetadata(md)
	}
	return
}

func AddTokenToContext(ctx context.Context, token string) context.Context {
	return metadata.NewContext(ctx, CreateMetadataWithToken(token))
}

func CreateMetadataWithToken(token string) metadata.MD {
	return metadata.Pairs(TOKEN_KEY, token)
}

func GetTokenFromMetadata(md metadata.MD) (token string) {
	if val, ok := md[TOKEN_KEY]; ok {
		if len(val) == 1 {
			token = val[0]
		}
	}
	return
}

func AddTokenToHeader(ctx context.Context, token string) {
	header := CreateMetadataWithToken(token)
	grpc.SetHeader(ctx, header)
}
