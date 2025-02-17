package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/database"
	"go_final_project/taskstruct"
	"go_final_project/utilits"
	"net/http"
	"time"
)

func EditTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var task taskstruct.TaskObject
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	today := now.Format("20060102")
	if task.Date == "" || task.Date == "today" {
		task.Date = today
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(time.Now()) {
			if task.Repeat == "" {
				task.Repeat = today
			} else {
				nextDate, err := utilits.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Неверное правило повторения"}`, http.StatusBadRequest)
					return
				}
				if task.Date != today {
					task.Date = nextDate
				}
			}
		}
	}
	if task.Repeat != "" {
		_, err := utilits.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Неверное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}

	err := database.UpdateTask(task)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Ошибка обновления"}`, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}
