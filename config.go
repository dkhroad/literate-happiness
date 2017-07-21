package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

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
	Port       int            `json:"port"`
	Env        string         `json:"env"`
	PepperHash string         `json:"pepper_hash"`
	HMACKey    string         `json:"hmac_key"`
	Database   PostgresConfig `json:"database"`
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
		Database:   DefaultPostgresConfig(),
	}
}

func LoadConfig(isProd bool) (*Config, error) {
	var cfg Config
	if isProd {
		log.Println("Loading config from file config.json")
		f, err := os.Open("./config.json")
		if err != nil {
			return nil, err
		}
		jsonDecoder := json.NewDecoder(f)
		if err = jsonDecoder.Decode(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	log.Println("Loading default config")
	cfg = DefaultConfig()
	return &cfg, nil
}
