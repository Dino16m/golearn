package config

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

// LoggerConfig is the configuration for the log file rotation
type LoggerConfig struct {
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int //days
	Compress   bool
}
