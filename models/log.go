package models

import "time"

type PostLog struct {
	ID        int       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	PostID    int       `json:"post_id"`
	Operation string    `json:"operation"`
	Create_at time.Time `json:"create_at"`
}

func (r PostLog) TableName() string {
	return "post_logs"
}