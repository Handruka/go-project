package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

var DBFile string
var db *sql.DB

type Tsk struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Planner struct {
	Tasks []Tsk `json:"tasks"`
}

func mainHandle(res http.ResponseWriter, req *http.Request) {
	var out string

	switch {
	case req.URL.Path == "/api/nextdate":
		now := req.URL.Query().Get("now")
		date := req.URL.Query().Get("date")
		repeat := req.URL.Query().Get("repeat")
		nowParse, err := time.Parse("20060102", now)
		if err != nil {
			log.Printf("ошибка конвертации формата nowParse: %v", err)
			http.Error(res, "ошибка конвертации формата nowParse", http.StatusBadRequest)
			return
		}
		out, err = NextDate(nowParse, date, repeat)
		if err != nil {
			log.Printf("ошибка в NextDate: %v", err)
			http.Error(res, "ошибка в NextDate: "+err.Error(), http.StatusBadRequest)
			return
		}
		res.Write([]byte(out))
	case req.URL.Path == "/api/task":
		if req.Method == http.MethodPost {
			AddTask(res, req)
		}
		if req.Method == http.MethodGet {
			Task(res, req)
		}
		if req.Method == http.MethodPut {
			EditTask(res, req)
		}
		if req.Method == http.MethodDelete {
			Delete(res, req)
		}
	case req.URL.Path == "/api/tasks":
		Tasks(res, req)
	case req.URL.Path == "/api/task/done":
		Done(res, req)
	}

}

func main() {

	var install bool
	err := godotenv.Load() // загружаем переменные из .env
	if err != nil {
		log.Fatal("ошибка загрузки .env файла")
	}
	port := os.Getenv("TODO_PORT") // считываем переменную из .env и присваиваем её переменной port
	if port == "" {
		port = "7540"
	}
	appPath, err := os.Getwd() //Получаем путь к файлу main.go
	if err != nil {
		log.Printf("Функция os.Getwd() выполнилась с ошибкой: %v", err)
	}
	DBFile = filepath.Join(appPath, "scheduler.db") // Создаём путь к файлу БД
	_, err = os.Stat(DBFile)

	if err != nil {
		install = true
	}
	// если install равен true(файл базы данных не существует)

	db, err = CreatDB(install) // создаём базу данных create.go
	if err != nil {
		log.Printf("Функция CreatDB выполнилась с ошибкой: %v", err)
	}

	defer db.Close()

	http.HandleFunc("/api/", mainHandle)

	WebDir := "./web"
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(WebDir))))
	log.Printf("Сервер запущен на порту:%s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
