package repository

import (
	"context"
	"errors"
	"inter-protocol-auth-server/internal/models"
	"inter-protocol-auth-server/pkg/auth/credential"
	"inter-protocol-auth-server/pkg/auth/repository"

	authError "inter-protocol-auth-server/internal/auth/error"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection string) repository.IRepository {
	return &UserRepository{
		db: db.Collection(collection),
	}
}

func (r *UserRepository) Insert(ctx context.Context, user credential.ICredential) error {
	newUser := new(models.User)
	newUser.Username = user.GetLogin()
	newUser.Email = user.GetEmail()
	newUser.Password = user.GetPassword()

	var emailTest *models.User = nil

	err := r.db.FindOne(ctx, bson.M{"email": newUser.Email}).Decode(emailTest)
	if err != mongo.ErrNoDocuments {
		errEmailExist := errors.New("user with this email already exist")
		log.Errorf("%s", errEmailExist.Error())
		return errEmailExist
	}

	_, err = r.db.InsertOne(ctx, newUser)
	if err != nil {
		log.Errorf("error on inserting user: %s", err.Error())
		return authError.ErrUserAlreadyExists
	}

	return nil
}

func (r *UserRepository) Get(ctx context.Context, username, password string) (credential.ICredential, error) {
	user := new(models.User)

	if err := r.db.FindOne(ctx, bson.M{"_id": username, "password": password}).Decode(user); err != nil {
		log.Errorf("error occured while getting user from db: %s", err.Error())
		if err == mongo.ErrNoDocuments {
			return nil, authError.ErrUserDoesNotExist
		}

		return nil, err
	}

	return user, nil
}
