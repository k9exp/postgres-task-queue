package main

import (
	"encoding/json"
	"fmt"
	"k9exp/postgres-task-queue/data"
	"log"
	"net/http"
	"time"
)

type RequestPayload struct {
	Text string `json:"text"`
	Time uint32 `json:"time"`
}

// POST /process
func producer(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	text := requestPayload.Text
	time := requestPayload.Time

	err = data.InsertTask(text, time)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := "Task added in the queue\n"
	w.Write([]byte(msg))
	return
}

func worker(err chan error, worker_id uint16) {
	for {
		data, e := data.GetTask()
		if e != nil {
			err <- e
		}
		if data == nil {
			time.Sleep(10 * time.Second)
			continue
		}

		printText(data, worker_id)

		time.Sleep(1 * time.Second)
	}
}

func printText(data *data.TaskData, worker_id uint16) {
	for i := 0; i < int(data.Time); i++ {
		fmt.Printf("\ttask: %d, by worker: %d> %s [%d/%d]\n", data.Task_id, worker_id, data.Text, i+1, data.Time)
		time.Sleep(500 * time.Millisecond)
	}
}

func app(err chan error) {
	http.HandleFunc("/producer", producer)

	PORT := "4545"
	log.Printf("Serving on http://localhost:%s\n", PORT)
	err <- http.ListenAndServe(":"+PORT, nil)
}

func main() {
	y := data.SetupQueue()
	if y != nil {
		log.Fatal(y)
	}

	err := make(chan error, 1)

	go app(err)

	for i := 0; i < 3; i++ {
		go worker(err, uint16(i+1))
	}

	e := <-err
	log.Printf("Got error: %v\n", e.Error())
}
