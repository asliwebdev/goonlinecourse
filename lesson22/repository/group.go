package repository

import (
	"database/sql"
	"lesson22/model"

	"github.com/google/uuid"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) CreateGroup(group model.Group) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Commit()

	id := uuid.NewString()
	_, err = tx.Exec("insert into groups (id, name, created_at) values ($1, $2, NOW())", id, group.Name)
	return err
}

func (r *GroupRepository) GetGroup(id string) (*model.Group, error) {
	group := &model.Group{}
	err := r.db.QueryRow("select id, name, created_at from groups where id = $1", id).Scan(
		&group.Id, &group.Name, &group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) UpdateGroup(group model.Group) (*model.Group, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	_, err = tx.Exec("update groups set name = $1 where id = $2", group.Name, group.Id)

	if err != nil {
		return nil, err
	}

	updatedGroup := &model.Group{}
	err = r.db.QueryRow(`
	select id, name, created_at
	from course
	where id = $1`, group.Id).Scan(
		&updatedGroup.Id,
		&updatedGroup.Name,
		&updatedGroup.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return updatedGroup, nil
}

func (r *GroupRepository) DeleteGroup(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("delete from groups where id = $1", id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *GroupRepository) GetListGroups() ([]*model.Group, error) {
	rows, err := r.db.Query("select id, name, created_at from groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		group := &model.Group{}
		err := rows.Scan(&group.Id, &group.Name, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}
