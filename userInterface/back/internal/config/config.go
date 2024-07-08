package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type FiberConfig struct {
	Port string
}

type JwtConfig struct {
	Secret string
}

type RPCConfig struct {
	Port string
}

type Config struct {
	Env     string
	Fiber   FiberConfig
	Jwt     JwtConfig
	Dsn     string
	GethUrl string
	Rpc     RPCConfig
}

func NewConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}

	fiberPort := os.Getenv("FIBER_PORT")
	if fiberPort == "" {
		fiberPort = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable not set")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic("DSN environment variable not set")
	}

	gethUrl := os.Getenv("GETH_URL")
	if gethUrl == "" {
		panic("GETH_URL environment variable not set")
	}

	rpcPort := os.Getenv("RPC_PORT")
	if rpcPort == "" {
		panic("RPC_PORT environment variable not set")
	}

	return &Config{
		Env: env,
		Fiber: FiberConfig{
			Port: fiberPort,
		},
		Jwt: JwtConfig{
			Secret: jwtSecret,
		},
		Dsn:     dsn,
		GethUrl: gethUrl,
		Rpc: RPCConfig{
			Port: rpcPort,
		},
	}
}

func (c Config) GetEnv() string {
	return c.Env
}

func (c Config) GetFiberPort() string {
	return c.Fiber.Port
}

func (c Config) GetJwtSecret() string {
	return c.Jwt.Secret
}

func (c Config) GetDsn() string {
	return c.Dsn
}

func (c Config) GetGethUrl() string {
	return c.GethUrl
}

func (c Config) GetRpcPort() string {
	return c.Rpc.Port
}
