package conf

import (
	"log"
	"os"
)

type config struct {
	DbUsername string
	DbPassword string
	DbName     string
	DbHostname string
}

// mustGetEnv returns an env variable value if present and fails othwewise
func mustGetEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

// Config keeps an exposed configuration structure
var Config config

// InitConfig populates config variable and supposed to be called when application started
func InitConfig() {
	m := mustGetEnv
	Config = config{
		DbUsername: m("DB_USERNAME"),
		DbPassword: m("DB_PASSWORD"),
		DbName:     m("DB_NAME"),
		DbHostname: m("DB_HOSTNAME"),
	}
}
