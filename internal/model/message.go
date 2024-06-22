package model

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	Room      string `gorm:"index"`
	User      string
	Content   string
	Timestamp int64 `gorm:"autoCreateTime"`
}
