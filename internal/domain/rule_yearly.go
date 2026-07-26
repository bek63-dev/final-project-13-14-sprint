// Package domain: файл rule_yearly.go реализует правило повторения "y" — ежегодно.
package domain

import "time"

// nextDateByYears сдвигает date на год вперёд, пока результат не окажется строго позже now.
func nextDateByYears(date, now time.Time) (time.Time, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if AfterNow(date, now) {
			return date, nil
		}
	}
}
