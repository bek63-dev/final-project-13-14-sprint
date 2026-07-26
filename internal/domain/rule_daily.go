// Package domain: файл rule_daily.go реализует правило повторения "d <интервал>" — через заданное число дней.
package domain

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Границы допустимого интервала в днях для правила "d".
const (
	minDayInterval = 1
	maxDayInterval = 400
)

// nextDateByDays сдвигает date на interval дней, пока результат не окажется строго позже now.
func nextDateByDays(date, now time.Time, ruleArgs []string) (time.Time, error) {
	if len(ruleArgs) == 0 || ruleArgs[0] == "" {
		return time.Time{},
			errors.New("nextDateByDays: не указан интервал дней")
	}

	interval, err := strconv.Atoi(ruleArgs[0])
	if err != nil {
		return time.Time{},
			fmt.Errorf("strconv.Atoi: некорректный интервал дней %q: %w", ruleArgs[0], err)
	}

	if interval < minDayInterval || interval > maxDayInterval {
		return time.Time{},
			fmt.Errorf("nextDateByDays: интервал дней %d вне допустимого диапазона [%d, %d]",
				interval, minDayInterval, maxDayInterval)
	}

	for {
		date = date.AddDate(0, 0, interval)
		if AfterNow(date, now) {
			return date, nil
		}
	}
}
