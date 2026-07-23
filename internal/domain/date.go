// Package domain: файл date.go содержит базовые утилиты для работы с датами,
// общие для всех правил повторения
package domain

import "time"

const DateLayout = "20060102"         // стандартный формат даты (YYYYMMDD), используемый на сервере и в БД
const searchDateLayout = "02.01.2006" // формат даты (DD.MM.YYYY), используемый при поиске пользовательским вводом

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

// ParseSearchDate пытается разобрать строку search как дату в формате 02.01.2006
// и вернуть её в формате хранения 20060102. Если search не похож на дату - возвращает false.
func ParseSearchDate(search string) (string, bool) {
	t, err := time.Parse(searchDateLayout, search)
	if err != nil {
		return "", false
	}
	return t.Format(DateLayout), true
}
