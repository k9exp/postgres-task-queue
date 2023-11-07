package main

import (
	"encoding/json"
	"k9exp/postgres-task-queue/data"
	"log"
	"net/http"
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

	w.Write([]byte("Task added in the queue\n"))
	return
}

func worker(err chan error) {
	// pool the first element the queue
	// do what it required to do
	// repeat
}

func app(err chan error) {
	http.Handle("/", http.FileServer(http.Dir("ui")))

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
	go worker(err)

	e := <-err
	log.Printf("Got error: %v\n", e.Error())
}
