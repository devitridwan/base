package models

import "time"

type BackoffRetry struct {
	SourceTopic string    `json:"source_topic"`
	Data        []byte    `json:"data"`
	Timestamp   time.Time `json:"timestamp"`
	Error       string    `json:"error"`
}
