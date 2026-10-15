package api

import (
	"encoding/json"
	"final-project/pkg/db"
	"fmt"
	"log"
	"net/http"
	"time"
)

func addTaskHandle(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if task == (db.Task{}) {
		writeJson(w, map[string]string{"error": "Пустое тело запроса"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "поле Title обязательно для заполнения"})
		return
	}

	if err := checkDate(&task); err != nil {
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	writeJson(w, response)
}

func putTaskHandle(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if task == (db.Task{}) {
		writeJson(w, map[string]string{"error": "Пустое тело запроса"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "поле Title обязательно для заполнения"})
		return
	}

	if err := checkDate(&task); err != nil {
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		log.Printf("Ошибка %v", err)
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	respUp := map[string]interface{}{}
	writeJson(w, respUp)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStr := today.Format("20060102")

	if task.Date == "" {
		task.Date = todayStr
	}

	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}

	if task.Repeat == "" {
		task.Date = todayStr
		return nil
	}

	if parsedDate.After(today) {
		task.Date = todayStr
		return nil
	}

	nextDateStr, err := NextDay(today, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	parsedNext, err := time.Parse("20060102", nextDateStr)
	if err != nil {
		return fmt.Errorf("ошибка формата даты от NextDay: %w", err)
	}

	if !parsedNext.After(today) {
		return fmt.Errorf("вычисленная дата (%s) не больше сегодняшней (%s)", nextDateStr, todayStr)
	}

	task.Date = nextDateStr

	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		log.Printf("[CRITICAL] Ошибка записи в сокет: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)

	resp := ErrorResp{
		Error: message,
	}

	writeJson(w, resp)
}
