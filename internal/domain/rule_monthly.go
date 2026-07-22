// Package domain: файл rule_monthly.go реализует правило повторения "m <дни месяца> [месяцы]"
package domain

import (
	"errors"
	"fmt"
	"time"
)

const (
	minMonthDay = 1
	maxMonthDay = 31

	lastMonthDay       = -1 // означает "последний день месяца"
	secondLastMonthDay = -2 // означает "предпоследний день месяца"

	minMonth = 1
	maxMonth = 12
)

// nextDateByMonth реализует правило "m <days> [months]": дата сдвигается день за днём,
// пока не попадёт на один из указанных дней месяца (1..31, либо -1 — последний день месяца, -2 — предпоследний), и,
// если указан список месяцев, месяц даты входит в этот список, и сама дата строго позже now.
func nextDateByMonth(date, now time.Time, ruleArgs []string) (time.Time, error) {
	if len(ruleArgs) == 0 || ruleArgs[0] == "" {
		return time.Time{},
			errors.New("nextDateByMonth: не указаны дни месяца")
	}

	days, err := parseRepeatArgs(ruleArgs[0])
	if err != nil {
		return time.Time{},
			fmt.Errorf("parseRepeatArgs: некорректный список дней месяца: %w", err)
	}

	var allowedDay [maxMonthDay + 1]bool
	wantLastDay, wantSecondLastDay := false, false
	for _, day := range days {
		switch day {
		case lastMonthDay:
			wantLastDay = true
		case secondLastMonthDay:
			wantSecondLastDay = true
		default:
			if day < minMonthDay || day > maxMonthDay {
				return time.Time{},
					fmt.Errorf("недопустимый день месяца %d, допустимо [%d, %d], %d или %d",
						day, minMonthDay, maxMonthDay, lastMonthDay, secondLastMonthDay)
			}
			allowedDay[day] = true
		}
	}

	var allowedMonth [maxMonth + 1]bool
	anyMonth := true
	if len(ruleArgs) > 1 && ruleArgs[1] != "" {
		months, err := parseRepeatArgs(ruleArgs[1])
		if err != nil {
			return time.Time{},
				fmt.Errorf("некорректный список месяцев: %w", err)
		}

		anyMonth = false
		for _, month := range months {
			if month < minMonth || month > maxMonth {
				return time.Time{},
					fmt.Errorf("недопустимый месяц %d, допустимо [%d, %d]",
						month, minMonth, maxMonth)
			}
			allowedMonth[month] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)

		if !anyMonth && !allowedMonth[int(date.Month())] {
			continue
		}

		day := date.Day()
		daysInMonth := lastDayOfMonth(date).Day()

		dayMatches := allowedDay[day] || (wantLastDay && day == daysInMonth) ||
			(wantSecondLastDay && day == daysInMonth-1)

		if dayMatches && AfterNow(date, now) {
			return date, nil
		}
	}
}
