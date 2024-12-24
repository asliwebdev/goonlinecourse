package model

import (
	"github.com/google/uuid"
)

type StudentCourse struct {
	ID        uuid.UUID `json:"id"`
	StudentID uuid.UUID `json:"student_id"`
	CourseID  uuid.UUID `json:"course_id"`
}
