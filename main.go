package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mail-api/models"
)

var mails = []models.Email{}

func main() {
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})

	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			json.NewEncoder(w).Encode(mails)
		
		case "POST":
			var msg models.Email
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				http.Error(w, `{"err": "json broken"}`, 400) // ПЕРЕДЕЛать
				return
			}

			msg.ID = fmt.Sprintf("%d", len(mails)+1)
			mails = append(mails, msg)

			w.WriteHeader(201)
			json.NewEncoder(w).Encode(msg)
		
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, `{"error": "only delete here"}`, 405)
			return
		}

		id := r.URL.Path[len("/api/messages/"):]

		for i, m := range mails {
			if m.ID == id {
				mails = append(mails[:i], mails[i+1:]...)
				fmt.Fprintf(w, `{"ok": true}`)
				return
			}
		}

		http.Error(w, `{"error": "no such email"}`, 404)
	})

	fmt.Println("пашЫт и ебашЫт")
	http.ListenAndServe("localhost:8080", nil)
}