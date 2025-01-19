package config

import (
	"os"
	"strconv"
)

// LoadConfig carrega todas as configurações do ambiente
func LoadConfig() *Config {
	return &Config{
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
		RateLimit: RateLimitConfig{
			IPLimit:       getEnvInt("IP_LIMIT", 5),
			APIKeyLimit:   getEnvInt("API_KEY_LIMIT", 10),
			BlockDuration: getEnvInt("BLOCK_DURATION", 300),
			WindowSize:    getEnvInt("WINDOW_SIZE", 20),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

// getEnvInt obtém uma variável de ambiente e converte para int.
// Se a variável não existir ou não puder ser convertida, retorna o valor padrão
func getEnvInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnv obtém uma variável de ambiente ou retorna o valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type Config struct {
	Redis     RedisConfig
	RateLimit RateLimitConfig
	Server    ServerConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type RateLimitConfig struct {
	IPLimit       int
	APIKeyLimit   int
	BlockDuration int
	WindowSize    int
}

type ServerConfig struct {
	Port string
}
