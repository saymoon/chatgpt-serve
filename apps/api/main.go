package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
    db := InitDB()
    logrus.Infoln(db)

	http.HandleFunc("/tasks/create", authMiddleware(CreateTaskHandler, db))
	http.HandleFunc("/tasks/get", authMiddleware(GetTaskHandler, db))
	http.HandleFunc("/tasks/update", authMiddleware(UpdateTaskHandler, db))
	http.HandleFunc("/tasks/get-pending", authMiddleware(GetPendingTaskHandler, db))

	log.Fatal(http.ListenAndServe(":8081", nil))
}

