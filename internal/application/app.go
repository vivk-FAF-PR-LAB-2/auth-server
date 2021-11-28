package application

import (
	"context"
	"github.com/spf13/viper"
	"inter-protocol-auth-server/internal/auth/authorization"
	"inter-protocol-auth-server/internal/auth/repository"
	"inter-protocol-auth-server/internal/controller"
	authorization2 "inter-protocol-auth-server/pkg/auth/authorization"
	"inter-protocol-auth-server/pkg/mongo"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IApp interface {
	Start()
	Shutdown()
}

type serverApp struct {
	server *http.Server

	authorizer authorization2.IAuthorization
}

func New(ctx context.Context) IApp {
	router := gin.New()

	db := mongo.Init(viper.GetString("mongo.uri"), viper.GetString("mongo.name"))

	userRepo := repository.NewUserRepository(db, viper.GetString("mongo.collection"))
	authorizer := authorization.NewAuthorizer(
		userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetDuration("auth.token_ttl")*time.Second,
	)

	ctrl := controller.NewController(authorizer)
	ctrl.RegisterRoutes(router)

	return &serverApp{
		server: &http.Server{
			Addr:    ":" + viper.GetString("port"),
			Handler: router,
		},
		authorizer: authorizer,
	}
}

func (app *serverApp) Start() {
	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error while running `auth` server: %v\n", err)
	}
}

func (app *serverApp) Shutdown() {
	if err := app.server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Unable to shutdown `auth` server: %v\n", err)
	}
}
