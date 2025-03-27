package modelx

import (
	"gorm.io/gorm"
	"time"
)

type Basic struct {
	AutoID    int64          `gorm:"column:auto_id;uniqueIndex;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
