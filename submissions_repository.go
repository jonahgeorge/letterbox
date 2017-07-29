package main

import (
	"database/sql"
)

const (
	SUBMISSIONS_INSERT_SQL          = `insert into submissions (form_id, body) values ($1, $2) returning *`
	SUBMISSIONS_FIND_BY_FORM_ID_SQL = `select * from submissions where form_id = $1`
)

type SubmissionsRepository struct {
	db *sql.DB
}

func NewSubmissionsRepository(db *sql.DB) *SubmissionsRepository {
	return &SubmissionsRepository{
		db: db,
	}
}

func (repo *SubmissionsRepository) FindByFormId(formId int) ([]Submission, error) {
	var submissions []Submission
	rows, err := repo.db.Query(SUBMISSIONS_FIND_BY_FORM_ID_SQL, formId)

	for rows.Next() {
		submission := new(Submission)
		err = repo.scanRows(rows, submission)
		submissions = append(submissions, *submission)
	}

	return submissions, err
}

func (repo *SubmissionsRepository) Create(formId int, body string) (*Submission, error) {
	submission := new(Submission)
	row := repo.db.QueryRow(SUBMISSIONS_INSERT_SQL, formId, body)
	err := repo.scanRow(row, submission)
	return submission, err
}

func (repo *SubmissionsRepository) scanRow(row *sql.Row, submission *Submission) error {
	return row.Scan(&submission.id, &submission.formId, &submission.body, &submission.createdAt, &submission.updatedAt)
}

func (repo *SubmissionsRepository) scanRows(rows *sql.Rows, submission *Submission) error {
	return rows.Scan(&submission.id, &submission.formId, &submission.body, &submission.createdAt, &submission.updatedAt)
}
