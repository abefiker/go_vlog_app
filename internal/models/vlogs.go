package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Vlog struct to hold the data for an individual vlog.
type Vlog struct {
	VlogID      int       `db:"vlog_id"`
	UserID      int       `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	VideoFile   string    `db:"video_file"`
	PhotoFile   string    `db:"photo_file"`
	Views       int       `db:"views"`
	Likes       int       `db:"likes"`
	CreatedAt   time.Time `db:"created_at"`
	// Add other fields as needed
}

// Define a VlogModel type which wraps a sql.DB connection pool.
type VlogModel struct {
	DB *sql.DB
}

// This will insert a new vlog into the database.
func (m *VlogModel) Insert(user_id int, title string, description string, photoFile string, views int, likes int) (int, error) {
	// Implement your insert logic here
	stmt := `INSERT INTO vlogs (user_id ,title, description, photo_file, views, likes, created_at) VALUES (?,?, ?, ?, ?, ?, NOW())`
	result, err := m.DB.Exec(stmt, user_id, title, description, photoFile, views, likes)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// This will return a specific vlog based on its id.
func (m *VlogModel) Get(id int) (*Vlog, error) {
	// Implement your get logic here
	stmt := `SELECT vlog_id , user_id ,title, description, photo_file, views, likes, created_at FROM vlogs
		WHERE vlog_id = ?`
	row := m.DB.QueryRow(stmt, id)
	v := &Vlog{}
	err := row.Scan(&v.VlogID, &v.UserID, &v.Title, &v.Description, &v.PhotoFile, &v.Views, &v.Likes, &v.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return v, nil
}

// This will return the 10 most recently created vlogs.
func (m *VlogModel) Latest() ([]*Vlog, error) {
	// Implement your latest logic here
	stmt := `SELECT vlog_id, user_id, title, description, photo_file, views, likes, created_at 
         FROM vlogs 
         ORDER BY vlog_id DESC 
         LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	vlogs := []*Vlog{}
	for rows.Next() {
		v := &Vlog{}
		err := rows.Scan(&v.VlogID, &v.UserID, &v.Title, &v.Description, &v.PhotoFile, &v.Views, &v.Likes, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		vlogs = append(vlogs, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return vlogs, nil
}
