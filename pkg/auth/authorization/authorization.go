package authorization

import (
	"context"
	"inter-protocol-auth-server/pkg/auth/credential"
)

type IAuthorization interface {
	SignUp(ctx context.Context, user credential.ICredential) error
	SignIn(ctx context.Context, user credential.ICredential) (string, error)
}
