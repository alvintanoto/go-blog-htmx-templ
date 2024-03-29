package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/gorilla/securecookie"
)

var ErrMissingEnv = errors.New("missing env variable")

// List of server configurations
type Configurations struct {
	Server   ServerConfigurations   `json:"server"`
	Database DatabaseConfigurations `json:"database"`
	Redis    RedisConfigurations    `json:"redis"`
}

type RedisConfigurations struct {
	Url string `json:"redis_url"`
}

type ServerConfigurations struct {
	Port      int    `json:"port"`
	SecretKey string `json:"secret_key"`
}

type DatabaseConfigurations struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c *Configurations) ReadConfigurations() error {
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		log.Println("port env is empty")
		return ErrMissingEnv
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Println("failed to convert port string: ", err.Error())
		return err
	}

	var secretKey string
	secretKey = os.Getenv("BLOG_SECRET_KEY")
	if secretKey == "" {
		secretKey = string(securecookie.GenerateRandomKey(32))
	}

	dbHostEnv := os.Getenv("DB_HOST")
	dbPortEnv := os.Getenv("DB_PORT")
	dbUserEnv := os.Getenv("DB_USER")
	dbPasswordEnv := os.Getenv("DB_PASSWORD")
	dbNameEnv := os.Getenv("BLOG_DB_NAME")

	redisEnv := os.Getenv("DB_REDIS")

	c.Server.Port = port
	c.Server.SecretKey = secretKey
	c.Database.Host = dbHostEnv
	c.Database.Port = dbPortEnv
	c.Database.User = dbUserEnv
	c.Database.Password = dbPasswordEnv
	c.Database.Name = dbNameEnv
	c.Redis.Url = redisEnv

	return nil
}
