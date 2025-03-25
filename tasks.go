package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func Tasks(res http.ResponseWriter, req *http.Request) { // функция возвращает последние 50 созданных задач в JSON-формате
	switch req.Method {
	case http.MethodGet:
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		db, err := sql.Open("sqlite", "./scheduler.db")
		if err != nil {
			http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
			return
		}
		defer db.Close()
		var tasks Planner
		query := "SELECT * FROM scheduler ORDER BY date ASC LIMIT 50"
		rows, err := db.Query(query)
		if err != nil {
			http.Error(res, `{"error":"ошибка при чтении rows базы данных"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var planner Tsk
			var dateInt int
			err = rows.Scan(&planner.Id, &dateInt, &planner.Title, &planner.Comment, &planner.Repeat)
			if err != nil {
				http.Error(res, `{"error":"ошибка при записи данных в переменные, при чтении rows базы данных"}`, http.StatusInternalServerError)
				return
			}
			planner.Date = fmt.Sprintf("%d", dateInt)
			tasks.Tasks = append(tasks.Tasks, planner)
		}

		if len(tasks.Tasks) == 0 {
			// Возвращаем пустой список задач
			emptyResponse := Planner{Tasks: []Tsk{}}
			jsonResponse, err := json.Marshal(emptyResponse)
			if err != nil {
				http.Error(res, `{"error":"ошибка при создании JSON-ответа"}`, http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusOK)
			res.Write(jsonResponse)
			return
		}

		response := Planner{Tasks: tasks.Tasks}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(res, `{"error":"ошибка при создании JSON-ответа"}`, http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write(jsonResponse)
	}
}
