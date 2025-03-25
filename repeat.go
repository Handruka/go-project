package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) { // функция возвращает следующую дату в соответствии с правилом "repeat"
	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("ошибка парсинга dateTime в формат time.time ")
	}
	lng := len(repeat)
	str := strings.Split(repeat, " ")

	switch {
	case lng == 1 && str[0] == "y":
		dateTime = dateTime.AddDate(1, 0, 0)
		for now.After(dateTime) {
			dateTime = dateTime.AddDate(1, 0, 0)
		}
		return dateTime.Format("20060102"), nil
	case lng > 2 && lng < 6 && str[0] == "d":
		d, err := strconv.Atoi(str[1])
		if err != nil || d > 400 {
			return "", errors.New("ошибка конвертации dayString в число")
		}
		if d < 401 {
			dateTime = dateTime.AddDate(0, 0, d)
			for now.After(dateTime) {
				dateTime = dateTime.AddDate(0, 0, d)
			}
		} else {
			return "", errors.New("repeat больше 400")

		}
	default:
		return "", errors.New("ошибка передаваемых значений в функцию nextDate")
	}
	return dateTime.Format("20060102"), nil
}
