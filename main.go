package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var DBFile string

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
			http.Error(res, "ошибка конвертации формата nowParse", http.StatusBadRequest)
			return
		}
		out, err = NextDate(nowParse, date, repeat)
		if err != nil {
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
	appPath, err := os.Getwd() //Получаем путь к файлу main.go
	if err != nil {
		log.Fatal("Функция os.Getwd() выполнилась с ошибкой ", err)
	}
	DBFile = filepath.Join(appPath, "scheduler.db") // Создаём путь к файлу БД
	_, err = os.Stat(DBFile)
	var Install bool
	if err != nil {
		Install = true
	}
	CreatDB(Install) // создаём базу данных create.go

	http.HandleFunc("/api/", mainHandle)

	WebDir := "./web"
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(WebDir))))
	err = http.ListenAndServe("localhost:7540", nil)
	if err != nil {
		panic(err)
	}
}
