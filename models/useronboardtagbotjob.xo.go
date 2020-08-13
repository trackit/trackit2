// Package models contains the types for schema 'trackit'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// UserOnboardTagbotJob represents a row from 'trackit.user_onboard_tagbot_job'.
type UserOnboardTagbotJob struct {
	ID        int       `json:"id"`        // id
	Created   time.Time `json:"created"`   // created
	UserID    int       `json:"user_id"`   // user_id
	Completed time.Time `json:"completed"` // completed
	WorkerID  string    `json:"worker_id"` // worker_id
	JobError  string    `json:"job_error"` // job_error

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the UserOnboardTagbotJob exists in the database.
func (uotj *UserOnboardTagbotJob) Exists() bool {
	return uotj._exists
}

// Deleted provides information if the UserOnboardTagbotJob has been deleted from the database.
func (uotj *UserOnboardTagbotJob) Deleted() bool {
	return uotj._deleted
}

// Insert inserts the UserOnboardTagbotJob to the database.
func (uotj *UserOnboardTagbotJob) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if uotj._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO trackit.user_onboard_tagbot_job (` +
		`created, user_id, completed, worker_id, job_error` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, uotj.Created, uotj.UserID, uotj.Completed, uotj.WorkerID, uotj.JobError)
	res, err := db.Exec(sqlstr, uotj.Created, uotj.UserID, uotj.Completed, uotj.WorkerID, uotj.JobError)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	uotj.ID = int(id)
	uotj._exists = true

	return nil
}

// Update updates the UserOnboardTagbotJob in the database.
func (uotj *UserOnboardTagbotJob) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !uotj._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if uotj._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE trackit.user_onboard_tagbot_job SET ` +
		`created = ?, user_id = ?, completed = ?, worker_id = ?, job_error = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, uotj.Created, uotj.UserID, uotj.Completed, uotj.WorkerID, uotj.JobError, uotj.ID)
	_, err = db.Exec(sqlstr, uotj.Created, uotj.UserID, uotj.Completed, uotj.WorkerID, uotj.JobError, uotj.ID)
	return err
}

// Save saves the UserOnboardTagbotJob to the database.
func (uotj *UserOnboardTagbotJob) Save(db XODB) error {
	if uotj.Exists() {
		return uotj.Update(db)
	}

	return uotj.Insert(db)
}

// Delete deletes the UserOnboardTagbotJob from the database.
func (uotj *UserOnboardTagbotJob) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !uotj._exists {
		return nil
	}

	// if deleted, bail
	if uotj._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM trackit.user_onboard_tagbot_job WHERE id = ?`

	// run query
	XOLog(sqlstr, uotj.ID)
	_, err = db.Exec(sqlstr, uotj.ID)
	if err != nil {
		return err
	}

	// set deleted
	uotj._deleted = true

	return nil
}

// User returns the User associated with the UserOnboardTagbotJob's UserID (user_id).
//
// Generated from foreign key 'user_onboard_tagbot_job_ibfk_1'.
func (uotj *UserOnboardTagbotJob) User(db XODB) (*User, error) {
	return UserByID(db, uotj.UserID)
}

// UserOnboardTagbotJobsByUserID retrieves a row from 'trackit.user_onboard_tagbot_job' as a UserOnboardTagbotJob.
//
// Generated from index 'foreign_user'.
func UserOnboardTagbotJobsByUserID(db XODB, userID int) ([]*UserOnboardTagbotJob, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, created, user_id, completed, worker_id, job_error ` +
		`FROM trackit.user_onboard_tagbot_job ` +
		`WHERE user_id = ?`

	// run query
	XOLog(sqlstr, userID)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*UserOnboardTagbotJob{}
	for q.Next() {
		uotj := UserOnboardTagbotJob{
			_exists: true,
		}

		// scan
		err = q.Scan(&uotj.ID, &uotj.Created, &uotj.UserID, &uotj.Completed, &uotj.WorkerID, &uotj.JobError)
		if err != nil {
			return nil, err
		}

		res = append(res, &uotj)
	}

	return res, nil
}

// UserOnboardTagbotJobByID retrieves a row from 'trackit.user_onboard_tagbot_job' as a UserOnboardTagbotJob.
//
// Generated from index 'user_onboard_tagbot_job_id_pkey'.
func UserOnboardTagbotJobByID(db XODB, id int) (*UserOnboardTagbotJob, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, created, user_id, completed, worker_id, job_error ` +
		`FROM trackit.user_onboard_tagbot_job ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	uotj := UserOnboardTagbotJob{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&uotj.ID, &uotj.Created, &uotj.UserID, &uotj.Completed, &uotj.WorkerID, &uotj.JobError)
	if err != nil {
		return nil, err
	}

	return &uotj, nil
}
