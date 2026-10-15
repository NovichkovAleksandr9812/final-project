package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDayHandle(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		now = time.Now()
		_, err := NextDay(now, r.FormValue("date"), r.FormValue("repeat"))
		if err != nil {
			log.Printf("Ошибка %v", err)
			writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
			return
		}
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	repeat, err := NextDay(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	writeJson(w, repeat)
}

func afterDate(date, now time.Time) bool {
	return date.After(now)
}

func NextDay(now time.Time, dstart string, repeat string) (string, error) {
	var nextDate time.Time
	var interval int
	var rep []string

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("repeat пустая строка")
	}

	rep = strings.Split(repeat, " ")

	if len(rep) == 1 && rep[0] == "y" {
		interval = 1
	} else if len(rep) == 1 && rep[0] != "y" {
		return "", errors.New("недопустимый формат")
	} else if len(rep) == 2 && rep[0] == "d" {
		interval, err = strconv.Atoi(rep[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("недопустимый формат")
		}
	} else {
		return "", errors.New("недопустимый формат")
	}

	if rep[0] == "y" {
		var candidate time.Time
		yearsDiff := now.Year() - date.Year()
		if yearsDiff <= 0 {
			yearsDiff = 1
		}
		candidate = date.AddDate(yearsDiff, 0, 0)
		if afterDate(candidate, now) {
			nextDate = candidate
		} else {
			nextDate = candidate.AddDate(1, 0, 0)

		}
	} else {
		candidate := date.AddDate(0, 0, interval)

		for !afterDate(candidate, now) {
			candidate = candidate.AddDate(0, 0, interval)
		}

		nextDate = candidate
	}

	return nextDate.Format("20060102"), nil
}
