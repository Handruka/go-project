package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func Task(res http.ResponseWriter, req *http.Request) { // функция возвращает задачу из базы данных в JSON-формате
	id := req.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(res, `{"error":"неверный формат ID"}`, http.StatusBadRequest)
		return
	}
	var scanId int
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var tasks Tsk
	query := "SELECT * FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, idInt)
	err = row.Scan(&scanId, &tasks.Date, &tasks.Title, &tasks.Comment, &tasks.Repeat)
	if err != nil {
		http.Error(res, `{"error":"ошибка при сканировании row базы данных"}`, http.StatusInternalServerError)
		return
	}
	if scanId == idInt {
		tasks.Id = id
	}
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, `{"error":"ошибка при создании JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write(response)

}
