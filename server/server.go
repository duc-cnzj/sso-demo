package server

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"sso/app/models"
	"sso/config/env"
	"sso/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	_ "github.com/go-sql-driver/mysql"
	redis2 "github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Loader interface {
	Load(*Server) error
	GetWeight() int
}

type LoaderCollection []Loader

func (l LoaderCollection) Len() int {
	return len(l)
}

func (l LoaderCollection) Less(i, j int) bool {
	return l[i].GetWeight() < l[j].GetWeight()
}

func (l LoaderCollection) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Server struct {
	rootPath         string
	configPath       string
	env              *env.Env
	config           *env.Config
	db               *gorm.DB
	redis            *redis2.Pool
	session          sessions.Store
	engine           *gin.Engine
	loaders          LoaderCollection
	runningInConsole bool
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) SetRunningInConsole() {
	s.runningInConsole = true
}

func (s *Server) RunningInConsole() bool {
	return s.runningInConsole
}

func (s *Server) Config() *env.Config {
	return s.config
}

func (s *Server) Env() *env.Env {
	return s.env
}

func (s *Server) Init(configPath, rootPath string) error {
	gob.Register(&models.User{})

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	s.configPath = configPath
	s.rootPath = rootPath

	loaders := LoaderCollection{
		&ConfigLoader{},
		&DBLoader{},
		&RedisLoader{},
		&SessionLoader{},
		&EnvLoader{},
		&ValidatorLoader{},
		&EnvironmentLoader{},
	}

	sort.Sort(loaders)

	for _, loader := range loaders {
		if err := loader.Load(s); err != nil {
			log.Debug().Err(err).Msg("loader.Load")
			return err
		}
	}

	s.loaders = loaders

	if s.env.IsDebugging() {
		for _, lo := range s.loaders {
			log.Debug().Msg(fmt.Sprintf("register loader:%30s => weight: %d", reflect.TypeOf(lo).String(), lo.GetWeight()))
		}
	}

	r := gin.New()

	if s.RunningInConsole() {
		s.env.SkipLoadResources()
	}
	s.engine = routes.Init(r, s.env)

	log.Info().Msg("server inited.")

	return nil
}

func (s *Server) Run() error {
	log.Info().Msg(fmt.Sprintf("server running at :%d", s.config.AppPort))

	return s.engine.Run(fmt.Sprintf(":%d", s.config.AppPort))
}

func (s *Server) Shutdown() {

}

func (s *Server) ProductionMode() {
	gin.SetMode(gin.ReleaseMode)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func (s *Server) DebugMode() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	gin.DefaultWriter = log.Logger
	gin.SetMode(gin.DebugMode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Debug().Msgf("route:%10s\t%v", httpMethod, absolutePath)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("############### debug mode ###############")
	s.env.PrintConfig()
	s.db.LogMode(true)
}

func (s *Server) LoadTranslators() *ut.UniversalTranslator {
	zhLang := zh.New()
	enLang := en.New()
	uni := ut.New(enLang, zhLang, enLang)
	return uni
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadConfig(configPath string) (*env.Config, error) {
	var err error
	if configPath == "" {
		configPath = ".env"
	}
	if !path.IsAbs(configPath) {
		getwd, _ := os.Getwd()
		configPath = getwd + "/" + configPath
	}
	exists := fileExists(configPath)
	if !exists {
		return nil, errors.New("file not exists in " + configPath)
	}
	viper.SetConfigType("env")
	file, _ := ioutil.ReadFile(configPath)

	if err = viper.ReadConfig(bytes.NewReader(file)); err != nil {
		return &env.Config{}, err
	}

	config := &env.Config{
		AppPort:             viper.GetUint("APP_PORT"),
		AppEnv:              viper.GetString("APP_ENV"),
		Debug:               viper.GetBool("DEBUG"),
		DBConnection:        viper.GetString("DB_CONNECTION"),
		DBHost:              viper.GetString("DB_HOST"),
		DBPort:              viper.GetUint("DB_PORT"),
		DBDatabase:          viper.GetString("DB_DATABASE"),
		DBUsername:          viper.GetString("DB_USERNAME"),
		DBPassword:          viper.GetString("DB_PASSWORD"),
		RedisHost:           viper.GetString("REDIS_HOST"),
		RedisPassword:       viper.GetString("REDIS_PASSWORD"),
		RedisPort:           viper.GetUint("REDIS_PORT"),
		SessionLifetime:     viper.GetInt("SESSION_LIFETIME"),
		AccessTokenLifetime: viper.GetInt("ACCESS_TOKEN_LIFETIME"),
		ApiTokenLifetime:    viper.GetInt("API_TOKEN_LIFETIME"),
		JwtSecret:           viper.GetString("JWT_SECRET"),
		JwtExpiresSeconds:   viper.GetInt64("JWT_EXPIRES_AT"),
	}

	return config, nil
}
