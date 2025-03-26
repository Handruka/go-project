package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var dataPrs time.Time

func AddTask(res http.ResponseWriter, req *http.Request) { // функция добавляет задачу в базу данных, в соответствии с правилом повторения в JSON-формате
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Tsk

	now := time.Now()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("ошибка чтения Body в POST-запросе: %v", err)
		http.Error(res, `{"error":"ошибка чтения Body в POST-запросе"}`, http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &task); err != nil {
		log.Printf("ошибка декодирования в функции Unmarshal: %v", err)
		http.Error(res, `{"error":"ошибка декодирования в функции Unmarshal"}`, http.StatusInternalServerError)
		return
	}

	if task.Title == "" {
		log.Printf("title не может быть пустым: %v", err)
		http.Error(res, `{"error":"title не может быть пустым"}`, http.StatusBadRequest)
		return
	}
	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		dataPrs, err = time.Parse("20060102", task.Date)
		if err != nil {
			log.Printf("ошибка парсинга даты task.Date: %v", err)
			http.Error(res, `{"error":"ошибка парсинга даты task.Date"}`, http.StatusBadRequest)
			return
		}
		task.Date = dataPrs.Format("20060102")
	}
	if task.Repeat != "" && task.Date != time.Now().Format("20060102") {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("ошибка обновления даты: %v", err)
			http.Error(res, `{"error":"ошибка обновления даты"}`, http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else {
		task.Date = time.Now().Format("20060102")
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("ошибка при обращении к базе данных: %v", err)
		http.Error(res, `{"error":"ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
		return
	}
	taskIntId, err := result.LastInsertId()
	if err != nil {
		log.Printf("ошибка обработки LastInsertId: %v", err)
		http.Error(res, `{"error":"ошибка обработки LastInsertId"}`, http.StatusInternalServerError)
		return
	}
	idInt := int(taskIntId)
	task.Id = strconv.Itoa(idInt)
	jsonResponse, err := json.Marshal(task)
	if err != nil {
		log.Printf("ошибка при создании JSON-ответа: %v", err)
		http.Error(res, `{"error":"ошибка при создании JSON-ответа"}`, http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
	res.Write(jsonResponse)

}
