package models

import "time"

/*
This is the definition of the task struct that will be used
throughout the application. Along with the field names, the
json labels are provided to facilitate the binding process
between the model itself and the JSON format.
*/
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}
