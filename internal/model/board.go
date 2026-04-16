package model

import (
	"github.com/google/uuid"
	"time"
)

type Board struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_boards_name"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now();index:idx_boards_created_at"`
}

type BoardColumn struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_board_columns_board_id_name"`
	Position  int       `gorm:"not null;default:0"`
	BoardID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_board_columns_board_id_name"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}
