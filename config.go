package main

import "fmt"

type PostgressConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func (c PostgressConfig) Dialect() string {
	return "postgres"
}

func (c PostgressConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Dbname)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Dbname)
}

func DefaultPostgressConfig() PostgressConfig {
	return PostgressConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "lenslocked_dev",
	}
}

type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmacKey"`
}

func DefaultConfig() Config {
	return Config{
		Port:    4000,
		Env:     "dev",
		Pepper:  "mUGD8rTdJe",
		HMACKey: "the-secret-key",
	}
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

/*



db, err := gorm.Open(dialect, connectionInfo)
const userPwPepper = "mUGD8rTdJe"
const hmacSecretKey = "the-secret-key"

db.LogMode(true)
*/
