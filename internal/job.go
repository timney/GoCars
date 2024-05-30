package main

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          int
	CreatedAt   time.Time
	JobID       string
	Results     int
	Description string
	CompletedAt time.Time
	ModelID     int
}

func NewJob() *Job {
	return &Job{
		JobID: uuid.NewString(),
	}
}

/*
create table main.job_run
(
    id           integer
        primary key autoincrement,
    created_at   TIMESTAMP default CURRENT_TIMESTAMP not null,
    job_id       text,
    results      bigint,
    description  text,
    completed_at TIMESTAMP,
    model_id     integer
        references main.model
);
*/
