// Package server оборачивает стандартный net/http.Server, задавая
// безопасные тайм-ауты и логирование для HTTP-сервера планировщика задач.
package server

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// Server — HTTP-сервер приложения с предварительно настроенными тайм-аутами и логером ошибок.
type Server struct {
	server *http.Server
	logger *log.Logger
}

// Тайм-ауты по умолчанию для сетевых соединений сервера.
const (
	defaultReadTimeout       = 5 * time.Second  // максимум на чтение всего запроса
	defaultReadHeaderTimeout = 2 * time.Second  // максимум на чтение заголовков
	defaultWriteTimeout      = 10 * time.Second // максимум на запись ответа
	defaultIdleTimeout       = 30 * time.Second // максимум простоя keep-alive соединения
)

// New создаёт Server с указанными адресом, логером и обработчиком запросов,
// применяя тайм-ауты по умолчанию. Вызывается в main перед запуском сервера.
func New(address string, logger *log.Logger, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:              address,
			Handler:           handler,
			ErrorLog:          logger,
			ReadTimeout:       defaultReadTimeout,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			WriteTimeout:      defaultWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
		},
		logger: logger,
	}
}

// Run запускает прослушивание сетевого адреса и обработку входящих HTTP-запросов.
// Метод блокирует выполнение до остановки сервера или возникновения критической ошибки.
// Возвращает nil, если сервер был штатно остановлен (http.ErrServerClosed).
func (s *Server) Run() error {
	s.logger.Printf("\nСервер запущен на %s", s.server.Addr)

	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
