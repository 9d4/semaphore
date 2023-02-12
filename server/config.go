package server

import (
	"reflect"
	"strings"

	"github.com/9d4/semaphore/util"
	"github.com/spf13/viper"
)

var defaultConfig *Config = &Config{
	v:             nil,
	Address:       "0.0.0.0:3500",
	DBHost:        "127.0.0.1",
	DBPort:        5432,
	DBName:        "semaphore",
	DBUsername:    "semaphore",
	DBPassword:    "smphr",
	RedisAddress:  "127.0.0.1:6379",
	RedisUsername: "default",
	RedisPassword: "t00r",
	LogRequest:    false,
}

func init() {
	defaultConfig.Key = util.GenerateKey()
	defaultConfig.KeyBytes = util.StringToBytes(defaultConfig.Key)
}

// Config is configuration for the http server to be able to run.
type Config struct {
	v *viper.Viper

	Key      string
	KeyBytes []byte

	// Address to listen on
	Address    string
	DBHost     string
	DBPort     int
	DBName     string
	DBUsername string
	DBPassword string

	RedisAddress  string
	RedisUsername string
	RedisPassword string

	LogRequest bool
}

func (c *Config) Apply(conf *Config) error {
	if conf != c {
		*conf = *c
	}
	return nil
}

// Option gorm option interface
type Option interface {
	Apply(*Config) error
}

// ParseViper extract configuration from v and returns *Config.
func ParseViper(v *viper.Viper) *Config {
	c := parseViper(v, defaultConfig)

	return c
}

func parseViper(v *viper.Viper, defaultConf *Config) *Config {
	c := &Config{v: v}

	c.v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	c.Key = getOrDefault(v.GetString("app_key"), defaultConf.Key)

	// default of []byte is comparable to nil, that's why if converting empty
	// string yields non nil []byte
	if c.Key != "" {
		c.KeyBytes = getOrDefault([]byte(c.Key), defaultConf.KeyBytes)
	}

	c.Address = getOrDefault(v.GetString("address"), defaultConf.Address)
	c.DBHost = getOrDefault(v.GetString("db-host"), defaultConf.DBHost)
	c.DBPort = getOrDefault(v.GetInt("db-port"), defaultConf.DBPort)
	c.DBName = getOrDefault(v.GetString("db-name"), defaultConf.DBName)
	c.DBUsername = getOrDefault(v.GetString("db-username"), defaultConf.DBUsername)
	c.DBPassword = getOrDefault(v.GetString("db-password"), defaultConf.DBPassword)
	c.RedisAddress = getOrDefault(v.GetString("redis-address"), defaultConf.RedisAddress)
	c.RedisUsername = getOrDefault(v.GetString("redis-username"), defaultConf.RedisUsername)
	c.RedisPassword = getOrDefault(v.GetString("redis-password"), defaultConf.RedisPassword)
	c.LogRequest = getOrDefault(v.GetBool("log-request"), defaultConf.LogRequest)

	return c
}

func getOrDefault[T interface{}](val T, defaultVal T) T {
	refVal := reflect.ValueOf(val)

	if refVal.IsValid() && !refVal.IsZero() {
		return refVal.Interface().(T)
	}

	return defaultVal
}
