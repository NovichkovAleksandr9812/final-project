package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", NextDayHandle)
	mux.HandleFunc("/api/task", taskHandle)
	mux.HandleFunc("/api/tasks", tasksHandle)
	mux.HandleFunc("/api/task/done", taskDoneHandle)
}

func taskHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandle(w, r)
	case http.MethodGet:
		getTaskHandle(w, r)
	case http.MethodPut:
		putTaskHandle(w, r)
	case http.MethodDelete:
		deleteTaskHandle(w, r)
	default:
		http.Error(w, r.Method+" не поддерживается", http.StatusMethodNotAllowed)
	}
}
