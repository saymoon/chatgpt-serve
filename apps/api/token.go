package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"

	"github.com/sirupsen/logrus"
)

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func authMiddleware(next http.HandlerFunc, db *sql.DB) http.HandlerFunc {
    logrus.Infoln(db)
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		row := db.QueryRow("SELECT * FROM tokens WHERE token = ?", tokenStr)
		var token Token
		err := row.Scan(&token.Token, &token.IsAdmin)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid token.", http.StatusUnauthorized)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// 把 token 信息保存在 request context 中，以便在后续的处理函数中使用
		ctx := context.WithValue(r.Context(), "token", token)
		ctx = context.WithValue(ctx, "db", db)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

