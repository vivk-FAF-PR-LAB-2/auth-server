package authorization

import (
	"context"
	"inter-protocol-auth-server/pkg/auth/credential"
)

type IAuthorization interface {
	SignUp(ctx context.Context, user credential.ICredential) error
	Confirm(ctx context.Context, token string) error
	SignIn(ctx context.Context, user credential.ICredential) (string, error)
}
