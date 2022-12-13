package heartbeat

import (
	"fmt"
	"math"
	"strings"
)

func secondsToReadableUkrainian(input int) string {
	weeks := math.Floor(float64(input) / 60 / 60 / 24 / 7)
	seconds := input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)

	var formatWeeks, formatDays, formatHours, formatMinutes string

	if weeks > 0 {
		formatWeeks = fmt.Sprintf("%.f %s, ", weeks, ukrainianWeeks(weeks))
	}
	if days > 0 {
		formatDays = fmt.Sprintf("%.f %s, ", days, ukrainianDays(days))
	}
	if hours > 0 {
		formatHours = fmt.Sprintf("%.f %s і ", hours, ukrainianHours(hours))
	}
	if minutes > 0 {
		formatMinutes = fmt.Sprintf("%.f %s", minutes, ukrainianMinutes(minutes))
	}

	return fmt.Sprintf("%s%s%s%s", formatWeeks, formatDays, formatHours, formatMinutes)
}

type ukrainianFormatForms struct {
	// 1, 21, 31..
	first string
	// 2, 3, 4, 22, 23, 24, 32..
	second string
	// 0, 5, 6, 7, 8, 9, 10..20, 25, 26..
	third string
}

// 1 				тиждень
// 2, 3, 4 	тижні
// 5..		 	тижнів
func ukrainianWeeks(weeks float64) string {
	return ukrainianFormat(weeks, ukrainianFormatForms{
		first:  "тиждень",
		second: "тижні",
		third:  "тижнів",
	})
}

// 1				день
// 2, 3, 4	дні
// 5, 6, 7 	днів
func ukrainianDays(weeks float64) string {
	return ukrainianFormat(weeks, ukrainianFormatForms{
		first:  "день",
		second: "дні",
		third:  "днів",
	})
}

// 1, 21..								годину
// 2, 3, 4, 22, 23.. 			години
// 5, 6, 7, 8, 9, 10..20 	годин
func ukrainianHours(weeks float64) string {
	return ukrainianFormat(weeks, ukrainianFormatForms{
		first:  "годину",
		second: "години",
		third:  "годин",
	})
}

// 1, 21..										хвилину
// 2, 3, 4, 22, 23.. 					хвилини
// 0, 5, 6, 7, 8, 9, 10..20 	хвилин
func ukrainianMinutes(weeks float64) string {
	return ukrainianFormat(weeks, ukrainianFormatForms{
		first:  "хвилину",
		second: "хвилини",
		third:  "хвилин",
	})
}

func ukrainianFormat(time float64, form ukrainianFormatForms) string {
	asString := fmt.Sprintf("%.0f", time)
	var lastNumber byte = asString[len(asString)-1]

	switch true {
	case lastNumber == '1':
		return form.first
	case strings.ContainsAny(string(lastNumber), "234") && (time < 5 || time >= 22):
		return form.second
	case strings.ContainsAny(string(lastNumber), "56789") || (time >= 10 && time <= 20) || lastNumber == '0':
		return form.third
	default:
		return form.third
	}
}
