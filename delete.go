package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func Delete(res http.ResponseWriter, req *http.Request) { // функция удаляет задачу
	var task Tsk
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	idStr := req.URL.Query().Get("id")

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Id не может быть строкой или слишком длинным числом: %v", err)
		http.Error(res, `{"error":"Id не может быть строкой или слишком длинным числом"}`, http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = ?", idInt)
	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Printf("ошибка при сканировании row базы данных: %v", err)
		http.Error(res, `{"error":"ошибка при сканировании row базы данных"}`, http.StatusInternalServerError)
		return
	}

	queryDelete := "DELETE FROM scheduler WHERE id = ?"
	result, err := db.Exec(queryDelete, idInt)
	if err != nil {
		log.Printf("ошибка при обращении к базе данных: %v", err)
		http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("ошибка при получении результата Exec: %v", err)
		http.Error(res, `{"error":"ошибка при получении результата Exec"}`, http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("задача с таким id не найдена: %v", err)
		http.Error(res, `{"error":"задача с таким id не найдена"}`, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]string{})
}
