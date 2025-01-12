package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"go_final_project/taskstruct"
	"go_final_project/utilits"
	"net/http"
	"strconv"
	"time"
)

func PostTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var task = taskstruct.TaskObject{}
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
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
				task.Date = today
			} else {
				nextDate, err := utilits.NextDate(time.Now(), task.Date, task.Repeat)
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
	id, err := database.InsertInDB(task)
	idForResponse := strconv.Itoa(int(id))

	if err != nil {
		http.Error(w, `{"error":"Ошибка записи в БД"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{
		"id": idForResponse,
	})
}
