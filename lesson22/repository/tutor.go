package repository

import (
	"database/sql"
	"lesson22/model"

	"github.com/google/uuid"
)

type TutorRepository struct {
	db *sql.DB
}

func NewTutorRepository(db *sql.DB) *TutorRepository {
	return &TutorRepository{db: db}
}

func (r *TutorRepository) CreateTutor(tutor *model.Tutor) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Commit()

	id := uuid.NewString()
	_, err = tx.Exec("insert into tutors (id, name, subject, phone, created_at) values ($1, $2, $3, $4, NOW())", id, tutor.Name, tutor.Subject, tutor.Phone)

	if err != nil {
		return err
	}
	return nil
}

func (r *TutorRepository) GetTutor(id string) (*model.Tutor, error) {
	tutor := &model.Tutor{}
	err := r.db.QueryRow("select id, name, subject, phone, created_at from tutors where id = $1", id).Scan(
		&tutor.Id, &tutor.Name, &tutor.Subject, &tutor.Phone, &tutor.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tutor, nil
}

func (r *TutorRepository) UpdateTutor(tutor *model.Tutor) (*model.Tutor, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	_, err = tx.Exec("update tutors set name = $1, subject = $2, phone = $3 where id = $4", tutor.Name, tutor.Subject, tutor.Phone, tutor.Id)
	if err != nil {
		return nil, err
	}

	updatedTutor := &model.Tutor{}
	err = r.db.QueryRow(`
	select id, name, subject, phone, created_at 
	from tutor
	where id = $1`, tutor.Id).Scan(
		&updatedTutor.Id,
		&updatedTutor.Name,
		&updatedTutor.Subject,
		&updatedTutor.Phone,
		&updatedTutor.CreatedAt,
	)

	return updatedTutor, err
}

func (r *TutorRepository) DeleteTutor(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("delete from tutors where id = $1", id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *TutorRepository) GetListTutors() ([]*model.Tutor, error) {
	rows, err := r.db.Query("select id, name, subject, phone, created_at from tutors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tutors []*model.Tutor
	for rows.Next() {
		tutor := &model.Tutor{}
		err := rows.Scan(&tutor.Id, &tutor.Name, &tutor.Subject, &tutor.Phone, &tutor.CreatedAt)
		if err != nil {
			return nil, err
		}
		tutors = append(tutors, tutor)
	}

	return tutors, nil
}
