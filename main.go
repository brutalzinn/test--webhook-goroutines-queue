package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// lets suppose in some place of world i need to implement a queue to execute web hooks to all the requests failed.
// But i need to do this without dependencies because my i dont want much more thant a simple socket opened to receive do the works for me.
// So we use the gouroutines do this with channels. With this we attack my asks for a queue or a stack principle?!
// Sometimes i do some mistakes with overthinking about the design and patterns concepts before focus the main problem. Now we need to simplifies this cases.
// ITS TIME TO CHANGE THIS.

func worker(job_channel <-chan Job) {
	for job := range job_channel {
		fmt.Printf("start request to %s", job.Request.Url)
		fmt.Printf("start request body %s", job.Request.Body)
		job.Execute()
		fmt.Printf("Job request status %v", job.Status)

	}
}

// i create a main here
func main() {
	job_channel := make(chan Job, 10)
	http.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
		var job Job
		err := json.NewDecoder(r.Body).Decode(&job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		go worker(job_channel)
		job_channel <- job
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received Body: %+v", r.Body)
		fmt.Fprintf(w, "Body: %+v", r.Body)
	})

	http.ListenAndServe(":9000", nil)
}
