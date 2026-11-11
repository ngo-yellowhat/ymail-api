package models

import "fmt" //  добавил импорт

type Email struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// для логов
func (e *Email) ShortName() string {
	return fmt.Sprintf("[%s] %s", e.ID, e.Subject) //  заменил на Sprintf
}

func (e *Email) IsEmpty() bool {
	return e.From == "" && e.To == "" && e.Body == "" // Надо ли оно ?
}

// func (e *Email) Validate() error {
// 	return nil // потом сделаю
// }