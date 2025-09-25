package models

type Task struct {
	Id         string
	Payload    string
	MaxRetries int
	Status     string
}
