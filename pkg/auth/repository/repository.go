package repository

import (
	"context"
	"inter-protocol-auth-server/pkg/auth/credential"
)

type IRepository interface {
	Insert(ctx context.Context, user credential.ICredential) error
	Get(ctx context.Context, username, password string) (credential.ICredential, error)
}
