package main

import (
	"os"
	"strings"
	"net/http"
	"database/sql"
	"time"
	"fmt"
)
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			}
			next.ServeHTTP(w, r)
		})
}

func initDB() *sql.DB {
	db, _ := sql.Open("sqlite3", "./data.db")
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users (username TEXT, name TEXT, password TEXT)")
	stmt.Exec()
	stmt2, _ := db.Prepare("CREATE TABLE IF NOT EXISTS messages(id INTEGER, sender TEXT, receiver TEXT, content TEXT, senddate TEXT, readdate TEXT)")

	stmt2.Exec()
	stmt3, _ := db.Prepare("CREATE TABLE IF NOT EXISTS blockings(blocker_username TEXT, blocked_username TEXT)")
	stmt3.Exec()

	return db
}

func initLogger() *os.File {
	f, err := os.Create("server.log")
        if err != nil {
		fmt.Println("Log file was not open somehow.")
		return nil
	}
	return f
}

func authCheck(username, password string) bool {
	now := time.Now().String()
	result, err := db.Query("SELECT password FROM users WHERE username = ?", username)

	if err != nil {
		log("DBQUERYFAIL", "Database query process is failed", now)
		return false
	}

	defer result.Close()

	var realPassword string

	for result.Next() {
		err := result.Scan(&realPassword)
		if err != nil {
			log("VARSNOTMATCH", "Variables are not matched with values", now)
			return false
		}
	}
	if realPassword == password {
		log("AUTHSUCC", "Username and password is correct.", now)
		return true
	}
	log("AUTHFAIL", "Username or password is wrong", now)
	return false
}

func addReadInfo(messageID string, read_date string) {
	now := time.Now().String()
	stmt, err := db.Prepare("UPDATE messages SET readdate = ? WHERE id = ?")
	if err != nil {
		log("DBPREPERR", "Preparing Database is failed", now)
		return
	}
	_, err2 := stmt.Exec(read_date, messageID)
	if err2 != nil {
		log("DBUPTDATEERR", "No records found to update with WHERE statement", now)
		return
	}
}

func log(logtype, message, date string) {
	log_message := logtype+": "+message+" \t"+date+"\n"
	_, err := logger.WriteString(log_message)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(log_message)
}
