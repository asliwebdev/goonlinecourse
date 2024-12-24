package repository

import (
	"database/sql"
	"lesson22/model"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(id string, course model.Course) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec(`
		INSERT INTO course (id, name, started_date, tutor, number, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id, course.Name, course.StartedAt, course.Tutor, course.Number, course.CreatedAt, course.UpdatedAt,
	)

	return err
}

func (r *CourseRepository) GetCourse(courseID string) (*model.Course, error) {
	course := &model.Course{}
	err := r.db.QueryRow(`
		SELECT id, name, started_date, tutor, number, created_at, updated_at 
		FROM course
		WHERE id = $1`, courseID).Scan(
		&course.Id,
		&course.Name,
		&course.StartedAt,
		&course.Tutor,
		&course.Number,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (r *CourseRepository) UpdateCourse(courseID string, course model.Course) (*model.Course, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	_, err = tx.Exec(`
		UPDATE course 
		SET name = $1, started_date = $2, tutor = $3, number = $4, updated_at = $5 
		WHERE id = $6`,
		course.Name, course.StartedAt, course.Tutor, course.Number, course.UpdatedAt, courseID,
	)

	updatedCourse := &model.Course{}
	err = r.db.QueryRow(`
	select id, name, started_date, tutor, number, created_at, updated_at 
	from course
	where id = $1`, course.Id).Scan(
		&updatedCourse.Id,
		&updatedCourse.Name,
		&updatedCourse.StartedAt,
		&updatedCourse.Tutor,
		&updatedCourse.Number,
		&updatedCourse.CreatedAt,
		&updatedCourse.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return updatedCourse, nil
}

func (r *CourseRepository) DeleteCourse(courseID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec(`DELETE FROM course WHERE id = $1`, courseID)
	return err
}
