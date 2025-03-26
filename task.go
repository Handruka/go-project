package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func Task(res http.ResponseWriter, req *http.Request) { // функция возвращает задачу из базы данных в JSON-формате
	id := req.URL.Query().Get("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("неверный формат ID: %v", err)
		http.Error(res, `{"error":"неверный формат ID"}`, http.StatusBadRequest)
		return
	}
	var scanId int
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tasks Tsk
	query := "SELECT * FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, idInt)
	err = row.Scan(&scanId, &tasks.Date, &tasks.Title, &tasks.Comment, &tasks.Repeat)
	if err != nil {
		log.Printf("ошибка при сканировании row базы данных: %v", err)
		http.Error(res, `{"error":"ошибка при сканировании row базы данных"}`, http.StatusInternalServerError)
		return
	}
	if scanId == idInt {
		tasks.Id = id
	}
	response, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("ошибка при создании JSON-ответа: %v", err)
		http.Error(res, `{"error":"ошибка при создании JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write(response)

}
