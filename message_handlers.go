package main

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/gorilla/mux"
	"strconv"
)


func getAllMessages(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.")
		return
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth {
		fmt.Fprintf(w, "You have no permission to see your messages")
	} else {
		var messages []Message
		result, err := db.Query(`SELECT * FROM messages WHERE sender = ? or (receiver = ? and sender NOT IN
					(SELECT blocked_username FROM blockings WHERE blocker_username = ?))`,
					keyVal["username"], keyVal["username"], keyVal["username"])
		if err != nil {
			log("DBQUERYERR", "Database query process is failed")
			return
		}

		defer result.Close()

		for result.Next() {
			var message Message
			err := result.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.SendDate, &message.ReadDate)
			if err != nil {
				log("VARSNOTMATCHED", "Variables are not matched with values")
				return
			}
			messages = append(messages, message)
		}
		for i, _ := range messages {
			if messages[i].Receiver == keyVal["username"] && messages[i].ReadDate == "-" {
				messages[i].ReadDate = time.Now().String()
				addReadInfo(messages[i].ID)
			}
		}
		json.NewEncoder(w).Encode(messages)
		log("CHATFETCH", keyVal["username"] + " fetched all messages")
	}
}
func getMessage(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.")
		return
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	params := mux.Vars(r)
	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth {
		fmt.Fprintf(w, "You have no permission to see your messages")
	} else {
		var message Message
		result, err := db.Query(`SELECT * FROM messages WHERE id = ? and (sender = ? or (receiver = ? and sender NOT IN 
					(SELECT blocked_username FROM blockings WHERE blocker_username = ?)))`,
					params["id"], keyVal["username"], keyVal["username"], keyVal["username"])
		if err != nil {
			log("DBQUERYERR", "Database query process is failed")
			return
		}

		defer result.Close()

		for result.Next() {
			err := result.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.SendDate, &message.ReadDate)
			if err != nil {
				log("VARSNOTMATCHED", "Variables are not matched with values")
				return
			}
		}
		if message.Receiver == keyVal["username"] && message.ReadDate == "-" {
			message.ReadDate = now
			addReadInfo(message.ID)
		}
		json.NewEncoder(w).Encode(message)
		log("CHATFETCH", keyVal["username"] + " fetched a message with its id")
	}
}

func getUserMessages(w http.ResponseWriter, r *http.Request){
	now := time.Now().String()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.")
		return
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	params := mux.Vars(r)
	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth {
		fmt.Fprintf(w, "You have no permission to see your messages")
	} else {
		var messages []Message
		result, err := db.Query(`SELECT * FROM messages WHERE (sender = ? and receiver = ?) or 
					(receiver = ? and sender = ? and sender NOT IN 
					(SELECT blocked_username FROM blockings 
					WHERE blocker_username = ? and blocked_username = ?))`,
					keyVal["username"], params["username"], keyVal["username"],
					params["username"], keyVal["username"], params["username"])
		if err != nil {
			log("DBQUERYERR", "Database query process is failed")
			return
		}

		defer result.Close()

		for result.Next() {
			var message Message
			err := result.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.SendDate, &message.ReadDate)
			if err != nil {
				log("VARSNOTMATCHED", "Variables are not matched with values")
				return
			}
			messages = append(messages, message)
		}
		for i, _ := range messages {
			if messages[i].Receiver == keyVal["username"] && messages[i].ReadDate == "-" {
				messages[i].ReadDate = now
				addReadInfo(messages[i].ID)
			}
		}
		json.NewEncoder(w).Encode(messages)
		log("CHATFETCH", keyVal["username"] + " fetched all chat with " + params["username"])
	}
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	stmt, err := db.Prepare("INSERT INTO messages(id,sender,receiver,content,senddate,readdate) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log("DBPREPERR", "Database preparing failed")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.")
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	params := mux.Vars(r)
	t := time.Now()
	id := keyVal["username"] + "_" + params["username"] + "_" +strconv.FormatInt(t.UnixNano(),10)
	_, err = stmt.Exec(id,keyVal["username"], params["username"], keyVal["content"], now, "-")
	if err != nil {
		log("DBINSERTERR", "Something went wrong")
		return
	}
	fmt.Fprintf(w, "Message Successfully sent.")
	log("USERCREATE", keyVal["username"] + " sent a message to " + params["username"])
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.")
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth {
		fmt.Fprintf(w, "You don't have permission to update a message!")
	} else {

		stmt, err := db.Prepare("UPDATE message SET content = ? WHERE id = ? and sender = ?")
		if err != nil {
			log("DBPREPERR", "Database preparing is failed")
			return
		}
		newContent := keyVal["content"]

		_, err = stmt.Exec(newContent, params["id"], keyVal["username"])
		if err != nil {
			log("DBUPDATEERR", "No records found to update with WHERE statement")
			return
		}

		fmt.Fprintf(w, "Message was successfully updated")
		log("CHATUPDATE", keyVal["username"] + " updated a message")
	}
}
func deleteMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("READREQERR", "Reading request is failed.")
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	isAuth := authCheck(keyVal["username"], keyVal["password"])
	if !isAuth {
		fmt.Fprintf(w,"You don't have permission to delete a message!")
	} else {
		stmt, err := db.Prepare("DELETE FROM messages WHERE id = ? and sender = ?")
		if err != nil {
			log("DBPREPERR", "Database preparing is failed.")
			return
		}

		_, err = stmt.Exec(params["id"], keyVal["username"])
		if err != nil {
			log("DBDELETEERR", "No records found to delete with WHERE statement")
			return
		}

		fmt.Fprintf(w, "Message was successfully deleted")
		log("CHATDELETE", params["username"] + " deleted a message")
	}
}
