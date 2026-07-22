// Package domain: файл rule_yearly.go реализует правило повторения "y" — ежегодно.
package domain

import "time"

// nextDateByYears реализует правило "y": дата сдвигается на год вперёд, пока не окажется позже now
func nextDateByYears(date, now time.Time) (time.Time, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if AfterNow(date, now) {
			return date, nil
		}
	}
}
