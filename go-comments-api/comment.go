package main

import (
	"database/sql"
)

type comment struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Body     string `json:"body"`
	Modtime  string `json:"mod_time"`
}

func getComments(db *sql.DB, start, count int) ([]comment, error) {
	rows, err := db.Query(
		"SELECT id, username, body, mod_time FROM comments ORDER BY mod_time LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []comment{}
	for rows.Next() {
		var c comment
		if err := rows.Scan(&c.ID, &c.Username, &c.Body, &c.Modtime); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (c *comment) getComment(db *sql.DB) error {
	return db.QueryRow("SELECT username, body, mod_time FROM comments WHERE id=$1",
		c.ID).Scan(&c.Username, &c.Body, &c.Modtime)
}

func (c *comment) deleteComment(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM comments WHERE id=$1", c.ID)
	return err
}

func (c *comment) createComment(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO comments(username, body) VALUES($1, $2) RETURNING id, mod_time",
		c.Username, c.Body).Scan(&c.ID, &c.Modtime)

	if err != nil {
		return err
	}

	return nil
}

func (c *comment) updateComment(db *sql.DB) error {
	err := db.QueryRow("UPDATE comments SET body=$1 WHERE id=$2 RETURNING username, mod_time", c.Body, c.ID).Scan(&c.Username, &c.Modtime)

	return err
}
