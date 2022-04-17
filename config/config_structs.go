package config

import (
	"time"
)

// DatabaseOptions ...
type DatabaseOptions struct {
	URL string
}

// RedisStoreOptions ...
type RedisStoreOptions struct {
	Size     int
	Network  string
	Address  string
	Password string
	KeyPairs string
}

// MailOptions options required to set up smtp mailing
type MailOptions struct {
	SenderName string
	FromEmail  string
	Host       string
	Port       int
	Username   string
	Password   string
}

// AuthMailTemplates a struct containing links to template files
/*
	emailVerifTxt and emailVerifHTML templates should expect data containing fields
	{Appname string, Name string, VerificationLink string}
*/
/*
	passwordResetHTML and passwordResetTxt templates should expect data containing fields
	{Appname string, Name  string, ResetLink string}
*/
type AuthMailTemplates struct {
	EmailVerifTxt         string
	EmailVerifHTML        string
	PasswordResetCodeTxt  string
	PasswordResetCodeHTML string
	PasswordResetLinkTxt  string
	PasswordResetLinkHTML string
}

// CORSConfig contains the settings used to setup CORS for the app
type CORSConfig struct {
	AllowAllOrigins bool
	AllowOrigins    []string

	// AllowMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowMethods []string

	// AllowHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	AllowHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposeHeaders []string

	// Allows to add origins like http://some-domain/*, https://api.* or http://some.*.subdomain.com
	AllowWildcard bool
}

// LoggerConfig is the configuration for the log file rotation
type LoggerConfig struct {
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int //days
	Compress   bool
}
