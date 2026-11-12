package models

import "fmt"

const MaxBodyLength = 10 * 1024 * 1024 // 10 MB 

type Email struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// для логов
func (e *Email) ShortName() string {
	return fmt.Sprintf("[%s] %s", e.ID, e.Subject)
}

func (e *Email) IsEmpty() bool {
	// валид длины
	if len(e.Body) > MaxBodyLength {
		return true
	}
	// невалид если пустое
	return e.From == "" || e.To == "" || e.Body == ""
}

// func (e *Email) Validate() error {
// 	return nil // потом сделаю
// }
// TODO: вынести вложения в multipart/form-data или S3, в JSON столько места не надо