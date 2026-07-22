// Package domain: файл date.go содержит базовые утилиты для работы с датами,
// общие для всех правил повторения
package domain

import "time"

// DateLayout — формат даты, используемый во всём проекте
const DateLayout = "20060102"

// TruncateDate обнуляет время суток (часы/минуты/секунды/наносекунды),
// оставляя только календарную дату. Нужна, чтобы сравнивать даты по дням,
// не обращая внимания на конкретное время внутри дня.
func TruncateDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0, t.Location())
}

// AfterNow сообщает, что календарный день date строго позже календарного дня now.
func AfterNow(date, now time.Time) bool {
	return TruncateDate(date).After(TruncateDate(now))
}

// lastDayOfMonth возвращает дату последнего календарного дня месяца, которому принадлежит date.
// Берем 1-е число следующего месяца и отнимаем один день.
func lastDayOfMonth(date time.Time) time.Time {
	firstOfNextMonth := time.Date(date.Year(), date.Month(),
		1, 0, 0, 0, 0,
		date.Location()).AddDate(0, 1, 0)

	return firstOfNextMonth.AddDate(0, 0, -1)
}
