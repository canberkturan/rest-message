package main

import (
	"time"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.", now)
		return
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth || keyVal["username"] != "admin" {
		fmt.Fprintf(w, "You don't have permission to see this page!")
	} else {
		var users []User
		result, err := db.Query("SELECT username, name, password FROM users")
		if err != nil {
			log("DBQUERYERR", "Database query process is failed", now)
			return
		}

		defer result.Close()

		for result.Next() {
			var user User
			err := result.Scan(&user.Username, &user.Name, &user.Password)
			if err != nil {
				log("VARSNOTMATCHED", "Variables are not matched with values", now)
				return
			}
			users = append(users, user)
		}

		json.NewEncoder(w).Encode(users)
		log("USERFETCH", "admin fetched all users", now)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	stmt, err := db.Prepare("INSERT INTO users(username,name,password) VALUES(?,?,?)")
	if err != nil {
		log("DBPREPERR", "Database preparing is failed", now)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.", now)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	name := keyVal["name"]
	password := keyVal["password"]

	result, err := db.Query("SELECT username FROM users WHERE username = ?", username)
	if err != nil {
		log("DBQUERYERR", "Database query process is failed", now)
		return
	}
	defer result.Close()
	var user_name string
	for result.Next() {
		err := result.Scan(&user_name)
		if err != nil {
			log("VARSNOTMATCHED", "Variables are not matched with values", now)
			return
		}
	}
	if user_name == username {
		fmt.Fprintf(w, "A user with this username is already exist.")
		log("ALREADYEXIST", "Username is already recorded", now)
	} else {
		_, err = stmt.Exec(username,name,password)
		if err != nil {
			log("DBINSERTERR", "Something went wrong", now)
			return
		}
		fmt.Fprintf(w, "New user was created")
		log("USERCREATE", "A user was created with this username: " + keyVal["username"], now)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.", now)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth || keyVal["username"] != params["username"] {
		fmt.Fprintf(w, "You don't have permission to see this user's info")
	} else {
		result, err := db.Query("SELECT username,name,password FROM users WHERE username = ?", params["username"])
		if err != nil {
			log("DBQUERYERR", "Database query process is failed", now)
			return
		}

		defer result.Close()

		var user User

		for result.Next() {
			err := result.Scan(&user.Username, &user.Name, &user.Password)
			if err != nil {
				log("VARSNOTMATCHED", "Variables are not matched with values", now)
				return
			}
		}

		json.NewEncoder(w).Encode(user)
		log("USERFETCH",keyVal["username"] + " fetched its infos", now)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.", now)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth || keyVal["username"] != params["username"] {
		fmt.Fprintf(w, "You don't have permission to update this user!")
	} else {

		stmt, err := db.Prepare("UPDATE users SET name = ? WHERE username = ?")
		if err != nil {
			log("DBPREPERR", "Database preparing is failed", now)
			return
		}
		newName := keyVal["name"]

		_, err = stmt.Exec(newName, params["username"])
		if err != nil {
			log("DBUPDATEERR", "No records found to update with WHERE statement", now)
			return
		}

		fmt.Fprintf(w, "User with Username = %s was updated", params["username"])
		log("USERUPDATE", keyVal["username"] + " updated its profile", now)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.", now)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth || keyVal["username"] != params["username"] {
		fmt.Fprintf(w,"You don't have permission to delete this user!")
	} else {
		stmt, err := db.Prepare("DELETE FROM users WHERE username = ?")
		if err != nil {
			log("DBPREPERR", "Database preparing is failed.", now)
			return
		}

		_, err = stmt.Exec(params["username"])
		if err != nil {
			log("DBDELETEERR", "No records found to delete with WHERE statement", now)
			return
		}

		fmt.Fprintf(w, "User with Username = %s was deleted", params["username"])
		log("USERDEL", params["username"] + " was deleted", now)
	}
}

func blockUser(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.", now)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth {
		fmt.Fprintf(w,"You need to authenticate for blocking a user")
	} else if keyVal["username"] == params["username"] {
		fmt.Fprintf(w, "You can't block yourself :)")
		log("SELFTBLOCK", "Someone tried to block itself.", now)
	} else {
		result, err := db.Query("SELECT blocked_username FROM blockings WHERE blocker_username = ? and blocked_username = ?",
					keyVal["username"], params["username"])
		if err != nil {
			log("DBQUERYERR", "Database query process is failed.", now)
			return
		}
		defer result.Close()

		var blocked_username string
		for result.Next() {
			err := result.Scan(&blocked_username)
			if err != nil {
				log("VARSNOTMATCHED", "Variables are not matched with values", now)
				return
			}
		}
		if blocked_username != "" {
			fmt.Fprintf(w,"You are already blocked this username: %s",params["username"])
		} else {
			stmt, err := db.Prepare("INSERT INTO blockings(blocker_username, blocked_username) VALUES(?,?)")
			if err != nil {
				log("DBPREPERR", "Database preparing is failed.", now)
				return
			}

			_, err2 := stmt.Exec(keyVal["username"],params["username"])
			if err2 != nil {
				log("DBINSERTERR", "Something went wrong.", now)
				return
			}
			fmt.Fprintf(w, "You've blocked '%s' username successfully", params["username"])
			log("BLOCK", keyVal["username"] + " blocked " + params["username"], now)
		}
	}
}

func unblockUser(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.", now)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth {
		fmt.Fprintf(w,"You need to authenticate for unblocking a user")
	} else if keyVal["username"] == params["username"] {
		fmt.Fprintf(w, "You can't unblock yourself :)")
		log("SELFBLOCK", "Someone tried to unblock itself.", now)
	} else {
		result, err := db.Query("SELECT blocked_username FROM blockings WHERE blocker_username = ? and blocked_username = ?",
					keyVal["username"], params["username"])
		if err != nil {
			log("DBQUERYERR", "Database query process is failed.", now)
			return
		}
		defer result.Close()

		var blocked_username string
		for result.Next() {
			err := result.Scan(&blocked_username)
			if err != nil {
				log("VARSNOTMATCHED", "Variables are not matched with values", now)
				return
			}
		}
		if blocked_username == "" {
			fmt.Fprintf(w,"You are not blocked this username: %s",params["username"])
		} else {
			stmt, err := db.Prepare("DELETE FROM blockings WHERE blocker_username = ? AND blocked_username = ?")
			if err != nil {
				log("DBPREPERR", "Preparing Database is failed", now)
				return
			}

			_, err2 := stmt.Exec(keyVal["username"],params["username"])
			if err2 != nil {
				log("DBDELETEERR", "No records found to delete with WHERE statement", now)
				return
			}
			fmt.Fprintf(w, "You've unblocked '%s' username successfully", params["username"])
			log("BLOCK", keyVal["username"] + " unblocked " + params["username"], now)
		}
	}
}

