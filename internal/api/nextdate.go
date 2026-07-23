package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/domain"
)

const DateLayout = domain.DateLayout

const (
	nowParam    = "now"
	dateParam   = "date"
	repeatParam = "repeat"
)

type NextDateParams struct {
	now    time.Time
	date   string
	repeat string
}

// parseNextDateParams парсит и валидирует GET-параметры звпроса для эндпоинта /api/nextdate
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

func NextDayHandler(res http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed(http.MethodGet, res, req) {
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

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, nextDate)
}
