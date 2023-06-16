package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Queue struct {
	Name string
	Request
	Response any
	Priority
	Status
}

type QueueResponse struct {
	Id string
}

func (queue Queue) insert() (response QueueResponse) {
	database_url := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close(context.Background())
	err = conn.QueryRow(context.Background(), "INSERT INTO queue (name, request_payload, response_payload, priority, status) VALUES ($1, $2, $3, $4, $5) returning id", queue.Name, queue.Request, queue.Response, queue.Priority, queue.Status).Scan(&response.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	return
}

func (queue Queue) update(new_queue Queue) (response QueueResponse) {
	database_url := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "UPDATE queue set payload_request=$1, payload_response=$2, priority=$3, status=$4", queue.Request, queue.Response, queue.Priority, queue.Status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	return
}
