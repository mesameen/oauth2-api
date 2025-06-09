package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/mesameen/oauth2-api/internal/config"
	"github.com/mesameen/oauth2-api/internal/logger"
	"github.com/mesameen/oauth2-api/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}
	err = logger.InitiLogger()
	if err != nil {
		log.Panic(err)
	}
	// initializing the config
	config.InitConfig()
	goth.UseProviders(
		google.New(config.OAuthConfig.ClientID, config.OAuthConfig.ClientSecret, config.OAuthConfig.ClientCallbackURL),
	)
	r := gin.Default()
	handler, err := routes.NewHandler()
	if err != nil {
		logger.Panicf("%v", err)
	}
	r.LoadHTMLFiles("templates/*")
	r.GET("/", handler.Home)
	r.GET("/api/auth/:provider", handler.SignInWithProvider)
	r.GET("/api/auth/:provider/callback", handler.AuthProviderCallback)
	r.GET("/api/success", handler.Success)
	r.GET("/api/logout", handler.Logout)
	r.Use(handler.AuthMiddleware)
	r.GET("/api/protected", handler.Protected)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.CommonConfig.Port),
		Handler: r,
	}
	logger.Infof("App us up and running on :%v", config.CommonConfig.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Panicf("Failed to up and run server. Error: %v", err)
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Infof("Shutting down server")
	if err := server.Shutdown(ctxTimeout); err != nil {
		logger.Errorf("Failed to shutdown server. Error: %v", err)
	}
}
