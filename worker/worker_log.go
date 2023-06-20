package worker

import (
	"context"
	"os"
	"time"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/custom_types"
	"github.com/jackc/pgx/v5"
)

type WorkerLog struct {
	Worker          *Worker
	RequestPayload  map[string]any
	ResponsePayload map[string]any
	Status          custom_types.Status
}

func connection() (conn *pgx.Conn, err error) {
	database_url := os.Getenv("DATABASE_URL")
	conn, err = pgx.Connect(context.Background(), database_url)
	if err != nil {
		panic(err.Error())
	}
	return conn, nil
}

func (workerLog *WorkerLog) Insert() (id string, err error) {
	conn, err := connection()
	defer conn.Close(context.Background())
	err = conn.QueryRow(context.Background(), "INSERT INTO queue (id, request_payload, response_payload, service, priority, status) VALUES ($1, $2, $3, $4, $5, $6) returning id", &workerLog.Worker.Id, &workerLog.RequestPayload, &workerLog.ResponsePayload, &workerLog.Worker.ServiceType, &workerLog.Worker.Options.Priority, &workerLog.Status).Scan(&id)
	if err != nil {
		return "", err
	}
	return
}

func (workerLog *WorkerLog) Update(id string) (err error) {
	conn, err := connection()
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "UPDATE queue set request_payload=$1, response_payload=$2, priority=$3, status=$4, update_at=$5 where id=$6", workerLog.RequestPayload, workerLog.ResponsePayload, workerLog.Worker.Options.Priority, workerLog.Status, time.Now(), id)
	if err != nil {
		return err
	}
	return
}
