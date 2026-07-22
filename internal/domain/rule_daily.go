// Package domain: файл rule_daily.go реализует правило повторения "d <интервал>" - через N дней.
package domain

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Границы допустимого интервала в днях
const (
	minDayInterval = 1
	maxDayInterval = 400
)

// nextDateByDays реализует правило "d <interval>":
// дата сдвигается на interval дней, пока не окажется позже now
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

	// Сдвигаем дату шагами по interval дней, пока не пройдём now.
	for {
		date = date.AddDate(0, 0, interval)
		if AfterNow(date, now) {
			return date, nil
		}
	}
}
