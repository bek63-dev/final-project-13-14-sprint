// Package domain: файл repeat_rule.go содержит парсинг строки repeat:
// разбор на тип правила + аргументы, а также разбор списков чисел через запятую.
package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Типы правил повторения — первое "слово" строки repeat
const (
	RuleTypeDays   = "d"
	RuleTypeWeeks  = "w"
	RuleTypeMonths = "m"
	RuleTypeYears  = "y"
)

// parseRepeatRule разбивает строку repeat на тип правила ("d", "w", "m", "y")
// и оставшиеся аргументы. Например, "m 1,15 3,6" -> ("m", ["1,15", "3,6"]).
func parseRepeatRule(repeat string) (string, []string, error) {
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", nil,
			errors.New("parseRepeatRule: правило повторения не задано")
	}

	return parts[0], parts[1:], nil
}

// parseRepeatArgs парсит список целых чисел через запятую, например "1,15,25" -> [1,15,25].
func parseRepeatArgs(ruleArgs string) ([]int, error) {
	parts := strings.Split(ruleArgs, ",")
	values := make([]int, 0, len(parts))

	for _, part := range parts {
		value, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil,
				fmt.Errorf("parseRepeatArgs: некорректное число %q: %w", part, err)
		}

		values = append(values, value)
	}

	return values, nil
}
