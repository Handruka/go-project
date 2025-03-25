package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

func EditTask(res http.ResponseWriter, req *http.Request) { // функция редактирует созданную задачу
	var task Tsk
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	now := time.Now()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, `{"error":"ошибка чтения Body в POST-запросе"}`, http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		http.Error(res, `{"error":"JSON равен нулю"}`, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(res, `{"error":"ошибка декодирования в функции Unmarshal"}`, http.StatusInternalServerError)
		return
	}
	if task.Id == "" {
		http.Error(res, `{"error":"Id не может быть пустым"}`, http.StatusInternalServerError)
		return
	}

	idInt, err := strconv.ParseInt(task.Id, 10, 64)
	if err != nil {
		http.Error(res, `{"error":"Id не может быть строкой или слишком длинным числом"}`, http.StatusInternalServerError)
		return
	}

	if task.Title == "" {
		http.Error(res, `{"error":"title не может быть пустым"}`, http.StatusInternalServerError)
		return
	}
	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		dataPrs, err = time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(res, `{"error":"ошибка парсинга даты task.Date"}`, http.StatusInternalServerError)
			return
		}
		task.Date = dataPrs.Format("20060102")
	}
	if task.Repeat != "" && task.Date != time.Now().Format("20060102") {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(res, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		task.Date = nextDate
	} else {
		task.Date = time.Now().Format("20060102")
	}

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

}
