package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Tasks(res http.ResponseWriter, req *http.Request) { // функция возвращает последние 50 созданных задач в JSON-формате

	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tasks Planner
	query := "SELECT * FROM scheduler ORDER BY date ASC LIMIT 50"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("ошибка при чтении rows базы данных: %v", err)
		http.Error(res, `{"error":"ошибка при чтении rows базы данных"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var planner Tsk
		var dateInt int
		err = rows.Scan(&planner.Id, &dateInt, &planner.Title, &planner.Comment, &planner.Repeat)
		if err != nil {
			log.Printf("ошибка при записи данных в переменные, при чтении rows базы данных: %v", err)
			http.Error(res, `{"error":"ошибка при записи данных в переменные, при чтении rows базы данных"}`, http.StatusInternalServerError)
			return
		}
		planner.Date = fmt.Sprintf("%d", dateInt)
		tasks.Tasks = append(tasks.Tasks, planner)
	}
	if err = rows.Err(); err != nil {
		log.Printf("ошибка при обходе rows: %v", err)
		http.Error(res, `{"error":"ошибка при обходе rows:"}`, http.StatusInternalServerError)
		return
	}

	if len(tasks.Tasks) == 0 {
		// Возвращаем пустой список задач
		emptyResponse := Planner{Tasks: []Tsk{}}
		jsonResponse, err := json.Marshal(emptyResponse)
		if err != nil {
			log.Printf("ошибка при создании JSON-ответа: %v", err)
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
		log.Printf("ошибка при создании JSON-ответа: %v", err)
		http.Error(res, `{"error":"ошибка при создании JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write(jsonResponse)
}
