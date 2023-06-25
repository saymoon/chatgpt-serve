package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("db").(*sql.DB)
	// 检查权限
	_, ok = r.Context().Value("token").(Token)
	if !ok {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Infoln(task)

	stmt, err := db.Prepare("INSERT INTO tasks(id, instance_id, conversation_id, model, prompt, response, status, error_message, created_at, updated_at) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(task.Id, task.InstanceId, task.ConversationId, task.Model, task.Prompt, "", TASK_STATUS_PENDING, "", time.Now(), time.Now())
	if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New task created successfully.")
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("db").(*sql.DB)
	// 检查权限
	_, ok = r.Context().Value("token").(Token)
	if !ok {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "`id` parameter is required", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
	var task Task
	err := row.Scan(&task.Id, &task.InstanceId, &task.ConversationId, &task.Model, &task.Prompt, &task.Response, &task.Status, &task.ErrorMessage, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
        logrus.Errorln(err.Error())
		if err == sql.ErrNoRows {
			http.Error(w, "No task found.", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	js, err := json.Marshal(task)
	if err != nil {
        logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func GetPendingTaskHandler(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("db").(*sql.DB)
	// 检查权限
    token, ok := r.Context().Value("token").(Token)
	if !ok || !token.IsAdmin {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	row := db.QueryRow("SELECT * FROM tasks WHERE status = ? ORDER BY created_at ASC LIMIT 1", "pending")
	var task Task
	err := row.Scan(&task.Id, &task.InstanceId, &task.ConversationId, &task.Model, &task.Prompt, &task.Response, &task.Status, &task.ErrorMessage, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
        logrus.Errorln(err.Error())
		if err == sql.ErrNoRows {
			http.Error(w, "No task found.", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	js, err := json.Marshal(task)
	if err != nil {
        logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("db").(*sql.DB)
	// 检查权限
	token, ok := r.Context().Value("token").(Token)
	if !ok || !token.IsAdmin {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
        logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE tasks SET conversation_id = ?, response = ?, status = ?, error_message = ?, updated_at = ? WHERE id = ?")
	if err != nil {
        logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    logrus.Info(task)
	_, err = stmt.Exec(task.ConversationId, task.Response, task.Status, task.ErrorMessage, time.Now(), task.Id)
	if err != nil {
        logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Task updated successfully.")
}

