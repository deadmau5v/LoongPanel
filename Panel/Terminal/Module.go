package Terminal

import "time"

type Screen struct {
	Name   string    `json:"name"`
	Id     uint32    `json:"id"`
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}
