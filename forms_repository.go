package main

import (
	"database/sql"
	"log"
)

const (
	FORMS_FIND_BY_ID_SQL      = `select * from forms where id = $1`
	FORMS_FIND_BY_UUID_SQL    = `select * from forms where uuid = $1`
	FORMS_FIND_BY_USER_ID_SQL = `select * from forms where user_id = $1`
	FORMS_INSERT_SQL          = `insert into forms (user_id, name, description) values ($1, $2, $3) returning *`
)

type FormsRepository struct {
	db *sql.DB
}

func NewFormsRepository(db *sql.DB) *FormsRepository {
	return &FormsRepository{
		db: db,
	}
}

func (repo *FormsRepository) FindByUserId(userId int) ([]Form, error) {
	var forms []Form

	rows, err := repo.db.Query(FORMS_FIND_BY_USER_ID_SQL, userId)
	for rows.Next() {
		form := new(Form)
		err := repo.scanRows(rows, form)
		forms = append(forms, *form)
	}

	return forms
}

func (repo *FormsRepository) Create(userId int, name, description string) (*Form, error) {
	form := new(Form)
	row := repo.db.QueryRow(FORMS_INSERT_SQL, userId, name, description)
	err := repo.scanRow(row, form)
	return form, err
}

func (repo *FormsRepository) FindById(id int) (*Form, error) {
	form := new(Form)
	row := repo.db.QueryRow(FORMS_FIND_BY_ID_SQL, id)
	err := repo.scanRow(row, form)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) FindByUuid(uuid string) (*Form, error) {
	form := new(Form)
	row := repo.db.QueryRow(FORMS_FIND_BY_UUID_SQL, uuid)
	err := repo.scanRow(row, form)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) scanRow(row *sql.Row, form *Form) error {
	return row.Scan(&form.id, &form.userId, &form.uuid, &form.name, &form.description, &form.createdAt, &form.updatedAt)
}

func (repo *FormsRepository) scanRows(rows *sql.Rows, form *Form) error {
	return rows.Scan(&form.id, &form.userId, &form.uuid, &form.name, &form.description, &form.createdAt, &form.updatedAt)
}
