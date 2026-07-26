// Package domain: файл date.go содержит базовые утилиты работы с датами,
// общие для всех правил повторения и для валидации на уровне api.
package domain

import "time"

// DateLayout — формат хранения и передачи даты (YYYYMMDD), используемый на сервере, в БД и в JSON API
const DateLayout = "20060102"

// searchDateLayout — формат даты (DD.MM.YYYY), в котором пользователь вводит дату при поиске задач.
const searchDateLayout = "02.01.2006"

// TruncateDate обнуляет время суток, оставляя только календарную дату.
// Нужна, чтобы сравнивать даты по дням, не учитывая время внутри дня.
func TruncateDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0, t.Location())
}

// AfterNow сообщает, что календарный день date строго позже календарного дня now.
// Используется всеми правилами повторения как условие остановки перебора дат.
func AfterNow(date, now time.Time) bool {
	return TruncateDate(date).After(TruncateDate(now))
}

// lastDayOfMonth возвращает дату последнего календарного дня месяца, которому принадлежит date.
func lastDayOfMonth(date time.Time) time.Time {
	firstOfNextMonth := time.Date(date.Year(), date.Month(),
		1, 0, 0, 0, 0,
		date.Location()).AddDate(0, 1, 0)

	return firstOfNextMonth.AddDate(0, 0, -1)
}

// ParseSearchDate пытается разобрать search как дату в формате
// searchDateLayout и вернуть её в формате хранения DateLayout.
// Если search не похож на дату, возвращает false — тогда поиск задач
// выполняется по тексту, а не по дате.
func ParseSearchDate(search string) (string, bool) {
	t, err := time.Parse(searchDateLayout, search)
	if err != nil {
		return "", false
	}
	return t.Format(DateLayout), true
}
