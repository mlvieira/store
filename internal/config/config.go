package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Port int
	Env  string
	API  string
	DB   struct {
		DSN string
	}
	Stripe struct {
		Secret string
		Key    string
	}
}

func NewConfig() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.Port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.Env, "env", "development", "Application enviroment {development|production}")
	flag.StringVar(&cfg.DB.DSN, "dsn", "dev:dev@tcp(localhost:3306)/store?parseTime=true&tls=false", "DSN")
	flag.StringVar(&cfg.API, "api", "http://localhost:4001", "URL to api")

	flag.Parse()

	cfg.Stripe.Key = os.Getenv("STRIPE_KEY")
	cfg.Stripe.Secret = os.Getenv("STRIPE_SECRET_KEY")

	return cfg
}

func NewLoggers() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog
}
