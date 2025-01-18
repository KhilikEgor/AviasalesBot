package domain

type User struct {
	ChatId       int64  `gorm:"primaryKey"`
	UserName     string `gorm:"size:255"`
	Notification bool   `gorm:"type:boolean"`
}
