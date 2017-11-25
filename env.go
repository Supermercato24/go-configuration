// Package config implements methods for get config variables.
//
// env file.
package configuration

import (
	"errors"
	"strings"
)

const (
	EnvLocal      = "local"
	EnvStaging    = "staging"
	EnvProduction = "production"
	envFile       = ".env"
	envRegexp     = `(?:APP_ENV=(?P<app_env>local|staging|production))|
(?:DB_HOST=(?P<db_host>[\S]+))|
(?:DB_DATABASE=(?P<db_database>[\S]+))|
(?:DB_USERNAME=(?P<db_username>[\S]+))|
(?:DB_PASSWORD=(?P<db_password>[\S]+))|
(?:MAIL_USERNAME=(?P<mail_username>[\S]+))|
(?:MAIL_PASSWORD=(?P<mail_password>[\S]+))|
(?:MAIL_ENABLED=(?P<mail_enabled>true|false))`
)

var (
	envStruct env
)

type env struct {
	env  string
	Db   db
	Mail mail
}

func (e *env) setEnv(env string) {
	e.env = env
}

func (e env) IsLocal() bool {
	return e.env == EnvLocal
}

func (e env) IsStaging() bool {
	return e.env == EnvStaging
}

func (e env) IsProduction() bool {
	return e.env == EnvProduction
}

type credential struct {
	username string
	password string
}

func (c credential) Username() string {
	return c.username
}

func (c credential) Password() string {
	return c.password
}

type db struct {
	Host        string
	Database    string
	Credentials credential
}

type mail struct {
	enabled     bool
	Credentials credential
}

func (m mail) IsEnabled() bool {
	return m.enabled
}

func init() {
	buffer, err := cfg(envFile, envRegexp)

	if err != nil {
		return
	}

	envStruct = env{
		env: string(buffer["app_env"]),
		Db: db{
			Host:     string(buffer["db_host"]),
			Database: string(buffer["db_database"]),
			Credentials: credential{
				username: string(buffer["db_username"]),
				password: string(buffer["db_password"]),
			},
		},
		Mail: mail{
			enabled: strings.ToLower(string(buffer["mail_enabled"])) == "true",
			Credentials: credential{
				username: string(buffer["mail_username"]),
				password: string(buffer["mail_password"]),
			},
		},
	}
}

// Expose env configuration.
func EnvConfiguration() *env {
	if (env{}) == envStruct {
		err := errors.New("empty env configuration")
		panic(err)
	}

	return &envStruct
}