package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	MongoDBConnectionMethod string // manual
	MongoDBPort             string // 27017
	MongoDBHost             string // 127.0.0.1
	MongoDBUsername         string // rootuser
	MongoDBPassword         string // rootpass
	MongoDBConnectionString string // mongodb connection string will be used when MongoDBConnectionMethod is set to auto
	RedisHost               string // 127.0.0.1
	RedisPort               string // 6379
	AdminPort               string // 6002
}

func LoadENV() *ENV {
	godotenv.Load()
	var env ENV

	var value string
	var found bool

	value, found = os.LookupEnv("MONGODB_CONNECTION_METHOD")
	if found {
		env.MongoDBConnectionMethod = value
	} else {
		env.MongoDBConnectionMethod = "manual"
	}

	value, found = os.LookupEnv("MONGODB_PORT")
	if found {
		env.MongoDBPort = value
	} else {
		env.MongoDBPort = "27017"
	}

	value, found = os.LookupEnv("MONGODB_HOST")
	if found {
		env.MongoDBHost = value
	} else {
		env.MongoDBHost = "127.0.0.1"
	}

	value, found = os.LookupEnv("MONGODB_USERNAME")
	if found {
		env.MongoDBUsername = value
	} else {
		env.MongoDBUsername = "rootuser"
	}

	value, found = os.LookupEnv("MONGODB_PASSWORD")
	if found {
		env.MongoDBPassword = value
	} else {
		env.MongoDBPassword = "rootpass"
	}

	value, found = os.LookupEnv("MONGODB_CONNECTION_STRING")
	if found {
		env.MongoDBConnectionString = value
	} else {
		env.MongoDBConnectionString = ""
	}

	value, found = os.LookupEnv("REDIS_HOST")
	if found {
		env.RedisHost = value
	} else {
		env.RedisHost = "127.0.0.1"
	}

	value, found = os.LookupEnv("REDIS_PORT")
	if found {
		env.RedisPort = value
	} else {
		env.RedisPort = "6379"
	}

	value, found = os.LookupEnv("ADMIN_PORT")
	if found {
		env.AdminPort = value
	} else {
		env.AdminPort = "10223"
	}

	return &env
}
