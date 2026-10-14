package api

import (
	"final-project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

type ErrorResp struct {
	Error string `json:"error"`
}

func tasksHandle(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}
	writeJson(w, TasksResp{Tasks: tasks})
}

func getTaskHandle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор (параметр id)")
		return
	}

	task, err := db.GetTask(id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	if task == nil {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	writeJson(w, task)
}

func taskDoneHandle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор (параметр id)")
		return
	}

	task, err := db.GetTask(id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	if task == nil {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	if task.Repeat == "" {
		if err = db.DeleteTask(id); err != nil {
			writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
			return
		}
		writeJson(w, map[string]interface{}{})
		return
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	nextDateStr, err := NextDay(date, task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	if err = db.UpdateDate(nextDateStr, id); err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	writeJson(w, map[string]interface{}{})
}

func deleteTaskHandle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор (параметр id)")
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID должен быть числом")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка")
		return
	}

	writeJson(w, map[string]interface{}{})
}
