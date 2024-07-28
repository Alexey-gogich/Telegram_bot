package models

import "time"

type Reminder struct {
	Id      int
	Chat_id int64
	Text    string
	Time    time.Time
}
