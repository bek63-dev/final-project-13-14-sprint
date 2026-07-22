// Package config отвечает за загрузку конфигурации сервера:
// хост, порт и путь к файлу БД. Значения берутся из .env файла,
// при отсутствии — используются значения по умолчанию.
package config

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config хранит итоговые параметры запуска сервера
type Config struct {
	Host   string
	Port   string
	DBFile string
}

// Значения по умолчанию, если в .env файле не заданы значения переменных окружения
const (
	defaultHost   = ""
	defaultPort   = "7540"
	defaultDBFile = "scheduler.db"
)

// Имена переменных окружения
const (
	envTodoHost   = "TODO_HOST"
	envTodoPort   = "TODO_PORT"
	envTodoDBFile = "TODO_DBFILE"
)

// Load читает .env файл (если он есть) и переменные окружения,
// возвращает итоговую конфигурацию сервера
func Load() Config {
	// Если .env файла нет - это не критическая ошибка, просто логируем и идём дальше
	if err := godotenv.Load(); err != nil {
		log.Printf("godotenv.Load: %v", err)
	}

	// Стартуем со значениями по умолчанию
	cfg := Config{
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

	cfg.DBFile = filepath.Clean(strings.TrimSpace(cfg.DBFile))

	return cfg
}

// Address собирает host:port в формате, который принимает http.Server.Addr
func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
