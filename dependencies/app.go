// +build wireinject

package dependencies

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dino16m/GinSessionMW/middleware"
	"github.com/dino16m/golearn/adapters"
	"github.com/dino16m/golearn/config"
	"github.com/dino16m/golearn/lib/event"
	"github.com/dino16m/golearn/middlewares"
	"github.com/dino16m/golearn/services"
	"github.com/dino16m/golearn/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	gwa "github.com/gobuffalo/gocraft-work-adapter"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/wire"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func provideRedisStore(cfg config.SuperConfig) (sessions.Store, error) {
	options := cfg.RedisStoreOptions

	return redis.NewStore(
		options.Size,
		options.Network,
		options.Address,
		options.Password,
		[]byte(options.KeyPairs),
	)
}

// initializeSession get a session
func initializeSession() (gin.HandlerFunc, error) {
	wire.Build(
		ProvideSuperConfig, sessions.Sessions,
		wire.FieldsOf(new(config.SuperConfig), "AppName"), provideRedisStore,
	)
	return func(c *gin.Context) {}, nil
}

func provideAuthService(service services.AuthService) middlewares.AuthenticatorFunc {
	return service.Authenticate
}

func initializeJwMw(authService services.AuthService,
) (*jwt.GinJWTMiddleware, error) {
	wire.Build(
		provideAuthService, wire.FieldsOf(new(config.SuperConfig), "JwtOptions"),
		middlewares.GetJwtMiddleware, ProvideSuperConfig)
	return &jwt.GinJWTMiddleware{}, nil
}

func provideSessionFunc() func(c *gin.Context) sessions.Session {
	return sessions.Default
}

func provideAuthUserRepo(repo types.UserRepository) func(key interface{}) types.AuthUser {
	return repo.GetUserByAuthId
}

func initializeSessionAuthMW(
	userRepo types.UserRepository) *middleware.SessionMiddleware {
	wire.Build(
		wire.FieldsOf(new(config.SuperConfig), "SessionOptions"),
		provideSessionFunc, ProvideSuperConfig,
		provideAuthUserRepo, middlewares.GetSessionMw)

	return &middleware.SessionMiddleware{}
}

func provideSessionAuthAdapter(
	mw *middleware.SessionMiddleware) adapters.SessionAuthUserManager {
	return adapters.NewSessionAuthUserManager(mw)
}

func provideIdentityKey(cfg config.SuperConfig) string {
	return cfg.JwtOptions.IdentityKey
}
func initializeJwtAuthAdapter(
	repo types.UserRepository) adapters.JwtAuthUserManager {
	wire.Build(provideIdentityKey, ProvideSuperConfig,
		provideAuthUserRepo, adapters.NewJwtAuthUserManager)
	return adapters.JwtAuthUserManager{}
}

func provideEventDispatcher() *event.AuthEventDispatcher {
	return event.NewAuthEventDispatcher()
}

func ProvideSuperConfig() config.SuperConfig {
	if config.IsSet() == false {
		panic("Config must be set up")
	}
	return config.Config
}

func provideRedisPool(cfg config.SuperConfig) *redigo.Pool {
	options := cfg.RedisStoreOptions
	dialPasswordOption := redigo.DialPassword(options.Password)
	return &redigo.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial(
				options.Network, options.Address, dialPasswordOption)
		},
	}
}

func provideWorker(pool *redigo.Pool, cfg config.SuperConfig) *gwa.Adapter {
	return gwa.New(gwa.Options{
		Pool:           pool,
		Name:           cfg.AppName,
		MaxConcurrency: 25,
	})
}

func provideCSRFMiddleware(cfg config.SuperConfig) middlewares.CSRFMiddleware {
	return middlewares.NewCSRFMiddleware(cfg.SecretKey,
		cfg.Env, cfg.SessionOptions)
}

func provideWebCORSMiddleware(cfg config.SuperConfig) WebCORSMiddleware {
	defaultCfg := cors.DefaultConfig()
	defaultCfg.AllowHeaders = append(defaultCfg.AllowHeaders, "X-CSRF-TOKEN")
	validCfg := updateCORSCfg(cfg.WebCORSConfig, defaultCfg)
	return WebCORSMiddleware(cors.New(validCfg))
}

func provideAPICORSMiddleware(cfg config.SuperConfig) APICORSMiddleware {
	defaultCfg := cors.DefaultConfig()
	defaultCfg.AllowHeaders = append(defaultCfg.AllowHeaders, "X-CSRF-TOKEN")
	validCfg := updateCORSCfg(cfg.APICORSConfig, defaultCfg)
	return APICORSMiddleware(cors.New(validCfg))
}

func updateCORSCfg(update config.CORSConfig, original cors.Config) cors.Config {
	original.AllowAllOrigins = update.AllowAllOrigins
	original.AllowCredentials = update.AllowCredentials
	original.AllowWildcard = update.AllowWildcard
	if len(update.AllowHeaders) > 0 {
		original.AllowHeaders = append(original.AllowHeaders, update.AllowHeaders...)
	}
	if len(update.AllowMethods) > 0 {
		original.AllowMethods = update.AllowMethods
	}
	if len(update.AllowOrigins) > 0 {
		original.AllowOrigins = update.AllowOrigins
	}
	if len(update.ExposeHeaders) > 0 {
		original.ExposeHeaders = update.ExposeHeaders
	}
	return original
}

func provideLogger() *logrus.Logger {
	log := logrus.New()
	cfg := ProvideSuperConfig().LoggerConfig
	if cfg.Filename != "" {
		log.Out = &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
	}
	return log
}

// initializeAll Get all the objects needed to start app
func initializeAll(
	services ServicesContainer,
	repos RepositoriesContainer) (App, error) {
	wire.Build(
		wire.FieldsOf(new(RepositoriesContainer), "UserRepo"),
		wire.FieldsOf(new(ServicesContainer), "AuthService"),
		initializeJwtAuthAdapter, provideSessionAuthAdapter,
		initializeSessionAuthMW, initializeJwMw,
		initializeSession, wire.Struct(new(App), "*"),
		provideEventDispatcher, provideCSRFMiddleware,
		provideRedisPool, provideWorker, ProvideSuperConfig,
		provideAPICORSMiddleware, provideWebCORSMiddleware,
		provideLogger,
		wire.Bind(new(event.Dispatcher), new(*event.AuthEventDispatcher)),
		wire.Bind(new(types.Logger), new(*logrus.Logger)),
	)
	return App{}, nil
}

func InitApp(
	services ServicesContainer,
	repos RepositoriesContainer) (App, error) {
	return initializeAll(services, repos)
}

type APICORSMiddleware gin.HandlerFunc
type WebCORSMiddleware gin.HandlerFunc

// App the fully constructed app, containing required values
type App struct {
	SessionMw             gin.HandlerFunc
	JwtAuthMW             *jwt.GinJWTMiddleware
	SessionAuthMiddleware *middleware.SessionMiddleware
	ApiUserManager        adapters.JwtAuthUserManager
	SessionUserManager    adapters.SessionAuthUserManager
	EventDispatcher       event.Dispatcher
	Worker                *gwa.Adapter
	CSRFMiddleware        middlewares.CSRFMiddleware
	APICORS               APICORSMiddleware
	WebCORS               WebCORSMiddleware
	Logger                types.Logger
}
