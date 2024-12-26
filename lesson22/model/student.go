package model

import "time"

type Student struct {
	Id        string
	Name      string
	LastName  string
	Phone     string
	Age       int
	Grade     int
	Gender    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateStudent struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
	Age      int    `json:"age"`
	Grade    int    `json:"grade"`
	Gender   string `json:"gender"`
	CourseId string `json:"course_id"`
}

type GetStudentResp struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
	Age      int    `json:"age"`
	Grade    int    `json:"grade"`
	Gender   string `json:"gender"`
	Course   Course `json:"course"`
}

type UpdateStudentResp struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
	Age      int    `json:"age"`
	Grade    int    `json:"grade"`
	Gender   string `json:"gender"`
	Course   Course `json:"course"`
}
