// Package domain: файл next_date.go содержит основную публичную функцию пакета domain -
// NextDate, которая по правилу repeat вычисляет следующую дату выполнения задачи.
package domain

import (
	"fmt"
	"time"
)

// NextDate возвращает следующую дату выполнения задачи в формате DateLayout.
// Отсчёт всегда ведётся от startDate: правило применяется к startDate до тех пор,
// пока результат не станет строго больше now (даже если startDate уже сам по себе больше now)
//
// Параметры:
//   - now       — дата, от которой ищется ближайшая следующая дата;
//   - startDate — исходная дата задачи в формате DateLayout (20060102);
//   - repeat    — правило повторения ("d 3", "y", "w 1,2", "m -1,5 1,6" и т.д.).
func NextDate(now time.Time, startDate string, repeat string) (string, error) {
	ruleType, ruleArgs, err := parseRepeatRule(repeat)
	if err != nil {
		return "", err
	}

	date, err := time.Parse(DateLayout, startDate)
	if err != nil {
		return "",
			fmt.Errorf("time.Parse: некорректный формат исходной даты %q, ожидаемый формат %q: %w",
				startDate, DateLayout, err)
	}

	// Диспетчеризация по типу правила
	switch ruleType {
	case RuleTypeDays:
		date, err = nextDateByDays(date, now, ruleArgs)
	case RuleTypeWeeks:
		date, err = nextDateByWeekdays(date, now, ruleArgs)
	case RuleTypeMonths:
		date, err = nextDateByMonth(date, now, ruleArgs)
	case RuleTypeYears:
		date, err = nextDateByYears(date, now)
	default:
		err = fmt.Errorf("неподдерживаемый формат правила повторения: %q", ruleType)
	}

	if err != nil {
		return "", err
	}

	return date.Format(DateLayout), nil
}
