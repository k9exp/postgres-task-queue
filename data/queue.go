package data

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func SetupQueue() error {
	connString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	DB = db // Set the global variable

	// CREATE QUEUE TABLE
	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS queue (
			task_id SERIAL PRIMARY KEY,
			text TEXT NOT NULL,
			time INT NOT NULL
		);
	`)

	if err != nil {
		return err
	}

	log.Println("Table \"queue\" is created in the database")

	return nil
}

func InsertTask(text string, time uint32) error {
	_, err := DB.Exec("INSERT INTO queue (text, time) VALUES ($1, $2);", text, time)

	return err
}

type TaskData struct {
	Task_id uint32
	Text    string
	Time    uint32
}

func GetTask() (*TaskData, error) {
	var data TaskData

	err := DB.QueryRow(`SELECT task_id, text, time FROM queue FOR UPDATE SKIP LOCKED LIMIT 1;`).Scan(&data.Task_id, &data.Text, &data.Time)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	_, err = DB.Exec("DELETE FROM queue WHERE task_id = $1;", data.Task_id)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
