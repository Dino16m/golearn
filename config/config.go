package config

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/spf13/viper"
)

// SuperConfig is the default exported settings struct for the auth app.
type SuperConfig struct {
	JwtOptions        JwtOptions
	RedisStoreOptions RedisStoreOptions `mapstructure:"redis"`
	AppName           string
	SessionOptions    sessions.Options  `mapstructure:"session"`
	MailOptions       MailOptions       `mapstructure:"mail"`
	AuthMailTemplates AuthMailTemplates `mapstructure:"authtemplates"`
	SecretKey         string
	AppURL            string
	Env               string
	APICORSConfig     CORSConfig
	WebCORSConfig     CORSConfig
	LoggerConfig      LoggerConfig
}

// Config is the struct containing structs of all other configs
var Config SuperConfig

// IsSet is used to check if the app has been configured properly yet
func IsSet() bool {
	return !reflect.ValueOf(Config).IsZero()
}

// Initialize sets up config from file. It handles setting reasonable defaults
// as well as raising errors when sensible defaults are unavalaible
func Initialize(cfgType string, path string) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	setDefaults()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	samesiteStr := viper.GetString("session.samesite")
	samesite := mapSameSite(samesiteStr)
	viper.Set("session.samesite", samesite)
	viper.Unmarshal(&Config)
}

type templatePaths struct {
	text string
	html string
}

var defaultEmailVerif = templatePaths{
	text: "templates/mail/email.txt",
	html: "templates/mail/email.html",
}
var defaultPasswordReset = defaultEmailVerif

func setDefaults() {
	setDefaultTemplates()
}

func setDefaultTemplates() {
	txt := viper.GetString("authtemplates.EmailVerifTxt")
	html := viper.GetString("authtemplates.EmailVerifHTML")
	emailVerif := getTemplatePaths(txt, html, defaultEmailVerif)
	viper.Set("authtemplates.EmailVerifTxt", emailVerif.text)
	viper.Set("authtemplates.EmailVerifHTML", emailVerif.html)

	txt = viper.GetString("authtemplates.PasswordResetLinkTxt")
	html = viper.GetString("authtemplates.PasswordResetLinkHTML")
	passResetLink := getTemplatePaths(txt, html, defaultPasswordReset)

	viper.Set("authtemplates.PasswordResetLinkTxt", passResetLink.text)
	viper.Set("authtemplates.PasswordResetLinkHTML", passResetLink.html)

	txt = viper.GetString("authtemplates.PasswordResetCodeTxt")
	html = viper.GetString("authtemplates.PasswordResetCodeHTML")
	passResetCode := getTemplatePaths(txt, html, defaultPasswordReset)

	viper.Set("authtemplates.PasswordResetCodeTxt", passResetCode.text)
	viper.Set("authtemplates.PasswordResetCodeHTML", passResetCode.html)
}

func mapSameSite(val string) http.SameSite {
	val = strings.ToLower(val)
	var mappedTo http.SameSite
	switch val {
	case "none":
		mappedTo = http.SameSiteNoneMode
	case "strict":
		mappedTo = http.SameSiteStrictMode
	case "lax":
		mappedTo = http.SameSiteLaxMode
	default:
		mappedTo = http.SameSiteDefaultMode
	}
	return mappedTo
}
func getTemplatePaths(txt, html string, defPaths templatePaths) templatePaths {
	templatePaths := templatePaths{}
	if txt == "" && html == "" {
		templatePaths = defPaths
	} else {
		templatePaths.html = html
		templatePaths.text = txt
	}
	return templatePaths
}
