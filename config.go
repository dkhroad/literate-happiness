package main

import "fmt"

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		c.Host, c.Port, c.Name, c.User, c.Password)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Name:     "lenslocked_dev",
		Password: "",
	}
}

type Config struct {
	Port       int    `json:"port"`
	Env        string `json:"env"`
	PepperHash string `json:"pepper_hash"`
	HMACKey    string `json:"hmac_key"`
}

func (c Config) isProd() bool {
	return c.Env == "Prod"
}

func DefaultConfig() Config {
	return Config{
		Port:       3000,
		Env:        "Dev",
		PepperHash: "doormat-wrangle-scam-gating-shelve",
		HMACKey:    "hmac-secret-key",
	}
}

// // TODO: update this to be a config variable
// const (
// 	pepperHash    = "doormat-wrangle-scam-gating-shelve"
// 	hmacSecretKey = "hmac-secret-key"
// )
//

// # services.go
// 	db.LogMode(true)
// 	db, err := gorm.Open("postgres", connectionInfo)
//
