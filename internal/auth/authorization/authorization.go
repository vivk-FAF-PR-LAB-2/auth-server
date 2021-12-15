package authorization

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/spf13/viper"
	"inter-protocol-auth-server/internal/auth/claims"
	"inter-protocol-auth-server/pkg/auth/credential"
	"inter-protocol-auth-server/pkg/auth/repository"
	"inter-protocol-auth-server/pkg/rsa"
	"inter-protocol-auth-server/pkg/sendsmtp"
	"time"
)

type Authorizer struct {
	repo repository.IRepository

	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration

	sender sendsmtp.ISender
}

func NewAuthorizer(repo repository.IRepository, hashSalt string, signingKey []byte, expireDuration time.Duration) *Authorizer {
	sender := sendsmtp.NewSender(
		viper.GetString("email.from"),
		viper.GetString("email.password"),
		viper.GetString("email.host"),
		viper.GetString("email.port"))

	return &Authorizer{
		repo:           repo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
		sender:         sender,
	}
}

func (a *Authorizer) SignUp(ctx context.Context, user credential.ICredential) error {
	pwd := sha1.New()
	pwd.Write([]byte(user.GetPassword()))
	pwd.Write([]byte(a.hashSalt))
	user.SetPassword(fmt.Sprintf("%x", pwd.Sum(nil)))

	publicKeyPath, privateKeyPath := viper.GetString("rsa.public_key"), viper.GetString("rsa.private_key")
	privateKey := rsa.GenerateKeyPair(publicKeyPath, privateKeyPath)

	user.SetEmail(rsa.Encrypt(privateKey, user.GetEmail()))

	a.sender.Send(rsa.Decrypt(privateKey, user.GetEmail()),
		"Sign Up",
		"Hello, "+user.GetLogin()+"!"+"\n"+
			"Thank you for registering with our service.")

	return a.repo.Insert(ctx, user)
}

func (a *Authorizer) SignIn(ctx context.Context, user credential.ICredential) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.GetPassword()))
	pwd.Write([]byte(a.hashSalt))
	user.SetPassword(fmt.Sprintf("%x", pwd.Sum(nil)))

	user, err := a.repo.Get(ctx, user.GetLogin(), user.GetPassword())
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: user.GetLogin(),
	})

	tokenSigned, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}

	publicKeyPath, privateKeyPath := viper.GetString("rsa.public_key"), viper.GetString("rsa.private_key")
	privateKey := rsa.GenerateKeyPair(publicKeyPath, privateKeyPath)

	// user.SetEmail(rsa.Encrypt(privateKey, user.GetEmail()))

	result := map[string]string{
		"token": tokenSigned,
		"email": rsa.Decrypt(privateKey, user.GetEmail()),
	}

	jsonBytes, err := json.Marshal(result)
	return string(jsonBytes), err
}
