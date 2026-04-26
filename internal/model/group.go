package model

import (
	"github.com/google/uuid"
	"time"
)

type Group struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now();index:idx_groups_created_at"`
}
