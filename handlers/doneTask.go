package handlers

import (
	"database/sql"
	"go_final_project/database"
	"go_final_project/utilits"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func DoneTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}

	if idStr == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"Неверный идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Ошибка получения задачи из БД"}`, http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		err := database.DeleteTaskByID(id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := utilits.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка расчёта следующей даты повторения"}`, http.StatusInternalServerError)
			return
		}

		err = database.UpdateTaskDate(uint64(id), nextDate)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}
