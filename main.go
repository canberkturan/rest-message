package main

import (
	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
	"net/http"
	"html/template"
	"os"
)

var db *sql.DB
var logger *os.File

func main() {

	IS_SECURE := true

	db = initDB()
	logger = initLogger()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/",showIndex).Methods("GET")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{username}", getUser).Methods("GET")
	router.HandleFunc("/users/{username}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{username}", deleteUser).Methods("DELETE")
	router.HandleFunc("/users/{username}/block", blockUser).Methods("GET")
	router.HandleFunc("/users/{username}/unblock", unblockUser).Methods("GET")
	router.HandleFunc("/messages", getAllMessages).Methods("GET")
	router.HandleFunc("/messages/all", getAllMessages).Methods("GET")
	router.HandleFunc("/messages/{username}", getUserMessages).Methods("GET")
	router.HandleFunc("/messages/{username}", sendMessage).Methods("POST")
	router.HandleFunc("/message/{id}", getMessage).Methods("GET")
	router.HandleFunc("/message/{id}", updateMessage).Methods("PUT")
	router.HandleFunc("/message/{id}", deleteMessage).Methods("DELETE")
	router.NotFoundHandler = router.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	log("SERVERSTART", "Server is being started")
	if IS_SECURE {
		err := http.ListenAndServeTLS(":8080", "tls/server.crt", "tls/server.key", removeTrailingSlash(router))
		if err != nil {
			panic(err.Error())
		}
	} else {
		err := http.ListenAndServe(":8080", removeTrailingSlash(router))
		if err != nil {
			panic(err.Error())
		}
	}

	defer db.Close()
	defer logger.Close()
}

func showIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w,nil)
}
