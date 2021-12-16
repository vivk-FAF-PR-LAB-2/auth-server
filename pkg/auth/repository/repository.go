package repository

import (
	"context"
	"github.com/google/uuid"
	"inter-protocol-auth-server/pkg/auth/credential"
)

type IRepository interface {
	Insert(ctx context.Context, user credential.ICredential) error
	Get(ctx context.Context, username, password string) (credential.ICredential, error)

	GetConfirmationToken(user credential.ICredential) uuid.UUID
	SetConfirmationToken(token string) (credential.ICredential, bool)
}
