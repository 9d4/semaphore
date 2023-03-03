package database

type RedisConfig struct {
	Address  string
	Username string
	Password string

	DB int
}
