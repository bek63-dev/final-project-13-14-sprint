// Package api: файл response.go содержит общие функции формирования JSON-ответов
// в едином формате для всех обработчиков пакета api.
package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse — формат ответа при ошибке: {"error": "текст ошибки"}.
type errorResponse struct {
	Error string `json:"error"`
}

// idResponse — формат ответа при успешном добавлении задачи: {"id": "123"}.
type idResponse struct {
	ID string `json:"id"`
}

// writeJSON сериализует data в JSON и отправляет клиенту с заданным статус-кодом.
// Используется всеми обработчиками для успешных и ошибочных ответов.
func writeJSON(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(data); err != nil {
		log.Println("writeJSON: json.Encode:", err)
	}
}

// writeError отправляет JSON-ответ с ошибкой в формате {"error": "..."}.
// Используется обработчиками при некорректном запросе или сбое обработки.
func writeError(res http.ResponseWriter, status int, message string) {
	writeJSON(res, status, errorResponse{Error: message})
}
