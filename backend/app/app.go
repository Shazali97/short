package app

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/short/dep"
	"github.com/short-d/short/dep/provider"
)

// ServiceConfig represents require parameters for the backend APIs
type ServiceConfig struct {
	LogPrefix            string
	LogLevel             fw.LogLevel
	MigrationRoot        string
	RecaptchaSecret      string
	GithubClientID       string
	GithubClientSecret   string
	FacebookClientID     string
	FacebookClientSecret string
	FacebookRedirectURI  string
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURI    string
	JwtSecret            string
	WebFrontendURL       string
	GraphQLAPIPort       int
	HTTPAPIPort          int
	KeyGenBufferSize     int
	KgsHostname          string
	KgsPort              int
}

// Start launches the GraphQL & HTTP APIs
func Start(
	dbConfig fw.DBConfig,
	config ServiceConfig,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
) {
	db, err := dbConnector.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	err = dbMigrationTool.MigrateUp(db, config.MigrationRoot)
	if err != nil {
		panic(err)
	}

	graphqlAPI, err := dep.InjectGraphQLService(
		"GraphQL API",
		provider.LogPrefix(config.LogPrefix),
		config.LogLevel,
		db,
		"/graphql",
		provider.ReCaptchaSecret(config.RecaptchaSecret),
		provider.JwtSecret(config.JwtSecret),
		provider.KeyGenBufferSize(config.KeyGenBufferSize),
		provider.KgsRPCConfig{
			Hostname: config.KgsHostname,
			Port:     config.KgsPort,
		},
	)
	if err != nil {
		panic(err)
	}
	graphqlAPI.Start(config.GraphQLAPIPort)

	httpAPI := dep.InjectRoutingService(
		"Routing API",
		provider.LogPrefix(config.LogPrefix),
		config.LogLevel,
		db,
		provider.GithubClientID(config.GithubClientID),
		provider.GithubClientSecret(config.GithubClientSecret),
		provider.FacebookClientID(config.FacebookClientID),
		provider.FacebookClientSecret(config.FacebookClientSecret),
		provider.FacebookRedirectURI(config.FacebookRedirectURI),
		provider.GoogleClientID(config.GoogleClientID),
		provider.GoogleClientSecret(config.GoogleClientSecret),
		provider.GoogleRedirectURI(config.GoogleRedirectURI),
		provider.JwtSecret(config.JwtSecret),
		provider.WebFrontendURL(config.WebFrontendURL),
	)
	httpAPI.StartAndWait(config.HTTPAPIPort)
}
