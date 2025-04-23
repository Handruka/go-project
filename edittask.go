package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func EditTask(res http.ResponseWriter, req *http.Request) { // функция редактирует созданную задачу
	var task Tsk
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	now := time.Now()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("ошибка чтения Body в POST-запросе: %v", err)
		http.Error(res, `{"error":"ошибка чтения Body в POST-запросе"}`, http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		log.Printf("JSON равен нулю: %v", err)

		http.Error(res, `{"error":"JSON равен нулю"}`, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Printf("ошибка декодирования в функции Unmarshal: %v", err)
		http.Error(res, `{"error":"ошибка декодирования в функции Unmarshal"}`, http.StatusInternalServerError)
		return
	}
	if task.Id == "" {
		log.Printf("Id не может быть пустым: %v", err)
		http.Error(res, `{"error":"Id не может быть пустым"}`, http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(task.Id, 10, 64)
	if err != nil {
		log.Printf("Id не может быть строкой или слишком длинным числом: %v", err)
		http.Error(res, `{"error":"Id не может быть строкой или слишком длинным числом"}`, http.StatusBadRequest)
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
			log.Printf("ошибка проверки даты: %v", err)
			http.Error(res, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		task.Date = nextDate
	} else {
		task.Date = time.Now().Format("20060102")
	}

	dateInt, err := strconv.Atoi(task.Date)
	if err != nil {
		log.Printf("ошибка при конвертации task.Date в intId: %v", err)
		http.Error(res, `{"error":"ошибка при конвертации task.Date в intId"}`, http.StatusInternalServerError)
		return
	}

	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	result, err := db.Exec(query, dateInt, task.Title, task.Comment, task.Repeat, idInt)
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
