package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/google/uuid"
	"mail-api/models"
)

var (
	mails = []models.Email{}
	mu    sync.Mutex // гонки
)

// чтобы ошибки тоже json отдавали
func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": %q}`, msg)
}

func main() {
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})

	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			// чтобы медленный клиент не держал мьютекс
			mu.Lock()
			snapshot := slices.Clone(mails)
			mu.Unlock()
			
			json.NewEncoder(w).Encode(snapshot)

		case "POST":
			var msg models.Email
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				log.Println("Ошибка парсинга JSON:", err)
				jsonError(w, "json broken", 400) // ПЕРЕДЕЛать
				return
			}
			if msg.IsEmpty() {
				jsonError(w, "пустое или слишком длинное письмо", 400)
				return
			}

			// uuid вместо нано, чтоб не было дублей
			msg.ID = uuid.New().String()

			mu.Lock()
			mails = append(mails, msg)
			mu.Unlock()

			// %q 
			log.Printf("added %q", msg.ShortName())
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(msg)

		default:
			jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			jsonError(w, "only delete here", 405)
			return
		}

		id := r.URL.Path[len("/api/messages/"):]

		mu.Lock()
		defer mu.Unlock()

		for i, m := range mails {
			if m.ID == id {
				mails = append(mails[:i], mails[i+1:]...)
				log.Println("Удалили письмо с ID:", id)
				w.Write([]byte(`{"ok": true}`))
				return
			}
		}

		jsonError(w, "no such email", 404)
	})

	fmt.Println("Слушаю порт 8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}