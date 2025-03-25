package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

var dataPrs time.Time

func AddTask(res http.ResponseWriter, req *http.Request) { // функция добавляет задачу в базу данных, в соответствии с правилом повторения в JSON-формате
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var task Tsk

	now := time.Now()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, `{"error":"ошибка чтения Body в POST-запросе"}`, http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &task); err != nil {
		http.Error(res, `{"error":"ошибка декодирования в функции Unmarshal"}`, http.StatusInternalServerError)
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
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}
		task.Date = nextDate
	} else {
		task.Date = time.Now().Format("20060102")
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
		return
	}
	taskIntId, err := result.LastInsertId()
	if err != nil {
		http.Error(res, `{"error":"ошибка обработки LastInsertId"}`, http.StatusInternalServerError)
		return
	}
	idInt := int(taskIntId)
	task.Id = strconv.Itoa(idInt)
	jsonResponse, err := json.Marshal(task)
	if err != nil {
		http.Error(res, `{"error":"ошибка при создании JSON-ответа"}`, http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
	res.Write(jsonResponse)

}
