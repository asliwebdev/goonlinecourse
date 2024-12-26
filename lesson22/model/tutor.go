package model

import "time"

type Tutor struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
