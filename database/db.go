package database

import (
	"database/sql"
	"errors"
	"go_final_project/taskstruct"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func FindEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env файл не найден или не работает: %v", err)
	}
}

func OpenDB() {

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	log.Print("База данных запущена.")

	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	if install {
		err = CreateDB(DB)
		if err != nil {
			log.Fatalf("Ошибка создания базы данных: %v", err)
		}
		log.Println("База данных создана.")
	}

}

func CreateDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK (LENGTH(repeat) <= 128)
	);

	CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	_, err := db.Exec(query)
	return err
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
		log.Println("Соединение закрыто.")
	}
}

func InsertInDB(task taskstruct.TaskObject) (uint64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func FindInDb(search string, limit int) ([]taskstruct.TaskObject, error) {
	var query string
	var args []interface{}

	if search == "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
		args = append(args, limit)
	} else {
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?"
			args = append(args, parsedDate.Format("20060102"), limit)
		} else {
			likePattern := "%" + search + "%"
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			args = append(args, likePattern, likePattern, limit)
		}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, errors.New("ошибка выполнения запроса к базе данных")
	}
	defer rows.Close()

	var tasks []taskstruct.TaskObject
	for rows.Next() {
		var task taskstruct.TaskObject
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, errors.New("ошибка сканирования данных из базы")
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("ошибка обработки данных из базы")
	}
	return tasks, nil
}

func GetTaskByID(id int) (taskstruct.TaskObject, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := DB.QueryRow(query, id)
	var task taskstruct.TaskObject
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return taskstruct.TaskObject{}, errors.New("задача не найдена")
	}
	return task, nil
}
func UpdateTask(task taskstruct.TaskObject) error {
	query := `
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteTaskByID(id int) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := DB.Exec(query, id)
	return err
}

func UpdateTaskDate(id uint64, newDate string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	_, err := DB.Exec(query, newDate, id)
	return err
}
