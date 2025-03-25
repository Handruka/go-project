package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func Done(res http.ResponseWriter, req *http.Request) { // функция отмечает задачу выполненой и удаляет либо назначает новую дату в соответствии с правило "repeat"
	var task Tsk
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	idStr := req.URL.Query().Get("id")

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(res, `{"error":"Id не может быть строкой или слишком длинным числом"}`, http.StatusInternalServerError)
		return
	}

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = ?", idInt)
	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		http.Error(res, `{"error":"ошибка при сканировании row базы данных"}`, http.StatusInternalServerError)
		return
	}
	if task.Repeat != "" {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(res, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		task.Date = nextDate
		dateInt, err := strconv.Atoi(task.Date)
		if err != nil {
			http.Error(res, `{"error":"ошибка при конвертации task.Date в intId"}`, http.StatusInternalServerError)
			return
		}

		query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
		result, err := db.Exec(query, dateInt, task.Title, task.Comment, task.Repeat, idInt)
		if err != nil {
			http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(res, `{"error":"ошибка при получении результата Exec"}`, http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(res, `{"error":"задача с таким id не найдена"}`, http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(map[string]string{})

	} else {
		queryDelete := "DELETE FROM scheduler WHERE id = ?"
		result, err := db.Exec(queryDelete, idInt)
		if err != nil {
			http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(res, `{"error":"ошибка при получении результата Exec"}`, http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(res, `{"error":"задача с таким id не найдена"}`, http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(map[string]string{})
	}

}
