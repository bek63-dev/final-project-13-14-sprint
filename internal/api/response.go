// Package api: файл response.go содержит вспомогательные функции и структуры
// для формирования JSON-ответов API в едином формате.
package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse — формат ответа при ошибке: {"error": "текст ошибки"}
type errorResponse struct {
	Error string `json:"error"`
}

// idResponse — формат ответа при успешном добавлении задачи: {"id": "123"}
type idResponse struct {
	ID string `json:"id"`
}

// writeJSON сериализует data в JSON и отправляет клиенту с указанным статус-кодом.
// Используется и для успешных, и для ошибочных ответов.
func writeJSON(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(data); err != nil {
		log.Println("writeJSON: json.Encode:", err)
	}
}

// writeError отправляет JSON-ответ с ошибкой в формате {"error": "..."}.
func writeError(res http.ResponseWriter, status int, message string) {
	writeJSON(res, status, errorResponse{Error: message})
}

// isMethodAllowed проверяет, соответствует ли метод запроса ожидаемому.
// Если нет — отправляет ошибку 405 Method Not Allowed и возвращает false.
func isMethodAllowed(expectedMethod string, res http.ResponseWriter, req *http.Request) bool {
	if req.Method != expectedMethod {
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
		return false
	}
	return true
}
