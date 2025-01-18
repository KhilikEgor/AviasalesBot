package domain

import "time"

type Vacancy struct {
	VacancyId      int64     `gorm:"primaryKey"`
	Name           string    `gorm:"size:255"`
	Description    string    `gorm:"size:255"`
	Link           string    `gorm:"size:255"`
	PublishDate    time.Time `gorm:"type:timestamp"`
	WithdrawalDate time.Time `gorm:"type:timestamp"` 
	Active         bool      `gorm:"type:boolean"`
}
