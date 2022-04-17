package config

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/spf13/viper"
	"github.com/dino16m/golearn-core/config"
)

// SuperConfig is the default exported settings struct for the auth app.
type SuperConfig struct {
	JwtOptions        config.JWTConfig
	RedisStoreOptions RedisStoreOptions `mapstructure:"Redis"`
	AppName           string
	SessionOptions    sessions.Options  `mapstructure:"Session"`
	MailOptions       MailOptions       `mapstructure:"Mail"`
	AuthMailTemplates AuthMailTemplates `mapstructure:"Authtemplates"`
	SecretKey         string
	AppURL            string
	Env               string
	APICORSConfig     CORSConfig
	WebCORSConfig     CORSConfig
	LoggerConfig      LoggerConfig
	DatabaseOptions   DatabaseOptions
}

// Config is the struct containing structs of all other configs
var Config SuperConfig

// IsSet is used to check if the app has been configured properly yet
func IsSet() bool {
	return !reflect.ValueOf(Config).IsZero()
}

// Initialize sets up config from file. It handles setting reasonable defaults
// as well as raising errors when sensible defaults are unavailable
func Initialize(cfgType string, path string) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	samesiteStr := viper.GetString("session.samesite")
	samesite := mapSameSite(samesiteStr)
	viper.Set("session.samesite", samesite)
	viper.Unmarshal(&Config)
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
