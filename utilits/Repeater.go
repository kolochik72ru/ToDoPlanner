package utilits

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Функция ежедневного повтора
func DailyRepeat(now time.Time, timeDate time.Time, format []string) (string, error) {
	if len(format) != 2 {
		return "", errors.New("неверный формат повторения")
	}

	days, err := strconv.Atoi(format[1])
	if err != nil || days <= 0 || days > 400 {
		return "", errors.New("неверный формат интервала")
	}

	if timeDate.After(now) {
		timeDate = timeDate.AddDate(0, 0, days)
	} else {
		for !timeDate.After(now) {
			timeDate = timeDate.AddDate(0, 0, days)
		}
	}
	return timeDate.Format("20060102"), nil
}

// Парсинг дней недели
func parseWeekdays(days string) ([]time.Weekday, error) {
	dayStrings := strings.Split(days, ",")
	var weekdays []time.Weekday

	for _, dayStr := range dayStrings {
		dayInt, err := strconv.Atoi(dayStr)
		if err != nil || dayInt < 1 || dayInt > 7 {
			return nil, errors.New("неверный формат дней недели")
		}
		weekday := time.Weekday((dayInt % 7))
		weekdays = append(weekdays, weekday)
	}
	return weekdays, nil
}

// Функция проверки дня
func WeeklyRepeat(now time.Time, timeDate time.Time, format []string) (string, error) {
	if len(format) != 2 {
		return "", errors.New("неверный формат повтора недели")
	}
	weekdays, err := parseWeekdays(format[1])
	if err != nil {
		return "", err
	}
	for {
		for _, weekday := range weekdays {
			if timeDate.Weekday() == weekday {
				if timeDate.After(now) {
					return timeDate.Format("20060102"), nil
				}
			}
		}
		timeDate = timeDate.AddDate(0, 0, 1)
	}
}

// Парсинг дней
func parseMonthDays(days string) ([]int, error) {
	daysParts := strings.Split(days, ",")
	var dayInt []int
	for _, part := range daysParts {
		day, err := strconv.Atoi(part)
		if err != nil || day == 0 || day < -31 || day > 31 {
			return nil, errors.New("неверный день месяца")
		}
		dayInt = append(dayInt, day)
	}
	return dayInt, nil
}

// Разделение месяца
func parseMonth(monthStr string) ([]int, error) {
	monthsParts := strings.Split(monthStr, ",")
	var months []int
	for _, part := range monthsParts {
		month, err := strconv.Atoi(part)
		if err != nil || month < 1 || month > 12 {
			return nil, errors.New("неверный месяц")
		}
		months = append(months, month)
	}
	return months, nil
}

// Проверка правильности месяца
func isMonthRight(currentMonth int, months []int) bool {
	for _, month := range months {
		if currentMonth == month {
			return true
		}
	}
	return false
}

// Вычисление даты
func calculateDate(now, timeDate time.Time, days []int, allowedMonths []int) time.Time {
	year, month, _ := timeDate.Date()
	location := timeDate.Location()
	sort.Ints(days)
	var nearestDate *time.Time
	for {
		if len(allowedMonths) > 0 && !isMonthRight(int(month), allowedMonths) {
			month++
			if month > 12 {
				month = 1
				year++
			}
			continue
		}
		for _, day := range days {
			var targetDate time.Time
			if day < 0 {
				firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, location)
				lastDayThisMonth := firstDayNextMonth.AddDate(0, 0, -1)
				targetDate = lastDayThisMonth.AddDate(0, 0, day+1)
			} else {
				targetDate = time.Date(year, month, day, 0, 0, 0, 0, location)
				if targetDate.Month() != time.Month(month) {
					continue
				}
			}
			if targetDate.After(now) {
				if nearestDate == nil || targetDate.Before(*nearestDate) {
					nearestDate = &targetDate
				}
			}
		}
		if nearestDate != nil {
			return *nearestDate
		}
		month++
		if month > 12 {
			month = 1
			year++
		}
	}
}

// Главная обработка
func MonthRepeat(now time.Time, timeDate time.Time, format []string) (string, error) {
	if len(format) < 2 || len(format) > 3 {
		return "", errors.New("неверный формат месяца")
	}
	days, err := parseMonthDays(format[1])
	if err != nil {
		return "", err
	}
	var months []int
	if len(format) == 3 {
		months, err = parseMonth(format[2])
		if err != nil {
			return "", err
		}
	}
	for _, day := range days {
		if day < -2 || day == 0 || day > 31 {
			return "", errors.New("неверный день")
		}
	}
	for {
		if len(months) > 0 && !isMonthRight(int(timeDate.Month()), months) {
			timeDate = timeDate.AddDate(0, 1, 0)
			continue
		}
		targetDate := calculateDate(now, timeDate, days, months)
		if targetDate.After(now) {
			return targetDate.Format("20060102"), nil
		}
		timeDate = timeDate.AddDate(0, 1, 0)
	}
}

// Функция годового повтора
func YearRepeat(now time.Time, timeDate time.Time, format []string) (string, error) {
	if len(format) != 1 {
		return "", errors.New("неверный формат года")
	}
	if timeDate.After(now) {
		timeDate = timeDate.AddDate(1, 0, 0)
	} else {
		for !timeDate.After(now) {
			timeDate = timeDate.AddDate(1, 0, 0)
		}
	}
	return timeDate.Format("20060102"), nil
}

// Функция вычисляет следующую дату задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {
	timeDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("некорректный формат даты")
	}
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}
	format := strings.Split(repeat, " ")
	switch format[0] {
	case "d":
		return DailyRepeat(now, timeDate, format)
	case "y":
		return YearRepeat(now, timeDate, format)
	case "w":
		return WeeklyRepeat(now, timeDate, format)
	case "m":
		return MonthRepeat(now, timeDate, format)
	default:
		return "", errors.New("правило повторения не поддерживается")
	}
}
