package main

import (
	"go_final_project/database"
	"go_final_project/handlers"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {

	web := "./web"

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	database.FindEnv()
	database.OpenDB()
	defer database.CloseDB()

	r := chi.NewRouter()

	//Отображение сайта
	r.Handle("/", http.FileServer(http.Dir(web)))
	//Вывод новой даты
	r.Get("/api/nextdate", handlers.HandlerForNewDate)
	//Пост новой задачи
	r.Post("/api/task", handlers.PostTask)
	//Получение списка задач
	r.Get("/api/tasks", handlers.GetTasks)
	//Получение одной задачи
	r.Get("/api/task", handlers.GetTask)
	//Редактирование задачи
	r.Put("/api/task", handlers.EditTask)
	//Отметка о выполнении задачи
	r.Post("/api/task/done", handlers.DoneTask)
	//Удаление задачи
	r.Delete("/api/task", handlers.DeleteTask)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
