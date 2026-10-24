package models

type Email struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// для логов
func (e *Email) ShortName() string {
	return "[" + e.ID + "] " + e.Subject
}

func (e *Email) IsEmpty() bool {
	return e.From == "" && e.To == "" && e.Body == "" // хз, пригодится ли
}

// func (e *Email) Validate() error {
// 	return nil // потом сделаю
// }