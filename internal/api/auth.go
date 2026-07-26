// Package api: файл auth.go содержит обработчик входа в систему и
// middleware проверки JWT-токена для защищённых маршрутов API.
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/auth"
)

// tokenCookieName — имя куки, в которой фронтенд хранит JWT-токен.
const tokenCookieName = "token"

// signInRequest описывает тело запроса POST /api/signin.
type signInRequest struct {
	Password string `json:"password"`
}

// tokenResponse — формат успешного ответа аутентификации: {"token": "..."}.
type tokenResponse struct {
	Token string `json:"token"`
}

// SignInHandler обрабатывает POST /api/signin: сверяет пароль из запроса
// с паролем из конфигурации (TODO_PASSWORD из .env) и при совпадении выдаёт JWT-токен для
// последующих запросов к защищённым маршрутам.
func (h *Handler) SignInHandler(res http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed(http.MethodPost, res, req) {
		return
	}

	var signIn signInRequest
	if err := json.NewDecoder(req.Body).Decode(&signIn); err != nil {
		log.Println("SignInHandler: json.NewDecoder:", err)
		writeError(res, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}
	defer req.Body.Close()

	if signIn.Password != h.Config.Password {
		writeError(res, http.StatusBadRequest, "Неверный пароль")
		return
	}

	token, err := auth.GenerateToken(signIn.Password, h.Config.SecretKey)
	if err != nil {
		log.Println("SignInHandler: auth.GenerateToken:", err)
		writeError(res, http.StatusInternalServerError, "ошибка генерации токена")
		return
	}

	writeJSON(res, http.StatusOK, tokenResponse{Token: token})
}

// Auth — middleware, защищающее обработчик next проверкой JWT-токена из
// куки "token". Если пароль в конфигурации не задан (TODO_PASSWORD из .env),
// аутентификация считается отключённой и запрос пропускается без проверки.
// Используется маршрутизатором для оборачивания защищённых эндпоинтов задач.
func (h *Handler) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if h.Config.Password == "" {
			next(res, req)
			return
		}

		cookie, err := req.Cookie(tokenCookieName)
		if err != nil || !auth.ValidateToken(cookie.Value, h.Config.Password, h.Config.SecretKey) {
			http.Error(res, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(res, req)
	}
}
