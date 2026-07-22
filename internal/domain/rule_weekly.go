// Package domain: файл rule_weekly.go реализует правило повторения "w <дни недели>"
package domain

import (
	"errors"
	"fmt"
	"time"
)

// Дни недели: 1 - понедельник, ..., 7 - воскресенье (ISO 8601, не как в time.Weekday)
const (
	minWeekday = 1
	maxWeekday = 7
)

// nextDateByWeekdays реализует правило "w <weekdays>":
// дата сдвигается день за днём, пока не попадёт на один из указанных дней недели
// и одновременно не окажется строго позже now.
func nextDateByWeekdays(date, now time.Time, ruleArgs []string) (time.Time, error) {
	if len(ruleArgs) == 0 || ruleArgs[0] == "" {
		return time.Time{},
			errors.New("nextDateByWeekdays: не указаны дни недели")
	}

	weekdays, err := parseRepeatArgs(ruleArgs[0])
	if err != nil {
		return time.Time{},
			fmt.Errorf("parseRepeatArgs: некорректный список дней недели: %w", err)
	}

	var allowedWeekday [maxWeekday + 1]bool

	for _, weekday := range weekdays {
		if weekday < minWeekday || weekday > maxWeekday {
			return time.Time{},
				fmt.Errorf("nextDateByWeekdays: недопустимый день недели %d, допустимо [%d, %d]",
					weekday, minWeekday, maxWeekday)
		}
		allowedWeekday[weekday] = true
	}

	for {
		date = date.AddDate(0, 0, 1)

		// time.Weekday нумерует с воскресенья (0), у нас же воскресенье - 7.
		// Поэтому 0 (воскресенье в time.Weekday) переводим в 7.
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		if allowedWeekday[weekday] && AfterNow(date, now) {
			return date, nil
		}
	}
}
