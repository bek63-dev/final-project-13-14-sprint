package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/domain"
)

// DateLayout — формат даты, используемый во всём API (алиас domain.DateLayout).
const DateLayout = domain.DateLayout

// Имена GET-параметров запроса к /api/nextdate.
const (
	nowParam    = "now"
	dateParam   = "date"
	repeatParam = "repeat"
)

// NextDateParams — разобранные параметры запроса /api/nextdate.
type NextDateParams struct {
	now    time.Time
	date   string
	repeat string
}

// parseNextDateParams читает и валидирует GET-параметры запроса /api/nextdate.
// Если параметр "now" не передан, используется текущая дата.
func parseNextDateParams(req *http.Request) (*NextDateParams, error) {
	now := domain.TruncateDate(time.Now())

	if nowStr := req.FormValue(nowParam); nowStr != "" {
		parsed, err := time.Parse(DateLayout, nowStr)
		if err != nil {
			return nil, fmt.Errorf("некорректный параметр %s=%q, ожидаемый формат %q",
				nowParam, nowStr, DateLayout)
		}
		now = parsed
	}

	return &NextDateParams{
		now:    now,
		date:   req.FormValue(dateParam),
		repeat: req.FormValue(repeatParam),
	}, nil
}

// NextDayHandler обрабатывает GET /api/nextdate:
// вычисляет дату следующего повторения задачи по переданным параметрам now, date и repeat и
// возвращает её в формате JSON.
func NextDayHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	params, err := parseNextDateParams(req)
	if err != nil {
		log.Println("NextDayHandler:", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := domain.NextDate(params.now, params.date, params.repeat)
	if err != nil {
		log.Println("NextDayHandler:", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(nextDate))
}
