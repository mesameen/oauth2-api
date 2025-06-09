package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type CommonConfiguration struct {
	Port int `default:"5000"`
}

type OAuthConfiguration struct {
	ClientID          string
	ClientSecret      string
	ClientCallbackURL string
	SessionKey        string
}

var CommonConfig CommonConfiguration
var OAuthConfig OAuthConfiguration

func InitConfig() {
	if err := envconfig.Process("", &CommonConfig); err != nil {
		log.Panicf("%v", err)
	}
	log.Printf("CommonConfig : %+v", CommonConfig)

	if err := envconfig.Process("OAUTH", &OAuthConfig); err != nil {
		log.Panicf("%v", err)
	}
	log.Printf("OAuthConfig : %+v", OAuthConfig)
}
