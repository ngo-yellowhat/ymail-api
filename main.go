package main
 
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
 
	"mail-api/models"
)
 
var (
	mails = []models.Email{}
	mu    sync.Mutex //  гонки 
)
 
func main() {
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})
 
	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
 
		switch r.Method {
		case "GET":
			mu.Lock()
			json.NewEncoder(w).Encode(mails)
			mu.Unlock()
 
		case "POST":
			var msg models.Email
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				log.Println("Ошибка парсинга JSON:", err)
				http.Error(w, `{"err": "json broken"}`, 400) // ПЕРЕДЕЛать
				return
			}
						if msg.IsEmpty() {
				http.Error(w, `{"err": "пустое письмо"}`, 400)
				return
			}
 
			// Генер id
			msg.ID = fmt.Sprintf("%d", time.Now().UnixNano())
 
			mu.Lock()
			mails = append(mails, msg)
			mu.Unlock()
 
			log.Println("Добавлено письмо:", msg.ShortName())
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
 
		mu.Lock()
		defer mu.Unlock()
 
		for i, m := range mails {
			if m.ID == id {
				mails = append(mails[:i], mails[i+1:]...)
				log.Println("Удалили письмо с ID:", id)
				fmt.Fprintf(w, `{"ok": true}`)
				return
			}
		}
 
		http.Error(w, `{"error": "no such email"}`, 404)
	})
 
	fmt.Println("Готово")
	http.ListenAndServe("localhost:8080", nil)
}