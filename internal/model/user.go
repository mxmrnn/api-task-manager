package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName  string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}
