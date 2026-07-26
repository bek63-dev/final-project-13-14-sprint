// Package config загружает конфигурацию сервера — хост, порт, путь к файлу БД,
// пароль и секретный ключ аутентификации — из .env файла и переменных окружения,
// применяя значения по умолчанию там, где это уместно.
package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config хранит итоговые параметры запуска сервера, полученные из .env и окружения.
// Используется при инициализации БД, маршрутизатора и HTTP-сервера в main.
type Config struct {
	Host      string
	Port      string
	DBFile    string
	Password  string
	SecretKey string
}

// Значения по умолчанию, применяемые, если соответствующая переменная
// окружения не задана. У Password намеренно нет значения по умолчанию:
// пустой пароль осмысленно означает отключённую аутентификацию.
const (
	defaultHost   = ""
	defaultPort   = "7540"
	defaultDBFile = "scheduler.db"
)

// Имена переменных окружения, из которых читается конфигурация.
const (
	envTodoHost      = "TODO_HOST"
	envTodoPort      = "TODO_PORT"
	envTodoDBFile    = "TODO_DBFILE"
	envTodoPassword  = "TODO_PASSWORD"
	envTodoSecretKey = "TODO_SECRET_KEY"
)

// Load читает .env файл и переменные окружения и возвращает итоговую
// конфигурацию сервера. Вызывается один раз при старте приложения в main.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("godotenv.Load: %v", err)
	}

	cfg := &Config{
		Host:   defaultHost,
		Port:   defaultPort,
		DBFile: defaultDBFile,
	}

	// Каждая переменная окружения перекрывает дефолт, только если она непустая
	if host := os.Getenv(envTodoHost); host != "" {
		cfg.Host = host
	}

	if port := os.Getenv(envTodoPort); port != "" {
		cfg.Port = port
	}

	if dbFile := os.Getenv(envTodoDBFile); dbFile != "" {
		cfg.DBFile = dbFile
	}

	// Пароль читается как есть, без дефолта: пустая строка означает отключённую аутентификацию
	cfg.Password = os.Getenv(envTodoPassword)

	cfg.SecretKey = os.Getenv(envTodoSecretKey)

	cfg.DBFile = filepath.Clean(strings.TrimSpace(cfg.DBFile))

	return cfg, nil
}

// Address собирает Host и Port в формате host:port, который принимает http.Server.Addr
func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
